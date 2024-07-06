mod messages;
mod store;
mod server;

struct DB {
    server: server::Server,
    store: store::Store<String>
}

impl DB {
    fn start(&mut self) -> ! {
        loop {
            let ctx: server::ServerRequestContext = self.server.read_request();

            let result: store::Result<String> = match ctx.request.opcode {
                messages::Opcode::Set(key, value) => self.store.set(&key, value),
                messages::Opcode::Get(key) => self.store.get(&key),
                messages::Opcode::Contains(key) => self.store.contains(&key),
                messages::Opcode::CAS(key, expected, new) => self.store.cas(&key, expected, new),
                messages::Opcode::Delete(key) => self.store.delete(&key),
            };

            ctx.connection.respond(result);
        }
    }
}