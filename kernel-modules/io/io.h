#ifndef _IO_
#define _IO_

#include <linux/kernel.h>
#include <linux/module.h>
#include <linux/pci.h>
#include <linux/init.h>
#include <linux/ioport.h>
#include <asm/io.h>
#include <linux/ioport.h>

#define MY_IO_N_PORTS 8
#define MY_IO_MINOR 0

ssize_t io_read(struct file * file, char __user * buffer, size_t count, loff_t * f_pos);
ssize_t io_write(struct file * file, const char __user * buffer, size_t count, loff_t * f_pos);

static int __init io_init(void);
static void __exit io_exit(void);

#endif