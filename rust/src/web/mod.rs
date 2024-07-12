use std::io::{BufRead, Write};
use std::net::{TcpListener, TcpStream};

use crate::web::pool::ThreadPool;

mod pool;

const HTML_TO_RETURN: &str =
    r#"<!DOCTYPE html>
    <html lang="en">
        <head>
            <meta charset="utf-8">
            <title>Hello!</title>
        </head>
        <body>
            <h1>Hello!</h1>
            <p>Hi from Rust</p>
        </body>
    </html>"#;

pub fn main() {
    let listener = TcpListener::bind("127.0.0.1:7878").unwrap();
    let mut thread_pool = ThreadPool::new(8);

    for stream in listener.incoming() {
        let stream = stream.unwrap();
        thread_pool.execute(|| handle_connection(stream));
    }
}

fn handle_connection(mut stream: TcpStream) {
    let status_line = "HTTP/1.1 200 OK";
    let length = HTML_TO_RETURN.len();

    let response = format!(
        "{status_line}\r\n\
        Content-Length: {length}\r\n\r\n\
        {HTML_TO_RETURN}"
    );

    stream.write_all(response.as_bytes()).unwrap();
}