mod tests;

struct Node<T> {
    data: T,
    next: Option<Box<Node<T>>>,
}

pub struct LinkedList<T> {
    head: Option<Box<Node<T>>>,
    count: u32
}

impl<T> LinkedList<T> {
    pub fn new() -> LinkedList<T> {
        LinkedList {
            head: None,
            count: 0
        }
    }

    pub fn remove_first(&mut self) {
        self.remove(0);
    }

    pub fn remove_last(&mut self) {
        self.remove(self.count - 1);
    }

    //TODO
    pub fn remove(&mut self, index_to_remove: u32) -> Option<T> {
        None
    }

    pub fn len(&self) -> u32 {
        return self.count;
    }

    pub fn is_empty(&self) -> bool {
        return self.count == 0;
    }

    pub fn get_last(&self) -> Option<&T> {
        return self.get(self.count - 1);
    }

    pub fn get_first(&self) -> Option<&T> {
        return self.get(0);
    }

    pub fn get(&self, index_lookup: u32) -> Option<&T> {
        if self.is_out_of_bounds(index_lookup) {
            return Option::None;
        }

        let mut last: Option<&Box<Node<T>>> = None;

        for current_index in 0..(index_lookup + 1) {
            last = match last {
                Some(prev) => prev.next.as_ref(),
                None => self.head.as_ref(),
            };

            if current_index == index_lookup {
                return Some(&last.unwrap().data);
            }
        }

        None
    }

    pub fn add_first(&mut self, data: T) {
        let new_node = Box::new(Node {
            data: data,
            next: self.head.take()
        });

        self.head = Some(new_node);
        self.count = self.count + 1;
    }

    pub fn add_last(&mut self, data: T) {
        match self.head {
            None => self.add_first(data),
            Some(_) => {
                let mut last_node: &mut Box<Node<T>> = self.get_last_node();

                last_node.next = Some(Box::new(Node {
                    data: data,
                    next: None
                }));
                self.count = self.count + 1;
            }
        }
    }

    fn get_last_node(&mut self) -> &mut Box<Node<T>> {
        let mut current = self.head.as_mut().expect("");

        while current.next.is_some() {
            current = current.next.as_mut().unwrap();
        }

        current
    }

    fn is_out_of_bounds(&self, index: u32) -> bool {
        return self.count == 0 || index < 0 || index + 1 > self.count;
    }
}
