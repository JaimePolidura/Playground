#include <linux/init.h>
#include <linux/module.h>

static int __init hello_init(void){
        printk(KERN_ALERT, "Hello world\n");
	return 0;
}

static void __exit hello_exit(void) {
        printk(KERN_ALERT, "Bye bye\n");
	return;
}

module_init(hello_init);
module_exit(hello_exit);

MODULE_LICENSE("Dual BSD/GPL");
MODULE_DESCRIPTION("Module");
MODULE_AUTHOR("Jaime");
