#include "block.h"

static struct my_block * my_block;

static const struct block_device_operations fops = {
	.owner = THIS_MODULE,
	.open = block_open,
	.release = block_release,
	.ioctl = block_ioctl,
	.submit_bio = block_submit_bio,
};

static int block_open(struct block_device * bdev, fmode_t mode) {
    struct my_block * my_block = bdev->bd_disk->private_data;
    
    return 0;
}

static int __init block_init(void) {
    register_blkdev(MY_BLOCK_MAJOR, "myblock");

    my_block = kzalloc(sizeof(struct my_block), GFP_KERNEL);
    spin_lock_init(&my_block->lock);

    my_block->size = MY_BLOCK_N_SECTORS * MY_BLOCK_SECTOR_SIZE;
    my_block->data = vmalloc(my_block->size);
    my_block->queue = blk_mq_init_queue(&my_block->tag_set);
    my_block->gendisk = blk_alloc_disk(NUMA_NO_NODE);
    my_block->gendisk->major = MY_BLOCK_MAJOR;
    my_block->gendisk->first_minor = MY_BLOCK_MINORS * 1; 
    my_block->gendisk->fops = &block_ops;
    my_block->gendisk->queue = my_block->queue;
    my_block->gendisk->private_data = my_block;
	sprintf(my_block->gendisk->disk_name, "myblock");
    set_capacity(my_block->gendisk, MY_BLOCK_N_SECTORS * (MY_BLOCK_SECTOR_SIZE / KERNEL_SECTOR_SIZE));

	blk_queue_physical_block_size(my_block->queue, MY_BLOCK_SECTOR_SIZE);
	blk_queue_logical_block_size(my_block->queue, MY_BLOCK_SECTOR_SIZE);
	blk_queue_max_hw_sectors(my_block->queue, BLK_DEF_MAX_SECTORS);
	blk_queue_flag_set(QUEUE_FLAG_NOMERGES, my_block->queue);

    add_disk(my_block->gendisk);

    return 0;
}

static void __exit block_exit(void) {
    del_gendisk(my_block->gendisk);
    put_disk(my_block->gendisk);
    vfree(my_block->data);
    blk_cleanup_queue(my_block->queue);
    unregister_blkdev(MY_BLOCK_MAJOR, "myblock");
    kfree(my_block);
}

module_init(block_init);
module_exit(block_exit);

MODULE_LICENSE("Dual BSD/GPL");
MODULE_DESCRIPTION("Module");
MODULE_AUTHOR("Jaime");
