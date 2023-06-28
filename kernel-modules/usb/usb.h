#ifndef _USB_
#define _USB_

#include <linux/rwsem.h>
#include <linux/module.h>
#include <linux/sched.h>
#include <linux/types.h>
#include <linux/fs.h>
#include <linux/cdev.h>
#include <linux/mutex.h>
#include <linux/kernel.h>
#include <asm/uaccess.h>
#include <linux/ioctl.h>
#include <linux/capability.h>
#include <linux/fs.h>
#include <linux/poll.h>
#include <linux/usb.h>
#include <linux/kobject.h>

#define USB_PRODUCT_ID 0xfff0
#define USB_VENDOR_ID 0xfff0

struct my_usb_driver {
    struct usb_device * device;
    struct usb_interface * interface;
    unsigned char * buffer;
    size_t buffer_size;

    __u8 out_endpointAdr;
    __u8 in_endpointAdr;

    struct kobject kobj;
};

int usb_probe(struct usb_interface * interface, const struct usb_device_id * id);
void usb_disconnect(struct usb_interface * interface);

int usb_file_open(struct inode * inode, struct file * file);
ssize_t usb_file_read(struct file * file, char __user * buffer, size_t count, loff_t *f_pos);
ssize_t usb_file_write(struct file * file, const char __user * buffer, size_t count, loff_t *f_pos);
int usb_file_release(struct inode * inode, struct file * file);

static void usb_file_written_callback(struct urb * urb);

inline bool is_endpoint_in(struct usb_endpoint_descriptor * endpoint, struct my_usb_driver * driver);
inline bool is_endpoint_out(struct usb_endpoint_descriptor * endpoint, struct my_usb_driver * driver);

#endif