use crate::enums::{messages, store};
use crate::enums::messages::Request;

pub struct Server;
pub struct Connection;

pub struct ServerRequestContext {
    pub request: Request<String>,
    pub connection: Connection
}

impl Connection {
    pub fn respond(&self, result: store::Result<String>) {

    }
}

impl Server {
    pub fn read_request(&self) -> ServerRequestContext {
        ServerRequestContext {
         request: Request {
             timestamp: messages::LamportTimestamp{
                 counter: 100,
                 node_id: 1
             },
             auth_key: 9178629871261,
             opcode: messages::Opcode::CAS(String::from("Ejemplo"), String::from("A"), String::from("B"))
         },
            connection: Connection
        }
    }
}