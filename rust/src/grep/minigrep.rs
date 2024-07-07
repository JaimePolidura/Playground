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
    let mut matches: Vec<&'a str> = Vec::new();

    for line in content.lines() {
        if line.contains(query) {
            matches.push(line);
        }
    }

    return matches;
}

fn search_ignore_case<'a>(
    content: &'a str,
    query: &'a str
) -> Vec<&'a str> {
    let mut matches: Vec<&'a str> = Vec::new();

    for line in content.lines() {
        if line.to_lowercase().contains(&query.to_lowercase()) {
            matches.push(line);
        }
    }

    return matches;
}