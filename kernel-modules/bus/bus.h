#ifndef _BUS_
#define _BUS_

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

static char *bus_version = "$Revision: 1.9 $";

struct my_bus_driver {
	char *version;
	struct module *module;
	struct device_driver driver;
	struct driver_attribute version_attr;
};

#define to_my_bus_driver(driver) container_of(driver, struct my_bus_driver, driver);

struct my_bus_device {
    char * name;
    struct device device;
    struct my_bus_driver driver;
};

#define to_my_bus_device(device) container_of(device, struct my_bus_device, device);

int register_my_bus_driver(struct my_bus_driver * my_bus_driver);
int unregister_my_bus_driver(struct my_bus_driver * my_bus_driver);

int register_my_bus_device(struct my_bus_device * my_bus_device);
int unregister_my_bus_device(struct my_bus_device * my_bus_device);

static void bus_release(struct device *dev);
static int bus_match (struct device * device, struct device_driver * driver);
static int bus_uevent(struct device * device, struct kobj_uevent_env * env);

static int __init bus_init(void);
static void __exit bus_exit(void);

#endif