#include "io.h"

static unsigned long my_io_base = 0x378;
static int my_io_major = 0;
static unsigned long my_io_buffer = 0;
static int my_io_irq = 5;

DECLARE_WAIT_QUEUE_HEAD(my_io_queue);
DECLARE_TASKLET(my_io_tasklet, my_io_tasklet_handler);

struct file_operations my_io_fops = {
	.owner = THIS_MODULE,
    .read = io_read,
    .write = io_write,
    .open = io_open,
};

irqreturn_t my_io_interrupt(int irq, void * dev_id) {
    int value = inb(my_io_base);

    if(!(value & 0x80)) {
        return -IRQ_NONE;
    }

    outb(value & 0x7F, my_io_base); //Clearn device's pending interrupt bit
    
    printk(KERN_ALERT "Interrupt served!");
    wake_up_interruptible(&my_io_queue);
    tasklet_schedule(&my_io_tasklet);

    return IRQ_HANDLED;
}

void my_io_tasklet_handler(struct tasklet_struct * unused) {
    printk(KERN_ALERT "Tasklet served");
}

int io_open(struct inode * inode, struct file * file) {
    file->private_data = kzalloc(sizeof(struct my_io_device), GFP_KERNEL);
    
    return 0;
}

ssize_t io_write(struct file * file, const char __user * buffer, size_t count, loff_t * f_pos) {
    struct my_io_device * my_io_device = file->private_data;
    struct inode * inode = file_dentry(file)->d_inode;
    int minor = iminor(inode);
    void * address_io_to_write = (void *) my_io_base + (minor & 0x0F);

    unsigned char * ptr_data_to_write = kmalloc(count, GFP_KERNEL);
    copy_from_user(ptr_data_to_write, buffer, count);
    
    unsigned char * actual_ptr_data_to_write = ptr_data_to_write;

    my_io_device->can_read == 1;

    while (count--) {
        iowrite8(*actual_ptr_data_to_write++, address_io_to_write);
        wmb();
    }

    //Raise interrupt (no se como xd)

    kfree(ptr_data_to_write);

    return count;
}

ssize_t io_read(struct file * file, char __user * buffer, size_t count, loff_t * f_pos) {
    struct my_io_device * my_io_device = file->private_data;

    while(my_io_device->can_read == 0){
        if(wait_event_interruptible(my_io_queue, my_io_device->can_read == 0)){
            return -ERESTARTSYS;
        }
    }

    struct inode * inode = file_dentry(file)->d_inode;
    int minor = iminor(inode);
    void * address_io_to_read = (void *) my_io_base + (minor & 0x0f);
    unsigned char * ptr_data_buffer = kmalloc(count, GFP_KERNEL);

    unsigned char * actual_ptr_data_buffer = ptr_data_buffer;
    size_t actual_count = count;

    while (actual_count--) {
		*actual_ptr_data_buffer++ = ioread8(address_io_to_read);
		rmb();
	}

    copy_to_user(buffer, ptr_data_buffer, count);
    kfree(ptr_data_buffer);
    my_io_device->can_read = 0;

    return 0;
}

static int __init io_init(void) {
    request_region(my_io_base, MY_IO_N_PORTS, "my_io");
    
    request_mem_region(my_io_base, MY_IO_N_PORTS, "my_io");
    my_io_base = (unsigned long) ioremap(my_io_base, MY_IO_N_PORTS);
    
    my_io_major = register_chrdev(0, "my_io", &my_io_fops);
    my_io_buffer = __get_free_pages(GFP_KERNEL, 1);

    request_irq(my_io_irq, my_io_interrupt, IRQF_SHARED, "my_io", my_io_interrupt);
    outb(0x10, my_io_base + 2);

    return 0;
}

static void __exit io_exit(void) {
    free_irq(my_io_irq, NULL);
    outb(0x00, my_io_base + 2);

    release_mem_region(my_io_base, MY_IO_N_PORTS);
    iounmap(my_io_base);
    
	unregister_chrdev(my_io_major, "my_io");
}

MODULE_LICENSE("GPL");

module_init(io_init);
module_exit(io_exit);