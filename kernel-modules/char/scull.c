#include <linux/module.h>
#include <linux/sched.h>
#include <linux/types.h>
#include <linux/fs.h>
#include <linux/cdev.h>

int _SCULL_N_DEVS = 0;
int _SCULL_MAJOR;
int _SCULL_MINOR = 0;

struct cdev * _my_cdev;

struct file_operations my_fops = {
};

static int __init scull_init(void) {
    dev_t dev = 0;

    int result = alloc_chrdev_region(&dev, _SCULL_MINOR, _SCULL_N_DEVS, "Scull");
    _SCULL_MAJOR = MAJOR(dev);    
    
    if (result < 0) {
        printk(KERN_WARNING "scull: can't get major %d\n", _SCULL_MAJOR);
        return result;
    }

    _my_cdev = cdev_alloc();
    _my_cdev->ops = &my_fops;
    _my_cdev->owner = THIS_MODULE;

    result = cdev_add(_my_cdev, dev, 1);
    if (result < 0) {
        printk(KERN_WARNING "scull: can't add char driver\n");
        return result;
    }

    printk(KERN_ALERT "Initialized!");

    return 0;
}

static void __exit scull_exit(void) {
    unregister_chrdev_region(MKDEV(_SCULL_MAJOR, _SCULL_MINOR), _SCULL_N_DEVS);
    cdev_del(_my_cdev);
}

module_init(scull_init);
module_exit(scull_exit);

MODULE_LICENSE("Dual BSD/GPL");
MODULE_DESCRIPTION("Module");
MODULE_AUTHOR("Jaime");
