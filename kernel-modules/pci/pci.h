#include <linux/kernel.h>
#include <linux/module.h>
#include <linux/pci.h>
#include <linux/init.h>

static int probe(struct pci_dev *dev, const struct pci_device_id *id);
static void remove(struct pci_dev *dev);

static int __init pci_init(void);
static void __exit pci_exit(void);