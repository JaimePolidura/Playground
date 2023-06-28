#include "io.h"

static unsigned long my_io_base = 0x378;
static int my_io_major = 0;
static unsigned long my_io_buffer = 0;

struct file_operations my_io_fops = {
	.owner = THIS_MODULE,
    .read = io_read,
    .write = io_write,
};

ssize_t io_write(struct file * file, const char __user * buffer, size_t count, loff_t * f_pos) {
    struct inode * inode = file_dentry(file)->d_inode;
    int minor = iminor(inode);
    void * address_io_to_write = (void *) my_io_base + (minor & 0x0f);

    unsigned char * ptr_data_to_write = kmalloc(count, GFP_KERNEL);
    copy_from_user(ptr_data_to_write, buffer, count);
    
    unsigned char * actual_ptr_data_to_write = ptr_data_to_write;
    
    while (count--) {
        iowrite8(*actual_ptr_data_to_write++, address_io_to_write);
        wmb();
    }

    kfree(ptr_data_to_write);

    return count;
}

ssize_t io_read(struct file * file, char __user * buffer, size_t count, loff_t * f_pos) {
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

    return 0;
}

static int __init io_init(void) {
    request_region(my_io_base, MY_IO_N_PORTS, "my_io");
    
    request_mem_region(my_io_base, MY_IO_N_PORTS, "my_io");
    my_io_base = (unsigned long) ioremap(my_io_base, MY_IO_N_PORTS);
    
    my_io_major = register_chrdev(0, "my_io", &my_io_fops);
    my_io_buffer = __get_free_pages(GFP_KERNEL, 1);

    return 0;
}

static void __exit io_exit(void) {
    release_mem_region(my_io_base, MY_IO_N_PORTS);
    iounmap(my_io_base);
    
	unregister_chrdev(my_io_major, "short");
}

MODULE_LICENSE("GPL");

module_init(io_init);
module_exit(io_exit);