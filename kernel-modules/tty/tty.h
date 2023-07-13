#ifndef _TTY_
#define _TTY_

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
#include <linux/netdevice.h>
#include <linux/etherdevice.h>
#include <linux/etherdevice.h>
#include <linux/skbuff.h>
#include <linux/in6.h>
#include <asm/checksum.h>
#include <linux/ip.h>
#include <linux/tty_driver.h>

#define MY_TTY_MINORS 4	/* 4 Devices */
#define MY_TTY_MAJOR 250
#define MY_TTY_TIMEOUT (HZ * 2)
#define MY_TTY_DATA_CHARACTER	't'
#define MY_TTY_RELEVANT_IFLAG (iflag) ((iflag) & (IGNBRK|BRKINT|IGNPAR|PARMRK|INPCK))

struct my_tty {
    struct tty_struct * tty;
    int open_count;

    struct timer_list timer;
    struct mutex mutex;
};

int tty_open(struct tty_struct * tty, struct file * file);
void tty_close(struct tty_struct * tty, struct file * file);
int tty_write(struct tty_struct * tty, const unsigned char *buffer, int count);
unsigned int tty_write_room(struct tty_struct * tty);

static void tty_timeout(struct timer_list * timer_list);

static int __init tty_init(void);
static void __exit tty_exit(void);

#endif