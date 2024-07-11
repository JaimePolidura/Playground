use std::num::Add;
use std::ops::Add;

fn unsafe_call() {
    let mut num = 5;
    let ptr_1: * const i32 = &num as * const i32;
    let ptr_2: * mut i32 = &mut num as * mut i32;

    unsafe {
        println!("r1 is: {}", * ptr_1);
        println!("r2 is: {}", * ptr_2);

        mierdon();
    }

    // mierdon();
}

unsafe fn mierdon() {}

fn split_at_mut(
    values: &mut [i32],
    mid: usize,
) -> (&mut [i32], &mut [i32]) {
    let len = values.len();
    assert!(mid <= len);
    let ptr: * mut i32 = values.as_mut_ptr();

    unsafe {
        (
            std::slice::from_raw_parts_mut(ptr, mid),
            std::slice::from_raw_parts_mut(ptr.add(mid), len - mid),
        )
    }
}

extern "C" {
    fn abs(input: i32) -> i32;
}

fn external() {
    unsafe {
        println!(
            "Absolute value of -3 according to C: {}",
            abs(-3)
        );
    }
}

//Se puede llamar en C, No mangle hace que el compilador no cambie el nombre de la function
#[no_mangle]
pub extern "C" fn call_from_c() {
    println!("Just called a Rust function from C!");
}

//Al ser static, la direccion de memoria de la variable global siempre es la misma
static mut COUNTER: u32 = 0;

fn add_to_count(inc: u32) {
    unsafe {
        COUNTER += inc;
    }
}
