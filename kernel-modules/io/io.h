#ifndef _IO_
#define _IO_

#include <linux/kernel.h>
#include <linux/module.h>
#include <linux/pci.h>
#include <linux/init.h>
#include <linux/ioport.h>
#include <asm/io.h>
#include <linux/ioport.h>
#include <linux/interrupt.h>

#define MY_IO_N_PORTS 8
#define MY_IO_MINOR 0

struct my_io_device {
    int can_read;
};

irqreturn_t my_io_interrupt(int irq, void * dev_id);

ssize_t io_read(struct file * file, char __user * buffer, size_t count, loff_t * f_pos);
ssize_t io_write(struct file * file, const char __user * buffer, size_t count, loff_t * f_pos);
int io_open(struct inode * inode, struct file * file);

static int __init io_init(void);
static void __exit io_exit(void);

#endif