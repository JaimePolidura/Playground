mod messenger;

use std::rc::Rc;

struct Person {
    name: String,
    age: i8,
}

enum List {
    Data(i32, Rc<List>),
    Null
}

pub fn pointers() {
    let a = Rc::new(List::Data(5, Rc::new(List::Data(10, Rc::new(List::Null)))));
    let b = List::Data(3, Rc::clone(&a));
    let c = List::Data(4, Rc::clone(&a));

    let mut person: Rc<Person> = Rc::new(Person{name: String::from("Jaime"), age: 21});
    let mut copy = Rc::clone(&person);

    // copy.age = copy.age + 1;
    // person.age = person.age + 1;
}
