fn pattern_matching() {
    let mut stack = Vec::new();
    stack.push(1);
    stack.push(2);
    stack.push(3);

    while let Some(top) = stack.pop() {
        println!("{top}");
    }

    let (x, y, z) = (1, 2, 3);

    let x = 1;
    match x {
        1 => println!("one"),
        2..=3 => println!("two"),
        3..=10 => println!("three"),
        11 | 12 => println!("mierdon"),
        _ => println!("anything"),
    }

    let num = Some(4);
    match num {
        Some(x) if x % 2 == 0 => println!("The number {x} is even"),
        Some(x) => println!("The number {x} is odd"),
        None => (),
    }
}

fn print_coordinates(&(x, y): &(i32, i32)) {
    println!("Current location: ({x}, {y})");
}

struct Point {
    x: i32,
    y: i32,
}

fn main() {
    let point = (3, 5);
    print_coordinates(&point);

    let p = Point { x: 0, y: 7 };
    let Point { x, y } = p;

    match p {
        Point { x, y: 0 } => println!("On the x axis at {x}"),
        Point { x: 0, y } => println!("On the y axis at {y}"),
        Point { x, y } => {
            println!("On neither axis: ({x}, {y})");
        }
    }
}