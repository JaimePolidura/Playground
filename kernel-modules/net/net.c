#include "net.h"

struct net_device * net_devices [2];

void net_init_device(struct net_device * net_device) {
    ether_setup(net_device);

    net_device->watchdog_timeot = MY_NET_TIMEOUT;
    net_device->netdev_ops = &my_net_netdev_ops;
    net_device->header_ops = &my_net_header_ops;
    net_device->flags |= IFF_NOARP;
    net_device->features = |= NETIF_F_HW_CSUM;

    struct my_net * my_net = netdev_priv(net_device);
	spin_lock_init(&priv->lock);
	my_net->device = net_device;

    //TODO
    snull_rx_ints(dev, 1);		/* enable receive interrupts */
	snull_setup_pool(dev);
}

static int __init net_init(void) {
    net_devices[0] = alloc_netdev(sizeof(struct my_net), "sn%d", net_init_device);
    net_devices[1] = alloc_netdev(sizeof(struct my_net), "sn%d", net_init_device);
    
    register_netdev(net_devices[0]);
    register_netdev(net_devices[1]);

    return 0;
}

static void __exit net_exit(void) {
    unregister_netdev(net_devices[0]);
    unregister_netdev(net_devices[1]);
}

module_init(net_init);
module_exit(net_exit);

MODULE_LICENSE("Dual BSD/GPL");
MODULE_DESCRIPTION("Module");
MODULE_AUTHOR("Jaime");