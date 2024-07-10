mod post;

pub trait Draw {
    fn draw(&self);
}

pub struct Screen {
    pub components: Vec<Box<dyn Draw>>,
}

impl Screen {
    pub fn run(&self) {
        for component in self.components {
            component.draw();
        }
    }
}

pub struct Screen2<T: Draw> {
    pub components: Vec<T>
}

impl<T> Screen2<T>
where
    T: Draw
{
    pub fn draw(&self) {
        for component in self.components {
            component.draw();
        }
    }
}

fn oop() {

}