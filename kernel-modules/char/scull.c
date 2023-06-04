#include "scull.h"

/*
    Ignoro los errores que pude producir kmalloc
*/

static void initializeScullDataContentPtr(struct scull_data * scull_data);

#define MAX_CONTENT_SIZE 4096

struct scull_data * _scull_data;

int _SCULL_N_DEVS = 0;
int _SCULL_MAJOR;
int _SCULL_MINOR = 0;

struct file_operations my_fops = {
    .owner = THIS_MODULE,
    .open = open_scull,
    .release = release_scull,
    .read = read_scull,
    .write = write_scull,
};

ssize_t write_scull(struct file * file, const char __user * buffer, size_t count, loff_t *f_pos) {
    struct scull_data * scull_data = file->private_data;
	
    if(*f_pos > scull_data->max_size){
        return -ENOMEM;
    }
    if(*f_pos + count >= scull_data->max_size){
        count = scull_data->max_size - *f_pos - 1;   
    }     
    
    if(copy_from_user(scull_data->content + *f_pos, buffer, count)){
        return -EFAULT;    
    }    

    scull_data->last_written_index += count;
    *f_pos += count;

	return count;
}

ssize_t read_scull(struct file *file, char __user *buffer, size_t count, loff_t *f_pos) {
    struct scull_data * scull_data = file->private_data;
    
    if(*f_pos >= scull_data->max_size){
        return -ENOMEM;
    }
    if(*f_pos + count >= scull_data->max_size) {
        count = scull_data->max_size - *f_pos - 1;
    }
    if(!scull_data->content){
	    initializeScullDataContentPtr(scull_data);
    }

    if(copy_to_user(buffer, scull_data->content + *f_pos, count)){
    	return -EFAULT;
    }
	
    *f_pos += count;

    return count;
}

static void initializeScullDataContentPtr(struct scull_data * scull_data) {
    scull_data->content = kmalloc(sizeof(struct scull_data) * MAX_CONTENT_SIZE, GFP_KERNEL);
    memset(scull_data->content, 0, MAX_CONTENT_SIZE);
}

int release_scull(struct inode *inode, struct file *file) {
    return 0;
}

int open_scull(struct inode * inode, struct file * file) {
    file->private_data = container_of(inode->i_cdev, struct scull_data, cdev);

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

    _scull_data = kmalloc(sizeof(struct scull_data), GFP_KERNEL);
    memset(_scull_data, 0, sizeof(struct scull_data));

    cdev_init(&_scull_data->cdev, &my_fops);
    _scull_data->cdev.owner = THIS_MODULE;
    _scull_data->max_size = MAX_CONTENT_SIZE;    

    result = cdev_add(&_scull_data->cdev, dev, 1);
    if (result < 0) {
        printk(KERN_WARNING "scull: can't add char driver\n");
        return result;
    }

    printk(KERN_ALERT "Initialized!");

    return 0;
}

static void __exit scull_exit(void) {
    unregister_chrdev_region(MKDEV(_SCULL_MAJOR, _SCULL_MINOR), _SCULL_N_DEVS);
    cdev_del(&_scull_data->cdev);
    kfree(_scull_data->content);
    kfree(_scull_data);

    printk(KERN_ALERT "Exited!");
}


module_init(scull_init);
module_exit(scull_exit);

MODULE_LICENSE("Dual BSD/GPL");
MODULE_DESCRIPTION("Module");
MODULE_AUTHOR("Jaime");
