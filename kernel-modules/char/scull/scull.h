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

#define SCULL_IOCTL_MAGIC 'J'
#define SCULL_IOCTL_GROW _IOW(SCULL_IOCTL_MAGIC, 0, unsigned short)
#define SCULL_INITIAL_MAX_CONTENT_SIZE 4096

struct scull {
    char * content;
    int max_size;
    int last_written_index;
    struct cdev cdev;
    struct rw_semaphore sem;
};

int scull_release(struct inode *inode, struct file *file);
int scull_open(struct inode * inode, struct file * file);
ssize_t scull_write(struct file * file, const char __user * buffer, size_t count, loff_t *f_pos);
ssize_t scull_read(struct file * file, char __user * buffer, size_t count, loff_t *f_pos);
long scull_ioctl(struct file *file, unsigned int cmd, unsigned long arg);

#endif
