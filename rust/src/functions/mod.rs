use std::io;
use rand::Error;

pub fn functions() {
    let result = {
        match advance((1, 2), (-1, 0)) {
            Ok(vec) => println!("The new position ({}, {})", vec.0, vec.1),
            Error=> println!("Error")
        }
    };

    let sum = {
        let mut last = 0;

        loop {
            last = last + 1;

            if last > 10 {
                break last
            }
        }
    };
}

fn advance(position: (i32, i32), vec: (i32, i32)) -> io::Result<(i32, i32)> {
    let new_position = (position.0 + vec.0, position.1 + vec.1);
    Ok(new_position)
}