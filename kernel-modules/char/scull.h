#ifndef _SCULL_H_
#define _SCULL_H_

#include <linux/module.h>
#include <linux/sched.h>
#include <linux/types.h>
#include <linux/fs.h>
#include <linux/cdev.h>
#include <linux/kernel.h>
#include <asm/uaccess.h>

struct scull_qset {
	void **data;
	struct scull_qset *next;
};

struct scull_dev {
    struct scull_qset *data; /* Pointer to first quantum set */
    int quantum; /* Tamaño máximo de los buffers en **data */
    int qset; /* Tamaño máximo de punteros a buffers en **data */
    unsigned long size; /* Suma total de todos los bytes guardados void**data */
    unsigned int access_key; /* used by sculluid and scullpriv */
    struct semaphore sem; /* mutual exclusion semaphore */
    struct cdev cdev; /* Char device structure */
};

int release_scull(struct inode *inode, struct file *file);
int open_scull(struct inode * inode, struct file * file);
ssize_t write_scull(struct file * file, const char __user * buffer, size_t count, loff_t *f_pos);
ssize_t read_scull(struct file * file, char __user * buffer, size_t count, loff_t *f_pos);

struct scull_qset *scull_follow(struct scull_dev *dev, int n);

#endif