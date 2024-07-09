// pub struct Stream<T, I: Iterator<Item = T>> {
//     iterator: I
// }
//
// impl<T, I> Stream<T, I> {
//     pub fn of<Iterable>(iterable: Iterable) -> Stream<T, Iterable::IntoIter>
//     where
//         Iterable: IntoIterator<Item = T>
//     {
//         return Stream { iterator: iterable.into_iter() };
//     }
//
//     pub fn filter<F>(&mut self, predicate: F) -> Stream<T, I>
//     where
//         F: Fn(T) -> bool
//     {
//         let mut mierdon = &self.iterator;
//
//         let mut filtered: Vec<T> = Vec::new();
//
//         return Stream { iterator: filtered.into_iter() };
//     }
// }