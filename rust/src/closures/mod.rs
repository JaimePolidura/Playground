use std::thread;

fn closures() {
    let list = vec![1, 2, 3];
    println!("Before defining closure: {:?}", list);
    let only_borrows = || println!("From closure: {:?}", list);

    let mut list = vec![1, 2, 3];
    let mut borrows_mutably = || list.push(7);
    borrows_mutably();

    //Ownership
    thread::spawn(move || {
        println!("From thread: {:?}", list)
    }).join().unwrap();
}

pub fn apply<F>(value: u32, f: F) -> u32
where
    F: FnOnce() -> u32
{
    return 1;
}