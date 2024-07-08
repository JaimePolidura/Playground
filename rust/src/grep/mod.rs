use std::{env, process};
use std::error::Error;
use std::io::Write;

mod params;
mod minigrep;

pub fn minigrep () {
    let params = params::read_from_args(&env::args())
        .unwrap_or_else(|err| {
            eprintln!("Problem parsing arguments: {}", err);
            process::exit(1);
        });

    minigrep::run(params).inspect_err(|err| {
        eprintln!("Application error: {err}");
        process::exit(1);
    });
}