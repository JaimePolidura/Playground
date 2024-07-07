mod params;
mod minigrep;

use std::{env, fs, process};
use std::error::Error;
use std::io::{stderr, Write};

pub fn minigrep () {
    let args: Vec<String> = env::args().collect();

    let params = params::read_from_args(&args)
        .unwrap_or_else(|err| {
            eprintln!("Problem parsing arguments: {}", err);
            process::exit(1);
        });

    minigrep::run(params).inspect_err(|err| {
        eprintln!("Application error: {err}");
        process::exit(1);
    });
}