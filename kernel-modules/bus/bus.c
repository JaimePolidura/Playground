#include "bus.h"


struct bus_type my_bus_native = {
    .name = "mybus",
    .uevent = bus_uevent,
    .match = bus_match,
};

struct device my_device_native = {
    .release = bus_release,
};

static ssize_t show_version(struct device_driver *driver, char * buffer) {
	struct my_bus_driver * my_bus_driver = to_my_bus_driver(driver);
    
	sprintf(buffer, "%s\n", my_bus_driver->version);
	return strlen(buffer);
}

int register_my_bus_driver(struct my_bus_driver * my_bus_driver) {
    my_bus_driver->driver.bus = &my_bus_native;
    
    if(driver_register(&my_bus_driver->driver)){
        return -1; 
    }

	my_bus_driver->version_attr.attr.name = "version";
	my_bus_driver->version_attr.attr.mode = S_IRUGO;
	my_bus_driver->version_attr.show = show_version;
	my_bus_driver->version_attr.store = NULL;    

    return driver_create_file(&my_bus_driver->driver, &my_bus_driver->version_attr);
}

int unregister_my_bus_driver(struct my_bus_driver * my_bus_driver) {
    driver_unregister(&my_bus_driver->driver);
    return 0;
}

int register_my_bus_device(struct my_bus_device * my_bus_device) {
    my_bus_device->device.parent = &my_device_native;
    my_bus_device->device.release = bus_release;
    dev_set_name(&my_bus_device->device, my_bus_device->name);

    device_register(&my_bus_device->device);

    return 0;
}

int unregister_my_bus_device(struct my_bus_device * my_bus_device) {
    device_unregister(&my_bus_device->device);
    return 0;
}

static void bus_release(struct device *dev) {
    printk(KERN_ALERT "Released bus\n");
}

static int bus_uevent(struct device * device, struct kobj_uevent_env * env) {
    if (add_uevent_var(env, "LDDBUS_VERSION=%s", bus_version))
		return -ENOMEM;

	return 0;
}

static int bus_match (struct device * device, struct device_driver * driver) {
    return !strncmp(dev_name(device), driver->name, strlen(driver->name));
}

static int __init bus_init(void) {
    bus_register(&my_bus_native);

    dev_set_name(&my_device_native, "ldd-device-0");
    device_register(&my_device_native);

	return 0;
}

static void __exit bus_exit(void) {
    bus_unregister(&my_bus_native);
}

module_init(bus_init);
module_exit(bus_exit);

EXPORT_SYMBOL(register_my_bus_driver);
EXPORT_SYMBOL(unregister_my_bus_driver);
EXPORT_SYMBOL(register_my_bus_device);
EXPORT_SYMBOL(unregister_my_bus_device);

MODULE_LICENSE("Dual BSD/GPL");
MODULE_DESCRIPTION("Module");
MODULE_AUTHOR("Jaime");
