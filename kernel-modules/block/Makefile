CCFLAGS := -Wunused-variable -Wunused-result
CONFIG_MODULE_SIG=n 
obj-m := block.o

KDIR := /lib/modules/$(shell uname -r)/build
PWD := $(shell pwd)

all:
	make -C $(KDIR) M=$(PWD) CFLAGS_EXTRA='$(CCFLAGS)' modules

clean:
	make -C $(KDIR) M=$(PWD) clean
