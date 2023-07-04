#ifndef _BLOCK_
#define _BLOCK_

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
#include <linux/vmalloc.h>
#include <linux/blk-mq.h>
#include <linux/hdreg.h>
#include <linux/blkdev.h>

#define KERNEL_SECTOR_SIZE	512
#define MY_BLOCK_N_SECTORS 1024
#define MY_BLOCK_SECTOR_SIZE 512
#define MY_BLOCK_MINORS 16
#define MY_BLOCK_MAJOR 0

struct my_block {
    int size;
    u8 * data;
    spinlock_t lock;
    struct request_queue * queue;
    struct gendisk * gendisk;
    struct blk_mq_tag_set tag_set;
};

static int block_open(struct block_device * bdev, fmode_t mode);

extern void blk_cleanup_queue(struct request_queue * queue);

static int __init block_init(void);
static void __exit block_exit(void);

#endif