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
#define MY_NET_PACKET_POOL_SIZE 8
#define MY_NET_USE_NAPI 0

#define MY_NET_PENDING_READ_PACKET_INTERRUPT 0x0001
#define MY_NET_PENDING_TRANSMISSION 0x0002

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
   struct my_net_packet * packet_pool;
   struct my_net_packet * read_queue;  /* List of incoming packets */
   int read_interruptions_enabled;
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
void net_read(struct net_device * device, struct my_net_packet * my_net_packet);
int net_poll(struct napi_struct * napi, int budget);
int net_create_header_arp(struct sk_buff * sk_buff, struct net_device * device, unsigned short type, const void * destination_address, 
	const void * source_address, unsigned length);
void net_regular_interrupt_handler(int irq, void * dev_id, struct pt_regs * regs);
void net_napi_interrupt_handler(int irq, void * dev_id, struct pt_regs * regs);
struct net_device_stats * net_get_stats(struct net_device * device);

static void net_hardware_xmit(char * data_packet, int length_packet, struct net_device * device);
static struct my_net_packet * net_get_packet_from_packet_pool(struct net_device * device);
static void net_enqueue_packet_read_queue(struct net_device * device, struct my_net_packet * my_net_packet);
static void net_destroy_packet_pool(struct net_device * device);
static void net_setup_packet_pool(struct net_device * device);
static void net_read_enable_interrupts(struct net_device * net_device, int enable);
static void net_release_packet(struct my_net_packet * my_net_packet);
static struct my_net_packet * net_dequeue_packet_from_read_queue(struct net_device * device);

void net_init_device(struct net_device * net_device);

static int __init net_init(void);
static void __exit net_exit(void);

#endif
