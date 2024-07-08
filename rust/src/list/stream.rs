struct Stream<T> {
}

impl<T> Stream<T> {
    pub fn of(mut iterator: impl Iterator<Item = T>) -> Stream<T> {
        Stream {}
    }
}