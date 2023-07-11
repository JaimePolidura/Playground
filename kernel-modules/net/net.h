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
#include <linux/skbuff.h>
#include <linux/in6.h>
#include <asm/checksum.h>
#include <linux/ip.h>

#define MY_NET_TIMEOUT 5
#define MY_NET_POOL_SIZE 8
#define MY_NET_USE_NAPI 0

#define MY_NET_RX_INTR 0x0001
#define MY_NET_TX_INTR 0x0002

#define MIN(a, b) a > b ? b : a;

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
	struct sk_buff * sk_buff;
	spinlock_t lock;
	struct net_device * device;
	struct napi_struct napi;
};

int net_open(struct net_device * device);
int net_release(struct net_device * device);
void net_tx_timeout (struct net_device * device, unsigned int txqueue);
int net_start_xmit(struct sk_buff * sk_buff, struct net_device * device);
void net_rx(struct net_device * device, struct my_net_packet * my_net_packet);
int net_poll(struct napi_struct * napi, int budget);

void net_regular_interrupt_handler(int irq, void * dev_id, struct pt_regs * regs);
void net_napi_interrupt_handler(int irq, void * dev_id, struct pt_regs * regs);

static void net_hardware_xmit(char * data_packet, int length_packet, struct net_device * device);
static struct my_net_packet * net_get_packet(struct net_device * device);
static void net_enqueue_packet(struct net_device * device, struct my_net_packet * my_net_packet);
static void net_destroy_pool(struct net_device * device);
static void net_setup_pool(struct net_device * device);
static void net_rx_ints(struct net_device * net_device, int enable);
static void net_release_packet(struct my_net_packet * my_net_packet);
static struct my_net_packet* net_dequeue_packet(struct net_device * device);

void net_init_device(struct net_device * net_device);

static int __init net_init(void);
static void __exit net_exit(void);

#endif