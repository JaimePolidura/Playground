#ifndef _SCULL_H_
#define _SCULL_H_

#include <linux/rwsem.h>
#include <linux/module.h>
#include <linux/sched.h>
#include <linux/types.h>
#include <linux/fs.h>
#include <linux/cdev.h>
#include <linux/kernel.h>
#include <asm/uaccess.h>
#include <linux/ioctl.h>
#include <linux/capability.h>
#include <linux/mm.h>
#include <linux/aio.h>
#include <linux/workqueue.h>
#include <linux/fs.h>

#define SCULL_IOCTL_MAGIC 'J'
#define SCULL_IOCTL_GROW _IOW(SCULL_IOCTL_MAGIC, 0, unsigned short)
#define SCULL_INITIAL_MAX_CONTENT_SIZE 16384

struct scull {
    char * content;
    int max_size;
    int last_written_index;
    struct cdev cdev;
    struct rw_semaphore sem;
};

struct scull_aio_async_work {
    struct kiocb * iocb;
    struct work_struct work;
    const char * buffer;
    size_t count;
    loff_t pos;
};

void scull_mmap_vma_open(struct vm_area_struct * vma);
void scull_mmap_vma_close(struct vm_area_struct * vma);
vm_fault_t scull_mmap_vma_fault(struct vm_fault *vmf);

int scull_release(struct inode *inode, struct file *file);
int scull_open(struct inode * inode, struct file * file);
ssize_t scull_write(struct file * file, const char __user * buffer, size_t count, loff_t *f_pos);
ssize_t scull_read(struct file * file, char __user * buffer, size_t count, loff_t *f_pos);
long scull_ioctl(struct file *file, unsigned int cmd, unsigned long arg);
int scull_mmap(struct file * file, struct vm_area_struct * vma);
ssize_t scull_aio_write(struct kiocb * iocb, const char __user * buffer, size_t count, loff_t pos);

void aio_async_write(void * data);

#endif
