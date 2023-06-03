#include <linux/module.h>
#include <linux/sched.h>

static int __init hello_init(void){
        printk(KERN_ALERT "Hello world, The current process is %i\n", current->pid);
	return 0;
}

static void __exit hello_exit(void) {
        printk(KERN_ALERT "Bye bye\n");
}

module_init(hello_init);
module_exit(hello_exit);

MODULE_LICENSE("Dual BSD/GPL");
MODULE_DESCRIPTION("Module");
MODULE_AUTHOR("Jaime");
