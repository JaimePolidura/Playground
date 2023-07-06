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
#include <linux/timer.h>

#define KERNEL_SECTOR_SIZE	512
#define MY_BLOCK_N_SECTORS 1024
#define MY_BLOCK_SECTOR_SIZE 512
#define MY_BLOCK_MINORS 16
#define MY_BLOCK_MAJOR 0
#define MY_BLOCK_INVALIDATE_DELAY 30 * HZ

struct my_block {
    int size;
    u8 * data;
    spinlock_t lock;
    struct request_queue * queue;
    struct gendisk * gendisk;
    struct blk_mq_tag_set tag_set;
    struct timer_list timer;
    short media_changed;
    short n_users;
};

static int block_open(struct block_device * bdev, fmode_t mode);
static void block_release(struct gendisk *disk, fmode_t mode);
int block_ioctl (struct block_device *bdev, fmode_t mode, unsigned int cmd, unsigned long arg);
blk_status_t block_queue_request(struct blk_mq_hw_ctx * blk_mq_hw_ctx, const struct blk_mq_queue_data * blk_mq_queue_data);
void block_submit_bio(struct bio * bio);

void block_transfer(struct my_block * dev, unsigned long sector, unsigned long sector_number, char * buffer, int write);

void block_timer_timeout(struct timer_list * ldev);
extern void blk_cleanup_queue(struct request_queue * queue);

static int __init block_init(void);
static void __exit block_exit(void);

#endif