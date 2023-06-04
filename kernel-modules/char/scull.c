#include "scull.h"

/*
    Ignoro los errores que pude producir kmalloc
*/

int _SCULL_N_DEVS = 0;
int _SCULL_MAJOR;
int _SCULL_MINOR = 0;

struct scull_dev * _scull_device;

struct file_operations my_fops = {
    .owner = THIS_MODULE,
    .open = open_scull,
    .release = release_scull,
    .read = read_scull,
};

ssize_t read_scull(struct file *file, char __user *buffer, size_t count, loff_t *f_pos) {
    struct scull_dev * dev = filp->private_data;
    
    if(*f_pos > dev->size){
        return 0;
    }
    if(*f_pos + count > dev->size) {
        count = dev->size - *f_pos;
    }

    int qset = dev->qset;
    int quantum = dev->quantum;
    int itemSizeBytes = qset * quantum;
    int item = (long) *f_pos / itemsize;
    int rest = (long) *f_pos % itemsize;
    int s_pos = rest / quantum;
    struct scull_qset *dptr = scull_follow(dev, item);

    if (dptr == NULL || !dptr->data || ! dptr->data[s_pos]) {
        return 0;
    }

	if (count > quantum - q_pos) {
		count = quantum - q_pos;
    }

	if (copy_to_user(buf, dptr->data[s_pos] + q_pos, count)) {
        return -EFAULT;
	}

    *f_pos += count;
	retval = count; 

    return count;
}

struct scull_qset *scull_follow(struct scull_dev *dev, int n) {
	struct scull_qset *qs = dev->data;

	if (! qs) {
		qs = dev->data = kmalloc(sizeof(struct scull_qset), GFP_KERNEL);
		memset(qs, 0, sizeof(struct scull_qset));
	}

	while (n--) {
		if (!qs->next) {
			qs->next = kmalloc(sizeof(struct scull_qset), GFP_KERNEL);
			memset(qs->next, 0, sizeof(struct scull_qset));
		}
		qs = qs->next;
	}

	return qs;
}

int release_scull(struct inode *inode, struct file *file) {
    return 0;
}

int open_scull(struct inode * inode, struct file * file) {
    file->private_data = container_of(inode->i_cdev, struct scull_dev, cdev);

     if ((file->f_flags & O_ACCMODE) == O_WRONLY) {
        //scull_trim(...)
     }

    return 0;
}

static int __init scull_init(void) {
    dev_t dev = 0;

    int result = alloc_chrdev_region(&dev, _SCULL_MINOR, _SCULL_N_DEVS, "Scull");
    _SCULL_MAJOR = MAJOR(dev);    
    
    if (result < 0) {
        printk(KERN_WARNING "scull: can't get major %d\n", _SCULL_MAJOR);
        return result;
    }

    _scull_device = kmalloc(sizeof(struct scull_dev), GFP_KERNEL);
    memset(_scull_device, 0, sizeof(struct scull_dev));

    cdev_init(&_scull_device->cdev, &my_fops);
    _scull_device->cdev.owner = THIS_MODULE;

    result = cdev_add(&_scull_device->cdev, dev, 1);
    if (result < 0) {
        printk(KERN_WARNING "scull: can't add char driver\n");
        return result;
    }

    printk(KERN_ALERT "Initialized!");

    return 0;
}

static void __exit scull_exit(void) {
    unregister_chrdev_region(MKDEV(_SCULL_MAJOR, _SCULL_MINOR), _SCULL_N_DEVS);
    cdev_del(&_scull_device->cdev);
    kfree(_scull_device);

    printk(KERN_ALERT "Exited!");
}

module_init(scull_init);
module_exit(scull_exit);

MODULE_LICENSE("Dual BSD/GPL");
MODULE_DESCRIPTION("Module");
MODULE_AUTHOR("Jaime");
