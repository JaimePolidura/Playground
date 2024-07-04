use std::cmp::Ordering;
use std::io;
use rand::Rng;

const FROM: i32 = 1;
const TO: i32 = 100;

pub fn game() {
    let mut guess: String = String::new();
    let random: String = rand::thread_rng().gen_range(FROM..=TO).to_string();

    let a: (i8, i16, i32, i64, i128) = (TO as i8, 2, 3, 4, 5);

    loop {
        println!("Guess the number!");

        io::stdin()
            .read_line(&mut guess)
            .expect("Failed to read the line");
        let mut guess = guess.trim();

        match guess.cmp(&random) {
            Ordering::Less => println!("Greater!"),
            Ordering::Greater => println!("Less!"),
            Ordering::Equal => break
        }
    }

    println!("You won!");
}
