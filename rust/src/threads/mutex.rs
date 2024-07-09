use std::rc::Rc;
use std::sync::{Arc, Mutex};
use std::thread::JoinHandle;

fn mutex() {
    let mutex = Arc::new(Mutex::new(0));
    let mut handles: Vec<JoinHandle<()>> = vec![];
    
    for _ in 0..10 {
        let counter = Arc::clone(&mutex);

        let handle = std::thread::spawn(move || {
            let mut num = counter.lock().unwrap();
            *num = *num + 1;
        });

        handles.push(handle);
    }

    for handle in handles {
        handle.join().unwrap();
    }

}