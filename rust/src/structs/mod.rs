mod lifetimes;

use std::ptr::null;

#[derive(Debug)]
struct User {
    name: String,
    email: String,
    sign_count: u64,
    state: UserState,
    position: Position
}

#[derive(Debug)]
enum UserState {
    DELETED,
    BANNED,
    CREATED
}

impl User {
    fn increase_sign_count(&mut self) {
        self.sign_count = self.sign_count + 1;
    }

    fn create(name: &str, email: &str) -> Self {
        User {
            name: name.to_string(),
            email: email.to_string(),
            sign_count: 1,
            state: UserState::CREATED,
            position: Position(0, 0, 0),
        }
    }
}

#[derive(Debug)]
struct Position(i32, i32, i32);

struct Empty;

pub fn structs() {
    let mut user = User{
        name: String::from("Jaime"),
        email: String::from("prueba@gmail.com"),
        sign_count: 10,
        state: UserState::CREATED,
        position: Position(1, 1, 1)
    };

    User::create("jaime", "molon@gmail.com");

    user.increase_sign_count();

    println!("{:?}", user);
}

fn get_name(user: &User) -> &str {
    &user.name[..]
}