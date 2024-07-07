#[cfg(test)]
mod tests {
    use crate::list;
    use crate::list::LinkedList;

    #[test]
    fn initialization() {
        let ll: LinkedList<u32> =  list::LinkedList::new();
        assert!(ll.is_empty());
        assert_eq!(ll.len(), 0);
    }

    #[test]
    fn add() {
        let mut ll: LinkedList<u32> =  list::LinkedList::new();
        ll.add_first(1);
        ll.add_first(2);

        assert!(!ll.is_empty());
        assert_eq!(ll.len(), 2);
    }
    
    #[test]
    // #[should_panic]
    fn get() {
        let mut ll: LinkedList<u32> =  list::LinkedList::new();
        ll.add_first(1);
        ll.add_first(2);
        ll.add_first(3);

        assert_eq!(*(ll.get_first().unwrap()), 3);
        assert_eq!(*(ll.get_last().unwrap()), 1);
        assert_eq!(*(ll.get(1).unwrap()), 2);
    }
}