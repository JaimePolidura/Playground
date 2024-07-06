use std::io::Write;
use std::fmt::Display;

trait Serializable {
    fn serialize(&self) -> Vec<u8>;
}

enum Options {
    Compress,
    NotCompress
}

trait HasSerializationOptions {
    fn get_options(&self) -> Options;
}

struct Message<T: Serializable> {
    opcode: u8,
    body: T,
}

impl<T: Serializable> Serializable for Message<T> {
    fn serialize(&self) -> Vec<u8> {
        let mut result: Vec<u8> =  Vec::new();
        result.push(self.opcode);
        let mut body_serialized = self.body.serialize();
        result.write_all(&body_serialized);

        return result;
    }
}

fn serialize_and_send<T>(to_serialize: T)
where
    T: Serializable + HasSerializationOptions
{
    let options = to_serialize.get_options();
    let serialized = to_serialize.serialize();

    //TODO
}

struct Pair<T> {
    x: T,
    y: T,
}

impl<T> Pair<T> {
    fn new(x: T, y: T) -> Self {
        Self { x, y }
    }
}

impl<T: Display + PartialOrd> Pair<T> {
    fn cmp_display(&self) {
        if self.x >= self.y {
            println!("The largest member is x = {}", self.x);
        } else {
            println!("The largest member is y = {}", self.y);
        }
    }
}