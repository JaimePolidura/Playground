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

struct scull_data {
    char * content;
    int max_size;
    int last_written_index;
    struct cdev cdev;
    struct rw_semaphore sem;
};

int release_scull(struct inode *inode, struct file *file);
int open_scull(struct inode * inode, struct file * file);
ssize_t write_scull(struct file * file, const char __user * buffer, size_t count, loff_t *f_pos);
ssize_t read_scull(struct file * file, char __user * buffer, size_t count, loff_t *f_pos);
long ioctl_scull(struct file *file, unsigned int cmd, unsigned long arg);

struct scull_qset *scull_follow(struct scull_dev *dev, int n);

#endif
