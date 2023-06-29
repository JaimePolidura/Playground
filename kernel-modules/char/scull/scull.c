#include "scull.h"

struct scull * my_scull;

int _SCULL_MAJOR;
int _SCULL_MINOR = 0;

struct file_operations my_fops = {
    .owner = THIS_MODULE,
    .open = scull_open,
    .release = scull_release,
    .read = scull_read,
    .write = scull_write,
    .unlocked_ioctl = scull_ioctl,
};

long scull_ioctl(struct file *file, unsigned int cmd, unsigned long arg) {
    if(!capable(CAP_SYS_ADMIN))
        return -EPERM;

    struct scull * scull = file->private_data;

    switch(cmd){
        case SCULL_IOCTL_GROW:
            unsigned short newSize = (unsigned short) arg;            
            if(newSize <= 0 || newSize <= scull->max_size) 
                return -EINVAL;

            char * newPtr = kzalloc(sizeof(char) * newSize, GFP_KERNEL);
            if(newPtr == NULL)
                return -ENOMEM;

            memset(scull->content, 0, scull->max_size);

            up_write(&scull->sem);
	            
            memcpy(newPtr, scull->content, scull->max_size);            
            kfree(scull->content);
            scull->content = newPtr;
            scull->max_size = newSize;        

            down_write(&scull->sem);

            break;

        default:
            return -EPERM;
    }

    return 0;    
}

ssize_t scull_write(struct file * file, const char __user * buffer, size_t count, loff_t *f_pos) {
    struct scull * scull = file->private_data;
	
    if(*f_pos > scull->max_size){
        return -ENOMEM;
    }
    if(*f_pos + count >= scull->max_size){
        count = scull->max_size - *f_pos - 1;   
    }     
   
    down_write(&scull->sem);
 
    if(copy_from_user(scull->content + *f_pos, buffer, count)){
        up_write(&scull->sem);    
        return -EFAULT;    
    }    

    printk(KERN_ALERT "Written %i bytes to file on posicion %lld\n", count, *f_pos);

    scull->last_written_index += count;
    
    up_write(&scull->sem);

    *f_pos += count;

    return count;
}

ssize_t scull_read(struct file *file, char __user *buffer, size_t count, loff_t *f_pos) {
    struct scull * scull = file->private_data;
    
    if(*f_pos >= scull->max_size){
        return -ENOMEM;
    }
    if(*f_pos + count >= scull->max_size) {
        count = scull->max_size - *f_pos - 1;
    }
    
    down_read(&scull->sem);

    if(copy_to_user(buffer, scull->content + *f_pos, count)){
    	up_read(&scull->sem);
        return -EFAULT;
    }
	
    up_read(&scull->sem);

    printk(KERN_ALERT "Read %i bytes on position %lld\n", count, *f_pos);

    *f_pos += count;

    return count;
}

int scull_release(struct inode *inode, struct file *file) {
    return 0;
}

int scull_open(struct inode * inode, struct file * file) {
    file->private_data = container_of(inode->i_cdev, struct scull, cdev);
    
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

    my_scull = kzalloc(sizeof(struct scull), GFP_KERNEL);
    
    struct rw_semaphore * sem;
    init_rwsem(sem);

    cdev_init(&my_scull->cdev, &my_fops);
    my_scull->cdev.owner = THIS_MODULE;
    my_scull->max_size = SCULL_INITIAL_MAX_CONTENT_SIZE;    
    my_scull->sem = *sem;
    my_scull->content = kzalloc(SCULL_INITIAL_MAX_CONTENT_SIZE, GFP_KERNEL);

    result = cdev_add(&my_scull->cdev, dev, 1);
    if (result < 0) {
        printk(KERN_WARNING "scull: can't add char driver\n");
        return result;
    }

    printk(KERN_ALERT "Initialized with major number %d!\n", _SCULL_MAJOR);

    return 0;
}

static void __exit scull_exit(void) {
    unregister_chrdev_region(MKDEV(_SCULL_MAJOR, _SCULL_MINOR), 1);
    
    cdev_del(&my_scull->cdev);
    kfree(my_scull->content);
    kfree(my_scull);

    printk(KERN_ALERT "Exited!");
}


module_init(scull_init);
module_exit(scull_exit);

MODULE_LICENSE("Dual BSD/GPL");
MODULE_DESCRIPTION("Module");
MODULE_AUTHOR("Jaime");
