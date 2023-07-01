#include "scull.h"

struct scull * my_scull;

int _SCULL_MAJOR;
int _SCULL_MINOR = 0;

struct file_operations scull_ops = {
    .owner = THIS_MODULE,
    .open = scull_open,
    .release = scull_release,
    .read = scull_read,
    .write = scull_write,
    .unlocked_ioctl = scull_ioctl,
    .mmap = scull_mmap,
    .aio_write = scull_aio_write,
};

struct vm_operations_struct scull_mmap_vma_ops = {
    .open = scull_mmap_vma_open,
    .close = scull_mmap_vma_close,
    .fault = scull_mmap_vma_fault,
}; 

int scull_mmap(struct file * file, struct vm_area_struct * vma) {
    vma->vm_private_data = file->private_data;
    vma->vm_ops = &scull_mmap_vma_ops;

    scull_mmap_vma_open(vma);

    return 0;
}

vm_fault_t scull_mmap_vma_fault(struct vm_fault *vmf) {
    struct vm_area_struct * vma = vmf->vma;
    unsigned long page_offset = vmf->pgoff << PAGE_SHIFT;
    unsigned long bit_offset = vmf->address - vma->vm_start;
    unsigned long offset = page_offset + bit_offset;
    struct scull * scull = vma->vm_private_data;

    if(offset > scull->max_size){
        return NULL;
    }

    void * new_ptr = scull->content + offset;

    struct page * page = virt_to_page(new_ptr);
    get_page(page);
    vmf->page = page;

    return 0;
}

void scull_mmap_vma_open(struct vm_area_struct * vma) {
    printk(KERN_ALERT "Forked from other process");
}

void scull_mmap_vma_close(struct vm_area_struct * vma) {
    printk(KERN_ALERT "Closed");
}

ssize_t scull_aio_write(struct kiocb * iocb, const char __user * buffer, size_t count, loff_t pos) {
    if(is_sync_kiocb(iocb)){
        return scull_write(iocb->ki_filp, buffer, count, &pos);
    }

    struct scull * scull = iocb->ki_filp->private_data;
    
    if(pos > scull->max_size){
        return -ENOMEM;
    }
    if(pos + count >= scull->max_size){
        count = scull->max_size - pos - 1;   
    }

    struct scull_aio_async_work * async_work = kzalloc(sizeof(struct scull_aio_async_work), GFP_KERNEL);

    async_work->buffer = kzalloc(count, GFP_KERNEL);  
    async_work->iocb = iocb;
    async_work->count = count;
    async_work->pos = pos;
    copy_from_user(async_work->buffer, buffer, count);  

    DECLARE_WORK(aio_async_write, (void *) async_work);

    return -EIOCBQUEUED;
}

void aio_async_write(void * data) {
    struct scull_aio_async_work * async_write_work = (struct scull_aio_async_work *) data;
    struct scull * scull = (struct scull *) async_write_work->iocb->ki_filp->private_data;

    down_write(&scull->sem);
    memcpy(scull->content + async_write_work->pos, async_write_work->buffer, async_write_work->count);
    up_write(&scull->sem);

    async_write_work->iocb->ki_complete(async_write_work->iocb, async_write_work->count);

    kfree(async_write_work->buffer);
    kfree(async_write_work);
}

long scull_ioctl(struct file *file, unsigned int cmd, unsigned long arg) {
    if(!capable(CAP_SYS_ADMIN))
        return -EPERM;

    struct scull * scull = file->private_data;

    switch(cmd){
        case SCULL_IOCTL_GROW:
            unsigned short newSize = (unsigned short) arg;            
            if(newSize <= 0 || newSize <= scull->max_size) 
                return -EINVAL;

            char * newPtr = kzalloc(newSize, GFP_KERNEL);
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

    cdev_init(&my_scull->cdev, &scull_ops);
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
