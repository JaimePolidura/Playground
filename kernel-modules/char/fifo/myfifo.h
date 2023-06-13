#ifndef _FIFO_H_
#define _SCULL_H_

#include <linux/rwsem.h>
#include <linux/module.h>
#include <linux/sched.h>
#include <linux/types.h>
#include <linux/fs.h>
#include <linux/cdev.h>
#include <linux/mutex.h>
#include <linux/kernel.h>
#include <asm/uaccess.h>
#include <linux/ioctl.h>
#include <linux/capability.h>
#include <linux/fs.h>
#include <linux/poll.h>

#define FIFO_INITIAL_SIZE 1024

static char * allocate(size_t size);

struct fifo_device {
    char * content;
    struct cdev cdev;
    int max_size;
    struct mutex lock;
    wait_queue_head_t read_queue;
    int some_data_present;
    struct fasync_struct * async_queue;
};

int fifo_fasync(int fd, struct file * file, int mode);
int fifo_release(struct inode * inode, struct file * file);
int fifo_open(struct inode * inode, struct file * file);
ssize_t fifo_write(struct file * file, const char __user * buffer, size_t count, loff_t *f_pos);
ssize_t fifo_read(struct file * file, char __user * buffer, size_t count, loff_t *f_pos);
unsigned int fifo_poll(struct file * file, poll_table * wait);

#endif
