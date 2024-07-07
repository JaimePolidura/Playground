use std::collections::{HashMap, LinkedList};

pub fn vectors() {
    let mut vector: Vec<i32> = Vec::new();
    vector.push(2);
    vector.push(1);
    let numero: &i32 = &vector[1];
    println!("{}", numero);

    for i in vector {
        println!("{}", i);
    }

    let mut scores = HashMap::new();
    scores.insert(String::from("Blue"), 10);
    scores.insert(String::from("Yellow"), 50);
    let team_name = String::from("Blue");
    let score = scores.get(&team_name).copied().unwrap_or(0);

    for (key, value) in &scores {
        println!("{key}: {value}");
    }

    let ll: LinkedList<u32> = LinkedList::new();
    ll.back();

}