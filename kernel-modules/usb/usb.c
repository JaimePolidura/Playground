#include "usb.h"

static struct usb_device_id usb_devices [] = {
    {USB_DEVICE(USB_VENDOR_ID, USB_PRODUCT_ID)},
};

static struct usb_driver usb_driver_props = {
    .name = "myUsb",
    .id_table = usb_devices,
    .probe = usb_probe,
    .disconnect = usb_disconnect,
};

static struct file_operations usb_driver_file_ops = {
    .owner = THIS_MODULE,
    .open = usb_file_open,
    .read = usb_file_read,
    .write = usb_file_write,
    .release = usb_file_release,
};

static struct usb_class_driver usb_driver_file = {
    .name = "usb/skell%d",
    .fops = &usb_driver_file_ops,
    .minor_base = 192,
};

int usb_file_open(struct inode * inode, struct file * file) {
    int subminor = iminor(inode);
    struct usb_interface * interface = usb_find_interface(&usb_driver_props, subminor);

    file->private_data = interface;

    return 0;
}

ssize_t usb_file_read(struct file * file, char __user * buffer, size_t count, loff_t *f_pos) {
    struct my_usb_driver * driver = file->private_data;
    
    int result = usb_bulk_msg(driver->device, 
        usb_rcvbulkpipe(driver->device, driver->in_endpointAdr),
        driver->in_endpointAdr,
        min(driver->buffer_size, count),
        (int *) &count,
        HZ * 5); //5 Segundos
    
    if(result) {
        return 0;
    }

    copy_to_user(buffer, driver->in_endpointAdr, count);

    return count;
}

ssize_t usb_file_write(struct file * file, const char __user * buffer_user, size_t count, loff_t *f_pos) {
    struct my_usb_driver * driver = file->private_data;

    struct urb * urb = usb_alloc_urb(0, GFP_KERNEL);
    char * buffer_usb = usb_alloc_coherent(driver->device, driver->buffer_size, GFP_KERNEL, &urb->transfer_dma);
    copy_from_user(buffer_usb, buffer_user, count);

    usb_fill_bulk_urb(urb, 
        driver->device, 
        usb_sndbulkpipe(driver->device, driver->out_endpointAdr), 
        buffer_usb, 
        count, 
        usb_file_written_callback,
        driver);
    urb->transfer_flags |= URB_NO_TRANSFER_DMA_MAP;

    usb_submit_urb(urb, GFP_KERNEL);

    usb_free_urb(urb);

    return count;
}

static void usb_file_written_callback(struct urb * urb) {
    struct my_usb_driver * driver = urb->context;

    usb_free_coherent(urb->dev, urb->transfer_buffer_length, urb->transfer_buffer, urb->transfer_dma);
}

int usb_file_release(struct inode * inode, struct file * file) {
    return 0;
}

int usb_probe(struct usb_interface * interface, const struct usb_device_id * id) {
    struct my_usb_driver * my_driver = kzalloc(sizeof(struct my_usb_driver), GFP_KERNEL);
    struct usb_host_interface * interface_desc = interface->cur_altsetting;

    my_driver->device = usb_get_dev(interface_to_usbdev(interface));
    my_driver->interface = interface;
    
    for(int i = 0; i < interface_desc->desc.bNumEndpoints; i++){
        struct usb_endpoint_descriptor * endpoint = &interface_desc->endpoint[i].desc;
        
        if(is_endpoint_in(endpoint, my_driver)) {
			my_driver->buffer_size = endpoint->wMaxPacketSize;
			my_driver->in_endpointAdr = endpoint->bEndpointAddress;
			my_driver->buffer = kmalloc(endpoint->wMaxPacketSize, GFP_KERNEL);
        }

        if(is_endpoint_out(endpoint, my_driver)) {
            my_driver->out_endpointAdr = endpoint->bEndpointAddress;
        }
    }

    usb_set_intfdata(interface, my_driver);
    usb_register_dev(interface, &usb_driver_file);

    return 0;
}

inline bool is_endpoint_in(struct usb_endpoint_descriptor * endpoint, struct my_usb_driver * driver) {
    return !driver->in_endpointAdr &&
		    (endpoint->bEndpointAddress & USB_DIR_IN) &&
		    ((endpoint->bmAttributes & USB_ENDPOINT_XFERTYPE_MASK) == USB_ENDPOINT_XFER_BULK);
}

inline bool is_endpoint_out(struct usb_endpoint_descriptor * endpoint, struct my_usb_driver * driver) {
    return !driver->out_endpointAdr &&
		    !(endpoint->bEndpointAddress & USB_DIR_IN) &&
		    ((endpoint->bmAttributes & USB_ENDPOINT_XFERTYPE_MASK) == USB_ENDPOINT_XFER_BULK);
}

void usb_disconnect(struct usb_interface * interface) {
    usb_set_intfdata(interface, NULL);
    usb_deregister_dev(interface, &usb_driver_file);
}

static int __init usb_init(void) {
    usb_register(&usb_driver_props);

	return 0;
}

static void __exit usb_exit(void) {
    usb_deregister(&usb_driver_props);
}

module_init(usb_init);
module_exit(usb_exit);

MODULE_LICENSE("Dual BSD/GPL");
MODULE_DESCRIPTION("Module");
MODULE_AUTHOR("Jaime");
