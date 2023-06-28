#ifndef _PCI_
#define _PCI_

#include <linux/kernel.h>
#include <linux/module.h>
#include <linux/pci.h>
#include <linux/init.h>

static int pci_probe(struct pci_dev *dev, const struct pci_device_id *id);
static void pci_remove(struct pci_dev *dev);

static int __init pci_init(void);
static void __exit pci_exit(void);

#endif