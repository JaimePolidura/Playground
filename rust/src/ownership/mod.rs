pub fn ownership() {
    let hola = String::from("Hola");
    let adios = hola;
    // println!("Hola {adios}"); Invalid
    let adios = function(adios);
    println!("{}", adios);

    stringLength(&adios);
    println!("{}", adios);
}

fn first_word(s: &String) -> &str {
    let bytes = s.as_bytes();
    
    for (i, &item) in bytes.iter().enumerate() {
        if item == b' ' {
            &s[0..i];
        }
    }

    &s[..]
}

fn stringLength(mut string: &String) -> usize {
    return string.len();
}

fn function(string: String) -> String {
    string
}