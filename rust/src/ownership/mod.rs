pub fn ownership() {
    let hola = String::from("Hola");
    let adios = hola;
    // println!("Hola {adios}"); Invalid
    let adios = function(adios);
    println!("{}", adios);
}

fn function(string: String) -> String {
    string
}