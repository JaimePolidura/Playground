pub enum Opcode<T> {
    Set(String, T),
    Get(String),
    Contains(String),
    CAS(String, T, T),
    Delete(String)
}

pub struct LamportTimestamp {
    pub counter: u64,
    pub node_id: u16
}

pub struct Request<T> {
    pub timestamp: LamportTimestamp,
    pub auth_key: u64,
    pub opcode: Opcode<T>,
}

pub struct Response {
    pub status_code: u8
}