#include "myfifo.h"

int FIFO_MAJOR;
int FIFO_MINOR;

struct fifo_device * FIFO_DEVICE;

struct file_operations FIFO_OPS = {
    .open = fifo_open,
    .read = fifo_read,
    .write = fifo_write,
};

int fifo_open(struct inode * inode, struct file * file) {
    file->private_data = container_of(inode->i_cdev, struct fifo_device, cdev);

    return 0;
}

ssize_t fifo_write(struct file * file, const char __user * buffer, size_t count, loff_t *f_pos) {
    struct fifo_device * fifo_device = file->private_data;

    if(*f_pos > fifo_device->max_size){
        return -ENOMEM;
    }
    if(*f_pos + count >= fifo_device->max_size){
        count = fifo_device->max_size - *f_pos;   
    }     
   
    if(mutex_lock_interruptible(&fifo_device->lock)){
        return -ERESTARTSYS;
    }
    
    if(copy_from_user(fifo_device->content + *f_pos, buffer, count)){
        mutex_unlock(&fifo_device->lock);  
        return -EFAULT;    
    }    

    fifo_device->some_data_present = 1;

    mutex_unlock(&fifo_device->lock);

    wake_up_interruptible(&fifo_device->read_queue);
    
    *f_pos += count;

    return count;
}

ssize_t fifo_read(struct file * file, char __user * buffer, size_t count, loff_t *f_pos) {
    struct fifo_device * fifo_device = file->private_data;

    if(*f_pos >= fifo_device->max_size){
        return 0;
    }    
    if(*f_pos > fifo_device->max_size){
        return -ENOMEM;
    }
    if(*f_pos + count >= fifo_device->max_size){
        count = fifo_device->max_size - *f_pos;   
    }

    if(mutex_lock_interruptible(&fifo_device->lock)){
        return -ERESTARTSYS;
    }

    while(fifo_device->some_data_present == 0){
        mutex_unlock(&fifo_device->lock);
        if(file->f_flags & O_NONBLOCK){
            return -EGAIN;
        }
        if(wait_event_interruptible(fifo_device->read_queue, fifo_device->some_data_present == 0)){
            return -ERESTARTSYS;
        }
        if(mutex_lock_interruptible(&fifo_device->lock)){
            return -ERESTARTSYS;
        }
    }

    if(copy_to_user(buffer, fifo_device->content + *f_pos, count)) {
        mutex_unlock(&fifo_device->lock);
        return -ERESTARTSYS;
    }

    memset(fifo_device->content, 0, fifo_device->max_size);
    fifo_device->some_data_present = 0;

    mutex_unlock(&fifo_device->lock);

    *f_pos += count;

    return count;
}

static int __init fifo_init(void) {
    printk(KERN_ALERT "Initializing fifo!\n");

    FIFO_DEVICE = (struct fifo_device *) allocate(sizeof(struct fifo_device)); 
    
    dev_t dev = 0;
    struct cdev cdev;    
    
    alloc_chrdev_region(&dev, FIFO_MINOR, 1, "myfifo");
    cdev_init(&FIFO_DEVICE->cdev, &FIFO_OPS);
    cdev_add(&FIFO_DEVICE->cdev, dev, 1);

    FIFO_MAJOR = MAJOR(dev);
    FIFO_DEVICE->some_data_present = 0;
    FIFO_DEVICE->content = allocate(FIFO_INITIAL_SIZE);
    FIFO_DEVICE->max_size = FIFO_INITIAL_SIZE;
    FIFO_DEVICE->cdev.owner = THIS_MODULE;

    init_waitqueue_head(&FIFO_DEVICE->read_queue);
    mutex_init(&FIFO_DEVICE->lock);

    printk(KERN_ALERT "Initialized fifo with major number %i\n",FIFO_MAJOR);

    return 0;
}

static void __exit fifo_exit(void) {
    unregister_chrdev_region(MKDEV(FIFO_MAJOR, FIFO_MINOR), 1);
    
    cdev_del(&FIFO_DEVICE->cdev);
    kfree(FIFO_DEVICE->content);
    kfree(FIFO_DEVICE);

    printk(KERN_ALERT "Exited fifo!\n");
}

static char * allocate(size_t size) {
    char * ptr = kmalloc(size, GFP_KERNEL);
    memset(ptr, 0, size);
    
    return ptr;
}

module_init(fifo_init);
module_exit(fifo_exit);

MODULE_LICENSE("Dual BSD/GPL");
MODULE_DESCRIPTION("Module");
MODULE_AUTHOR("Jaime");
