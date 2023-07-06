#include "block.h"

static struct my_block * my_block;

static const struct block_device_operations block_ops = {
	.owner = THIS_MODULE,
	.open = block_open,
	.release = block_release,
	.ioctl = block_ioctl,
	.submit_bio = block_submit_bio,
};

static struct blk_mq_ops my_block_queue_ops = {
    .queue_rq = block_queue_request,
};

blk_status_t block_queue_request(struct blk_mq_hw_ctx * blk_mq_hw_ctx, const struct blk_mq_queue_data * blk_mq_queue_data) {
	struct request * request = blk_mq_queue_data->rq;
	struct my_block * my_block = request->q->queuedata;
    sector_t position_sector = blk_rq_pos(request);
    struct req_iterator req_iterator;
    struct bio_vec bio_vec;

    blk_mq_start_request(request);

    if(blk_rq_is_passthrough(request)) {
        blk_mq_end_request(request, BLK_STS_IOERR);
        return BLK_STS_IOERR;
    }

    rq_for_each_segment(bio_vec, request, req_iterator) {
	    size_t sector_num = blk_rq_cur_sectors(request);
        void * buffer = page_address(bio_vec.bv_page) + bio_vec.bv_offset;;
        
        block_transfer(my_block, position_sector, sector_num, buffer, rq_data_dir(request) == WRITE);

        position_sector += sector_num;
    }

    blk_mq_end_request(request, BLK_STS_OK);

    return BLK_STS_OK;
}

void block_submit_bio(struct bio * bio) {
    struct my_block * my_block = bio->bi_bdev->bd_disk->private_data; 
    loff_t position = bio->bi_iter.bi_sector << SECTOR_SHIFT;
	struct bvec_iter bvec_iterator;
    struct bio_vec bvec;

    bio_for_each_segment(bvec, bio, bvec_iterator) {
        unsigned int length = bvec.bv_len;
        void * buffer = page_address(bvec.bv_page) + bvec.bv_offset;

        block_transfer(my_block, length, position, buffer, bio_data_dir(bio));
    }
    
    bio_endio(bio);
}

void block_transfer(struct my_block * my_block, unsigned long sector, unsigned long sector_number, char * buffer, int write) {
	unsigned long nbytes = sector_number * KERNEL_SECTOR_SIZE;
    unsigned long offset = sector * KERNEL_SECTOR_SIZE;

    if (write)
		memcpy(my_block->data + offset, buffer, nbytes);
	else
		memcpy(my_block, my_block->data + offset, nbytes);
}

int block_ioctl (struct block_device * bdev, fmode_t mode, unsigned int cmd, unsigned long arg) {
    struct my_block * my_block = bdev->bd_disk->private_data;

    switch(cmd) {
        case HDIO_GETGEO:
            long size = my_block->size * (MY_BLOCK_SECTOR_SIZE / KERNEL_SECTOR_SIZE);

            struct hd_geometry geo = {
                .cylinders = (size & ~0x3F) >> 6,
                .heads = 4,
                .sectors = 16,
                .start = 4
            };

            copy_to_user((void __user *) arg, &geo, sizeof(geo));

            return 0;
    }

    return -ENOTTY;
}

static int block_open(struct block_device * bdev, fmode_t mode) {    
    del_timer_sync(&my_block->timer);

	spin_lock(&my_block->lock);
    my_block->n_users++;
    spin_unlock(&my_block->lock);

    return 0;
}

static void block_release(struct gendisk * disk, fmode_t mode) {
    struct my_block * my_block = disk->private_data;

    spin_lock(&my_block->lock);
    my_block->n_users--;
    if(my_block->n_users == 0){
        my_block->timer.expires = jiffies + MY_BLOCK_INVALIDATE_DELAY;
        add_timer(&my_block->timer);
    }
    spin_unlock(&my_block->lock);
}

void block_timer_timeout(struct timer_list * timer) {
    struct my_block *my_block = from_timer(my_block, timer, timer);

    spin_lock(&my_block->lock);
    my_block->media_changed = 1;
    spin_unlock(&my_block->lock);
}

static int __init block_init(void) {
    register_blkdev(MY_BLOCK_MAJOR, "myblock");

    my_block = kzalloc(sizeof(struct my_block), GFP_KERNEL);
    spin_lock_init(&my_block->lock);
    blk_mq_alloc_sq_tag_set(&my_block->tag_set, &my_block_queue_ops, 128, 0);
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
    timer_setup(&my_block->timer, block_timer_timeout, 0);

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
    //Clean queue
    unregister_blkdev(MY_BLOCK_MAJOR, "myblock");
    kfree(my_block);
}

module_init(block_init);
module_exit(block_exit);

MODULE_LICENSE("Dual BSD/GPL");
MODULE_DESCRIPTION("Module");
MODULE_AUTHOR("Jaime");
