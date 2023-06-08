#include "fifo.h"

int FIFO_MAJOR;
int FIFO_MINOR;

struct fifo_device * FIFO_DEVICE;

static int __init scull_init(void) {
    struct fifo_device * fifo_device = kmalloc(sizeof(struct fifo_device), GFP_KERNEL); 
        
    dev_t dev = 0;
    struct cdev cdev;    

    alloc_chrev_region(&dev, MINOR, 1, "fifo");
    cdev_init(&fifo_device->cdev, &FIFO_OPS);
    cdev_add(&fifo_device->cdev, dev, 1);

    FIFO_MAJOR = MAJOR(dev);
}

static void __exit scull_exit(void) {
}

module_init(fifo_init);
module_exit(fifo_exit);

MODULE_LICENSE("Dual BSD/GPL");
MODULE_DESCRIPTION("Module");
MODULE_AUTHOR("Jaime");
