use std::error::Error;
use std::fs;
use crate::grep::params::Params;

pub fn run(params: Params) -> Result<(), Box<dyn Error>>{
    let contents: String = fs::read_to_string(&params.file_path)?;

    let results = if params.ignore_case {
        search_case_sensitive(&contents, &params.query)
    } else {
        search_ignore_case(&contents, &params.query)
    };

    for line in results {
        println!("{}", line);
    }

    Ok(())
}

fn search_case_sensitive<'a>(
    content: &'a str,
    query: &'a str
) -> Vec<&'a str> {
    content.lines()
        .filter(|line| line.contains(query))
        .collect()
}

fn search_ignore_case<'a>(
    content: &'a str,
    query: &'a str
) -> Vec<&'a str> {
    content.lines()
        .filter(|line| line.to_lowercase().contains(&query.to_lowercase()))
        .collect()
}