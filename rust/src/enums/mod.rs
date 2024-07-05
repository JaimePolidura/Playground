mod messages;
mod store;
mod server;

use std::io;

struct DB {
    server: server::Server,
    store: store::Store<String>
}

impl DB {
    fn start(&mut self) -> ! {
        loop {
            let ctx: server::ServerRequestContext = self.server.read_request();

            let result: store::Result<String> = match self.read_request(&ctx.request) {
                messages::Opcode::Set(key, value) => self.store.set(&key, value),
                messages::Opcode::Get(key) => self.store.get(&key).ok_or(1),
                messages::Opcode::Contains(key) => self.store.contains(&key),
                messages::Opcode::CAS(key, expected, new) => self.store.cas(&key, expected, new),
                messages::Opcode::Delete(key) => Ok(self.store.delete(&key))
            };

            ctx.connection.respond(result);
        }
    }

    fn read_request(&self, request: &messages::Request<String>) -> io::Result<messages::Response> {
        return Ok(messages::Response{status_code: 200});
    }
}