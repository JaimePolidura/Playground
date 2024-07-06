struct Point<T> {
    x: T,
    y: T,
}

//Solo para los Point con tipo i32
impl Point<i32>  {
    fn get(&self) -> i32 {
        1
    }
}

fn genercis() {
}

fn get_largest<T: PartialOrd>(arr: &[T]) -> Option<&T> {
    if arr.len() == 0 {
        return None;
    }

    let mut largest: &T = &arr[0];

    for item in arr {
        if item > largest {
            largest = item;
        }
    }

    return Some(largest);
}