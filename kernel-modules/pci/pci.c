#include "pci.h"

static struct pci_device_id pci_ids[] = {
    {PCI_DEVICE(PCI_VENDOR_ID_INTEL, PCI_DEVICE_ID_INTEL_82801AA_3)},
    {0}
};
MODULE_DEVICE_TABLE(pci, pci_ids);

static struct pci_driver pci_driver = {
    .name = "my pci",
    .id_table = pci_ids,
    .probe = pci_probe,
	.remove = pci_remove,
};

static int pci_probe(struct pci_dev * device, const struct pci_device_id *id) {
    if(pci_enable_device(dev)) {
		return -ENODEV;
	}

    return 0;
}

static void pci_remove(struct pci_dev *dev) {
    //Nada necesario
}

static int __init pci_init(void) {
    return pci_register_driver();
}


static void __exit pci_exit(void) {
    pci_unregister_driver(&pci_driver);
}

MODULE_LICENSE("GPL");

module_init(pci_init);
module_exit(pci_exit);