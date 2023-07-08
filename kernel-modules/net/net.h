/**
 * /etc/networks
 * 
 * snullnet0 192.168.0.0
 * snullnet1 192.168.1.0
 * 
 * /etc/hosts
 * 192.168.0.1 local0
 * 192.168.0.2 remote0
 * 192.168.1.2 local1
 * 192.168.1.1 remote1
 * 
*/

#ifndef _NET_
#define _ENT_

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

#define MY_NET_TIMEOUT 5

struct my_net_packet {
    struct my_net_packet * next;
	struct net_device * device;
	int	length;
	u8 data [ETH_DATA_LEN];
};

struct my_net {
	struct net_device_stats stats;
	int status;
    struct my_net_packet * ppool;
	struct my_net_packet * rx_queue;  /* List of incoming packets */
	int rx_int_enabled;
	int tx_packetlen;
	u8 * tx_packetdata;
	struct sk_buff * skb;
	spinlock_t lock;
	struct net_device * device;
	struct napi_struct napi;
};

void net_init_device(struct net_device * net_device);

static int __init net_init(void);
static void __exit net_exit(void);

#endif