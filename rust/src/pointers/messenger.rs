use std::cell::RefCell;

pub trait Messenger {
    fn send(&self, msg: &str);
}

struct MockMessenger {
    sent_messages: RefCell<Vec<String>>,
}

impl MockMessenger {
    fn new() -> MockMessenger {
        MockMessenger {
            sent_messages: RefCell::new(vec![]),
        }
    }
}

impl Messenger for MockMessenger {
    fn send(&self, message: &str) {
        self.sent_messages
            .borrow_mut()
            .push(String::from(message));
    }
}

pub struct LimitTracker<T: Messenger> {
    messenger: T,
    value: usize,
    max: usize,
}

impl<T> LimitTracker<T>
where
    T: Messenger,
{
    pub fn new(
        messenger: T,
        max: usize
    ) -> LimitTracker<T> {
        LimitTracker {
            messenger,
            value: 0,
            max,
        }
    }
    
    pub fn set_value(&mut self, value: usize) {
        self.value = value;
        let percentage_of_max =
            self.value as f64 / self.max as f64;
        if percentage_of_max >= 1.0 {
            self.messenger
                .send("Error: You are over your quota!");
        } else if percentage_of_max >= 0.9 {
            self.messenger
                .send("Urgent: You're at 90% of your quota!");
        } else if percentage_of_max >= 0.75 {
            self.messenger
                .send("Warning: You're at 75% of your quota!");
        }
    }
}