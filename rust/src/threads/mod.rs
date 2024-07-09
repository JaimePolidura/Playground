mod mutex;

use std::sync::mpsc;

fn threads() {
    let (tx, rx) = mpsc::channel();

    let handle = std::thread::spawn(move || {
        let vals = vec![
            String::from("hi"),
            String::from("from"),
            String::from("the"),
            String::from("thread"),
        ];
        for val in vals {
            tx.send(val).unwrap();
            std::thread::sleep(std::time::Duration::from_secs(1));
        }
    });

    for received in rx {
        println!("Got: {received}");
    }

    handle.join();
}