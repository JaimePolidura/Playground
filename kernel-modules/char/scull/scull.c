#include "scull.h"

/*
    Ignoro los errores que pude producir kmalloc
*/

static void initializeScullDataContentPtr(struct scull_data * scull_data);

struct scull_data * _scull_data;

int _SCULL_MAJOR;
int _SCULL_MINOR = 0;

struct file_operations my_fops = {
    .owner = THIS_MODULE,
    .open = open_scull,
    .release = release_scull,
    .read = read_scull,
    .write = write_scull,
    .unlocked_ioctl = ioctl_scull,
};

long ioctl_scull(struct file *file, unsigned int cmd, unsigned long arg) {
    if(!capable(CAP_SYS_ADMIN))
        return -EPERM;

    struct scull_data * scull_data = file->private_data;

    switch(cmd){
        case SCULL_IOCTL_GROW:
            unsigned short newSize = (unsigned short) arg;            
            if(newSize <= 0 || newSize <= scull_data->max_size) 
                return -EINVAL;

            char * newPtr = kmalloc(sizeof(char) * newSize, GFP_KERNEL);
            if(newPtr == NULL)
                return -ENOMEM;

            memset(scull_data->content, 0, scull_data->max_size);

            up_write(&scull_data->sem);
	            
            memcpy(newPtr, scull_data->content, scull_data->max_size);            
            kfree(scull_data->content);
            scull_data->content = newPtr;
            scull_data->max_size = newSize;        

            down_write(&scull_data->sem);

            break;

        default:
            return -EPERM;
    }

    return 0;    
}

ssize_t write_scull(struct file * file, const char __user * buffer, size_t count, loff_t *f_pos) {
    struct scull_data * scull_data = file->private_data;
	
    if(*f_pos > scull_data->max_size){
        return -ENOMEM;
    }
    if(*f_pos + count >= scull_data->max_size){
        count = scull_data->max_size - *f_pos - 1;   
    }     
   
    down_write(&scull_data->sem);
 
    if(copy_from_user(scull_data->content + *f_pos, buffer, count)){
        up_write(&scull_data->sem);    
        return -EFAULT;    
    }    

    printk(KERN_ALERT "Written %i bytes to file on posicion %lld\n", count, *f_pos);

    scull_data->last_written_index += count;
    
    up_write(&scull_data->sem);

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
    
    down_read(&scull_data->sem);

    if(copy_to_user(buffer, scull_data->content + *f_pos, count)){
    	up_read(&scull_data->sem);
        return -EFAULT;
    }
	
    up_read(&scull_data->sem);

    printk(KERN_ALERT "Read %i bytes on position %lld\n", count, *f_pos);

    *f_pos += count;

    return count;
}

static void initializeScullDataContentPtr(struct scull_data * scull_data) {
    scull_data->content = kmalloc(sizeof(char) * SCULL_INITIAL_MAX_CONTENT_SIZE, GFP_KERNEL);
    memset(scull_data->content, 0, SCULL_INITIAL_MAX_CONTENT_SIZE);
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

    int result = alloc_chrdev_region(&dev, _SCULL_MINOR, 1, "scull");
    _SCULL_MAJOR = MAJOR(dev);    
 
    if (result < 0) {
        printk(KERN_WARNING "scull: can't get major %d\n", _SCULL_MAJOR);
        return result;
    }

    _scull_data = kmalloc(sizeof(struct scull_data), GFP_KERNEL);
    memset(_scull_data, 0, sizeof(struct scull_data));

    struct rw_semaphore * sem;
    init_rwsem(sem);

    cdev_init(&_scull_data->cdev, &my_fops);
    _scull_data->cdev.owner = THIS_MODULE;
    _scull_data->max_size = SCULL_INITIAL_MAX_CONTENT_SIZE;    
    _scull_data->sem = *sem;
    initializeScullDataContentPtr(_scull_data);

    result = cdev_add(&_scull_data->cdev, dev, 1);
    if (result < 0) {
        printk(KERN_WARNING "scull: can't add char driver\n");
        return result;
    }

    printk(KERN_ALERT "Initialized with major number %d!\n", _SCULL_MAJOR);

    return 0;
}

static void __exit scull_exit(void) {
    unregister_chrdev_region(MKDEV(_SCULL_MAJOR, _SCULL_MINOR), 1);
    
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
