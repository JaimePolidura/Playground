pub struct Params {
    pub ignore_case: bool,
    pub file_path: String,
    pub query: String,
}

pub fn read_from_args(args: &Vec<String>) -> Result<Params, &str> {
    if args.len() != 3 {
        return Err("Invalid arguments");
    }

    let query = args[1].clone();
    let file_path = args[2].clone();
    let ignore_case = true;

    Ok(Params { query, file_path, ignore_case })
}
