#include "net.h"

static void (*net_interrupt_handler)(int, void *, struct pt_regs *);

struct net_device * net_devices [2];

static const struct header_ops my_net_header_ops = {
    .create = net_create_header_arp
};

static const struct net_device_ops my_net_netdev_ops = {
    .ndo_open = net_open,
    .ndo_tx_timeout = net_tx_timeout,
    .ndo_start_xmit = net_start_xmit,
    .ndo_stop = net_release,
    .ndo_get_stats = net_get_stats,
};

void net_read(struct net_device * device, struct my_net_packet * my_net_packet) {
    struct my_net * my_net = netdev_priv(device);

    struct sk_buff * sk_buff = dev_alloc_skb(my_net_packet->length + 2);
    memcpy(skb_put(sk_buff, my_net_packet->length), my_net_packet->data, my_net_packet->length);

    sk_buff->dev = device;
    sk_buff->protocol = eth_type_trans(sk_buff, device);
    sk_buff->ip_summed = CHECKSUM_UNNECESSARY;
    	
    my_net->stats.rx_bytes += my_net_packet->length;
    my_net->stats.rx_packets++;

    netif_rx(sk_buff); //Manda el socket a capas superiores
}

int net_start_xmit(struct sk_buff * sk_buff, struct net_device * device) {
    struct my_net * my_net = netdev_priv(device);
    int length_packet = MIN(ETH_ZLEN, sk_buff->len);
    char * data_packet = sk_buff->data;

    if(sk_buff->len < ETH_ZLEN) {
        char shortpkt[ETH_ZLEN];
        memset(shortpkt, 0, ETH_ZLEN);
        memcpy(shortpkt, sk_buff->data, sk_buff->len);
        
        data_packet = shortpkt;
    }

    netif_trans_update(device); //Informar al kernel sobre cambios
    
    my_net->sk_buff = sk_buff;

    net_hardware_xmit(data_packet, length_packet, device);

    return 0;
}

static void net_hardware_xmit(char * data_packet, int length_packet, struct net_device * source_device) {
    struct iphdr * ip_header = (struct iphdr *)(data_packet + sizeof(struct ethhdr));
    u8 * destination_ip = (u8 *) &ip_header->daddr;
    u8 * source_ip = (u8 *) &ip_header->saddr;

    //Cambiamos 3º byte de la case C (255.255.255.0) que indica la red
    //El ^ es xor
	((u8 *)destination_ip)[2] ^= 1;
    ((u8 *)source_ip)[2] ^= 1; 

    ip_header->check = 0;
    ip_header->check = ip_fast_csum((unsigned char *) ip_header, ip_header->ihl); //Recalculamos en cs

    struct net_device * destination_device = net_devices[source_device == net_devices[0] ? 1 : 0];
    struct my_net * destination_my_net = netdev_priv(destination_device);
    
    struct my_net_packet * my_net_packet = net_get_packet_from_packet_pool(source_device);
    my_net_packet->length = length_packet;
    memcpy(my_net_packet->data, data_packet, length_packet);
    net_enqueue_packet_read_queue(destination_my_net->device, my_net_packet);

    if(destination_my_net->read_interruptions_enabled){
        destination_my_net->status |= MY_NET_PENDING_READ_PACKET_INTERRUPT;
    }

    struct my_net * my_net_source = netdev_priv(source_device);
    my_net_source->tx_packetlen = length_packet;
    my_net_source->tx_packetdata = data_packet;
    my_net_source->status |= MY_NET_PENDING_TRANSMISSION;
    net_interrupt_handler(0, source_device, NULL);
}

void net_regular_interrupt_handler(int irq, void * dev_id, struct pt_regs * regs) {
    struct net_device * net_device = (struct net_device *) dev_id;   
    struct my_net * my_net = netdev_priv(net_device);
    struct my_net_packet * my_net_packet_recieved = NULL;

    spin_lock(&my_net->lock);

    int status = my_net->status;
    my_net->status = 0;

    if(status & MY_NET_PENDING_READ_PACKET_INTERRUPT){
        my_net_packet_recieved = my_net->read_queue;
        my_net->read_queue = my_net_packet_recieved->next;

        net_read(net_device, my_net_packet_recieved);
    }
    if(status & MY_NET_PENDING_TRANSMISSION){
        my_net->stats.tx_packets++;
	my_net->stats.tx_bytes += my_net->tx_packetlen;

        dev_kfree_skb(my_net->sk_buff);
    }

    spin_unlock(&my_net->lock);

    if(my_net_packet_recieved){
        net_release_packet(my_net_packet_recieved);        
    }
}

void net_napi_interrupt_handler(int irq, void * dev_id, struct pt_regs * regs) {
    struct net_device * net_device = (struct net_device *) dev_id;   
    struct my_net * my_net = netdev_priv(net_device);

    spin_lock(&my_net->lock);
    
    int status = my_net->status;
    my_net->status = 0;

    if (status & MY_NET_PENDING_READ_PACKET_INTERRUPT) {
        net_read_enable_interrupts(net_device, 0); //Disable interrupciones
	napi_schedule(&my_net->napi);
    }
    if(status & MY_NET_PENDING_TRANSMISSION){
        my_net->stats.tx_packets++;
	my_net->stats.tx_bytes += my_net->tx_packetlen;
	dev_kfree_skb(my_net->sk_buff);
    }

    spin_unlock(&my_net->lock);
}

int net_poll(struct napi_struct * napi, int budget) { //Utilizado por napi cada vez que hay paquetes para leer
    struct my_net * my_net = container_of(napi, struct my_net, napi);
    struct net_device * device = my_net->device;
    int n_packets = 0;

    while(n_packets < budget && my_net->read_queue) {
        struct my_net_packet * my_net_packet = net_dequeue_packet_from_read_queue(device);
        struct sk_buff * sk_buff = dev_alloc_skb(my_net_packet->length + 2);

        skb_reserve(sk_buff, 2);
        
        memcpy(skb_put(sk_buff, my_net_packet->length), my_net_packet->data, my_net_packet->length);
	sk_buff->dev = device;
	sk_buff->protocol = eth_type_trans(sk_buff, device);
	sk_buff->ip_summed = CHECKSUM_UNNECESSARY;

        netif_receive_skb(sk_buff);
    
        n_packets++;
        my_net->stats.rx_packets++;
	my_net->stats.rx_bytes += my_net_packet->length;
	net_release_packet(my_net_packet);
    }

    if(n_packets < budget){ //Todos los paquetes han sido leidos, reestablecemos las interrupciones
        unsigned long flags;

        spin_lock_irqsave(&my_net->lock, flags);
        
        napi_complete_done(napi, n_packets);
        net_read_enable_interrupts(device, 1);
        
        spin_unlock_irqrestore(&my_net->lock, flags);
    }

    return n_packets;
}

int net_create_header_arp(struct sk_buff * sk_buff, struct net_device * device, unsigned short type, const void * destination_address, 
	    const void *source_address, unsigned length) {

    struct ethhdr * ethernet_header = (struct ethhdr *) skb_push(sk_buff, ETH_HLEN);

    ethernet_header->h_proto = htons(type);
    memcpy(ethernet_header->h_source, source_address ? source_address : device->dev-addr, device->addr_len);
    memcpy(ethernet_header->h_dest, destination_address ? destination_address : device->dev-addr, device->addr_len);

    ethernet_header->h_dest[ETH_ALEN - 1] ^= 0x01;

    return device->hard_header_len;
}

struct net_device_stats * net_get_stats(struct net_device * device) {
    struct my_net * my_net = netdev_priv(device);

    return &my_net->stats;
}

int net_open(struct net_device * device) {
    memcpy(device->dev_addr, "\0SNUL0", ETH_ALEN); //Asignamos una mac false
    if(device == net_devices[1]){
        unsigned char * mac = (unsigned char *) device->dev_addr; //Is declared as const
        mac[ETH_ALEN - 1]++;
    }

    if(MY_NET_USE_NAPI){
        struct my_net * my_net = netdev_priv(device);
	napi_enable(&my_net->napi);
    }

    netif_start_queue(device);

    return 0;
}

int net_release(struct net_device * device) {
    netif_stop_queue(device);

    if(MY_NET_USE_NAPI){
        struct my_net * my_net = netdev_priv(device);
	napi_enable(&my_net->napi);
    }

    return 0;
}

void net_tx_timeout(struct net_device * device, unsigned int txqueue) {
    struct my_net * my_net = netdev_priv(device);

    my_net->status |= MY_NET_PENDING_TRANSMISSION;
    net_interrupt_handler(0, device, NULL);
    my_net->stats.tx_errors++;

    spin_lock(&my_net->lock);
    net_destroy_packet_pool(device);
    net_setup_packet_pool(device);
    spin_unlock(&my_net->lock);
    netif_wake_queue(device);
}

static struct my_net_packet * net_get_packet_from_packet_pool(struct net_device * device) {
    struct my_net * my_net = netdev_priv(device);
    unsigned long flags;

    spin_lock_irqsave(&my_net->lock, flags);
    
    struct my_net_packet * my_net_packet = my_net->packet_pool;
    if(my_net_packet == NULL){
        spin_unlock_irqrestore(&my_net->lock, flags);
        return my_net_packet;
    }
    my_net->packet_pool = my_net_packet->next;
    if(my_net->packet_pool == NULL){
        netif_stop_queue(device);
    }

    spin_unlock_irqrestore(&my_net->lock, flags);
    
    return my_net_packet;
}

static struct my_net_packet * net_dequeue_packet_from_read_queue(struct net_device * device) {
	struct my_net * my_net = netdev_priv(device);
    unsigned long flags;

    spin_lock_irqsave(&my_net->lock, flags);

    struct my_net_packet * dequeued_packet = my_net->read_queue;
    if (dequeued_packet != NULL) {
	my_net->read_queue = dequeued_packet->next;
    }

    spin_unlock_irqrestore(&my_net->lock, flags);

    return dequeued_packet;
}

static void net_enqueue_packet_read_queue(struct net_device * device, struct my_net_packet * my_net_packet) {
    struct my_net * my_net = netdev_priv(device);
    unsigned long flags;

    spin_lock_irqsave(&my_net->lock, flags);
    my_net_packet->next = my_net->read_queue;
    my_net->read_queue = my_net_packet;
    spin_unlock_irqrestore(&my_net->lock, flags);
}

static void net_setup_packet_pool(struct net_device * device) {
    struct my_net * my_net = netdev_priv(device);
    my_net->packet_pool = NULL;

    for(int i = 0; i < MY_NET_PACKET_POOL_SIZE; i++){
        struct my_net_packet * my_net_packet = kmalloc(sizeof(struct my_net_packet), GFP_KERNEL);       
        my_net_packet->device = device;
        my_net_packet->next = my_net->packet_pool;
        my_net->packet_pool = my_net_packet;
    }
}

static void net_destroy_packet_pool(struct net_device * device) {
    struct my_net * my_net = netdev_priv(device);
	struct my_net_packet * packet;

    while(packet != NULL) {
        packet = my_net->packet_pool;
        my_net->packet_pool = packet->next;

        kfree(packet);
    }
}

static void net_read_enable_interrupts(struct net_device * device, int enable) {
    struct my_net * my_net = netdev_priv(device);
    my_net->read_interruptions_enabled = enable;
}

static void net_release_packet(struct my_net_packet * my_net_packet) {
    struct my_net * my_net = netdev_priv(my_net_packet->device);
    unsigned long flags;

    spin_lock_irqsave(&my_net->lock, flags);
    my_net_packet->next = my_net->packet_pool;
    my_net->packet_pool = my_net_packet;
    spin_unlock_irqrestore(&my_net->lock, flags);

    if (netif_queue_stopped(my_net_packet->device) && my_net_packet->next == NULL){
	netif_wake_queue(my_net_packet->device);
    }
}

void net_init_device(struct net_device * device) {
    ether_setup(device);
    
    device->watchdog_timeo = MY_NET_TIMEOUT; //Timout de cada paquete en Jiffies
    device->netdev_ops = &my_net_netdev_ops;
    device->header_ops = &my_net_header_ops;
    device->flags |= IFF_NOARP; //Quitamos ARP. No hay direcciones MAC, son dispositivos falsos
    device->features |= NETIF_F_HW_CSUM; //El checksum lo hace el hardware no la CPU  

    struct my_net * my_net = netdev_priv(device);
    memset(my_net, 0, sizeof(struct my_net));

    if(MY_NET_USE_NAPI){
        netif_napi_add(device, &my_net->napi, net_poll, 2);
    }

    spin_lock_init(&my_net->lock);
    my_net->device = device;

    net_read_enable_interrupts(device, 1);		/* permitimos interruptciones */
    net_setup_packet_pool(device);
}

static int __init net_init(void) {
    net_devices[0] = alloc_netdev(sizeof(struct my_net), "sn%d", NET_NAME_UNKNOWN, net_init_device);
    net_devices[1] = alloc_netdev(sizeof(struct my_net), "sn%d", NET_NAME_UNKNOWN, net_init_device);
    
    register_netdev(net_devices[0]);
    register_netdev(net_devices[1]);

    net_interrupt_handler = MY_NET_USE_NAPI ? net_napi_interrupt_handler : net_regular_interrupt_handler;

    return 0;
}

static void __exit net_exit(void) {
    unregister_netdev(net_devices[0]); //Borra la interfaz
    unregister_netdev(net_devices[1]);

    net_destroy_packet_pool(net_devices[0]);
    net_destroy_packet_pool(net_devices[1]);

    free_netdev(net_devices[0]); 
    free_netdev(net_devices[1]);
}

module_init(net_init);
module_exit(net_exit);

MODULE_LICENSE("Dual BSD/GPL");
MODULE_DESCRIPTION("Module");
MODULE_AUTHOR("Jaime");
