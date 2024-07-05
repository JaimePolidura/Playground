use std::collections::HashMap;

pub struct Store<T: PartialEq> {
    data: HashMap<String, T>
}

pub enum Result<T> {
    Failed,
    Success,
    SuccessWithValue(T)
}

impl<T: PartialEq> Store<T> {
    pub fn get(&self, key: &str) -> Result<T> {
        return match self.data.get(key) {
            Some(value) => Result::SuccessWithValue(value),
            None => Result::Failed
        }
    }

    pub fn set(&mut self, key: &str, value: T) -> Result<T> {
        self.data.insert(key.to_string(), value);
        return Result::Success;
    }

    pub fn delete(&mut self, key: &str) -> Result<T> {
        self.data.remove(key);
        return Result::Success;
    }

    pub fn cas(&mut self, key: &str, expected: T, new: T) -> Result<T> {
        match self.data.get(key) {
            Some(value) if *value == expected => {
                self.data.insert(key.to_string(), new);
                return Result::Success
            },
            _ => Result::Failed
        }
    }
}