#include "tty.h"

static struct tty_driver * my_tty_driver;
static struct tty_port my_tty_port[MY_TTY_MINORS];
static struct my_tty * my_tty_table[MY_TTY_MINORS];

static struct tty_operations my_tty_ops = {
    .open = tty_open,
    .close = tty_close,
    .write = tty_write,
    .write_room = tty_write_room,
    .set_termios = tiny_set_termios,
};

struct termios my_tty_std_termios = {
    .c_iflag = ICRNL | IXON,
    .c_oflag = OPOST | ONLCR,
    .c_cflag = B38400 | CS8 | CREAD | HUPCL,
    .c_lflag = ISIG | ICANON | ECHO | ECHOE | ECHOK | ECHOCTL | ECHOKE | IEXTEN,
    .c_cc = INIT_C_CC
};

//Usuario --> tty core -> mi driver write() --> hardware. Los datos se suelen poner en un buffer
int tiny_write(struct tty_struct * tty, const unsigned char *buffer, int count) {
    //Enviar datos al hardware
    return count;
}

static void tiny_set_termios(struct tty_struct * tty, struct ktermios * old_termios)  {
    if(!tty || !old_termios){
        return;
    }
    if(old_termios->c_flag == tty->termios.c_iflag && 
        (MY_TTY_RELEVANT_IFLAG(tty->termios.c_iflag) == MY_TTY_RELEVANT_IFLAG(old_termios->c_iflag))){
            return; //No ha cambiado nada
    }
}

//Se llamda cada vez que tty core quiere saber cuanto espacio disponible hay en el buffer
unsigned int tty_write_room(struct tty_struct * tty) {
    return 255;
}

static void tty_timeout(struct timer_list * timer_list) {
    struct my_tty * my_tty = from_timer(my_tty, timer_list, timer);

    char data[1] = {MY_TTY_DATA_CHARACTER};
    int port = my_tty->port;

    for(int i = 0; i < 1; i++){
        if(!tty_buffer_request_room(port, 1)){ //Buffer lleno
            tty_flip_buffer_push(port); //Mandar al usuario
        }
        tty_insert_flip_char(port, data[i], TTY_NORMAL); //Flip buffer
    }

    tty_flip_buffer_push(port);

    my_tty->timer.expires = jiffies + MY_TTY_TIMEOUT;
	add_timer(&my_tty->timer);
}


int tty_open(struct tty_struct * tty, struct file * file) {
    int index = tty->index;
    struct my_tty * my_tty_entry_in_table = my_tty_table[index];

    if(my_tty_entry_in_table == NULL){
        my_tty_entry_in_table = kmalloc(sizeof(struct my_tty), GFP_KERNEL | __GFP_ZERO);
        mutex_init(&my_tty_entry_in_table->mutex);

        my_tty_table[index] = my_tty_entry_in_table;
    }

    mutex_lock(&my_tty_entry_in_table->mutex);

    tty->driver_data = my_tty_entry_in_table;
    my_tty_entry_in_table->tty = tty;
    my_tty_entry_in_table->open_count++;

    if(my_tty_entry_in_table->open_count == 1){
        timer_setup(my_tty_entry_in_table->timer, tty_timeout, 0);
        my_tty_entry_in_table->timer.expires = jiffies + MY_TTY_TIMEOUT;
        add_timer(&my_tty_entry_in_table->timer);
    }

    mutex_unlock(&my_tty_entry_in_table->mutex);

    return 0;
}

void tty_close(struct tty_struct * tty, struct file * file) {
    struct my_tty * my_tty = tty->driver_data;
    
    mutex_lock(&my_tty->mutex);

    my_tty->open_count--;

    if(my_tty->open_count <= 0){
        del_timer(&my_tty->timer);
    }

    mutex_unlock(&my_tty->mutex);
}

static int __init tty_init(void) {
    my_tty_driver = alloc_tty_driver(MY_TTY_MINORS);

    my_tty_driver->owner = THIS_MODULE;
    my_tty_driver->driver_name = "my_tty";
    my_tty_driver->name = "my_tty";
    my_tty_driver->major = MY_TTY_MAJOR;
    my_tty_driver->type = TTY_DRIVER_TYPE_SERIAL;
    my_tty_driver->subtype = SERIAL_TYPE_NORMAL;
    my_tty_driver->flags = TTY_DRIVER_REAL_RAW | TTY_DRIVER_DYNAMIC_DEV,
	my_tty_driver->init_termios = my_tty_std_termios;
	my_tty_driver->init_termios.c_cflag = B9600 | CS8 | CREAD | HUPCL | CLOCAL;
	tty_set_operations(my_tty_driver, &my_tty_ops); 
    
    for(int i = 0; i < MY_TTY_MINORS; i++){
        tty_port_init(my_tty_port + i);
        tty_port_link_device(my_tty_port + i, my_tty_driver, i);
    }

    tty_register_driver(my_tty_driver);

    for(int i = 0; i < MY_TTY_MINORS; i++){
        tty_register_device(my_tty_driver, i, NULL);
    }

    return 0;
}

static void __exit tty_exit(void) {
	for (int i = 0; i < MY_TTY_MINORS; ++i) {
		tty_unregister_device(my_tty_driver, i);
		tty_port_destroy(my_tty_port + i);
	}
    
    for (int i = 0; i < MY_TTY_MINORS; ++i) {
        tty_unregister_device(my_tty_driver, i);
    }
    
    tty_unregister_driver(my_tty_driver);
}

module_init(tty_init);
module_exit(tty_exit);

MODULE_LICENSE("Dual BSD/GPL");
MODULE_DESCRIPTION("Module");
MODULE_AUTHOR("Jaime");