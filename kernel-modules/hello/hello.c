#include <linux/module.h>
#include <linux/sched.h>
#include <asm/msr.h>
#include <linux/timex.h>
#include <asm/page.h>

static int __init hello_init(void){
        unsigned long long tsc;
        tsc = rdtsc();
        
        printk(KERN_ALERT "Hello world, The current process is %i and %i CPU cycles. The current page size %i\n", current->pid, tsc, PAGE_SIZE);

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
