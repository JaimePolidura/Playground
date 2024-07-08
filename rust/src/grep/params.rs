pub struct Params {
    pub ignore_case: bool,
    pub file_path: String,
    pub query: String,
}

pub fn read_from_args(
    mut args: impl Iterator<Item = String>
) -> Result<Params, &'static str> {
    args.next(); //Ignore program name

    let ignore_case = true;
    let query = match args.next() {
        Some(arg) => arg,
        None => return Err("Didn't get a query string"),
    };
    let file_path = match args.next() {
        Some(arg) => arg,
        None => return Err("Didn't get a file path"),
    };

    Ok(Params { query, file_path, ignore_case })
}
