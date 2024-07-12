#[cfg(test)]
mod test {
    use std::fmt::Pointer;
    use crate::list::stream::Stream;

    #[test]
    fn to_vec() {
        let values: Vec<i32> = Stream::of(vec![1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
            .filter(|item| item % 2 == 0)
            .map(|item| item * 2)
            .to_vec();

        assert_eq!(values, vec![4, 8, 12, 16, 20]);
    }

    #[test]
    fn count() {
        let n_values: usize = Stream::of(vec![1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
            .filter(|item| item % 2 == 0)
            .map(|item| item * 2)
            .count();

        assert_eq!(n_values, 5);
    }

    #[test]
    fn all_match() {
        let all_match: bool = Stream::of(vec![1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
            .filter(|item| item % 2 == 0)
            .map(|item| item * 2)
            .all_match(|item| *item > 0);
        assert!(all_match);

        let all_match: bool = Stream::of(vec![1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
            .filter(|item| item % 2 == 0)
            .map(|item| item * 2)
            .all_match(|item| *item <= 0);
        assert!(!all_match);
    }

    #[test]
    fn none_match() {
        let none_match: bool = Stream::of(vec![1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
            .filter(|item| item % 2 == 0)
            .map(|item| item * 2)
            .none_match(|item| *item == 0);
        assert!(none_match);
    }

    #[test]
    fn any_match() {
        let any_match: bool = Stream::of(vec![1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
            .filter(|item| item % 2 == 0)
            .map(|item| item * 2)
            .any_match(|item| *item == 20);
        assert!(any_match);

        let any_match: bool = Stream::of(vec![1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
            .filter(|item| item % 2 == 0)
            .map(|item| item * 2)
            .any_match(|item| *item == 3);
        assert!(!any_match);
    }

    #[test]
    fn first() {
        let first: Option<i32> = Stream::of(vec![1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
            .filter(|item| item % 2 == 0)
            .map(|item| item * 2)
            .first();
        assert!(first.is_some());
        assert_eq!(first.unwrap(), 4);
    }

    #[test]
    fn last() {
        let last: Option<i32> = Stream::of(vec![1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
            .filter(|item| item % 2 == 0)
            .map(|item| item * 2)
            .last();
        assert!(last.is_some());
        assert_eq!(last.unwrap(), 20);
    }

    #[test]
    fn limit() {
        let limitted: Vec<i32> = Stream::of(vec![1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
            .filter(|item| item % 2 == 0)
            .map(|item| item * 2)
            .limit(2)
            .to_vec();

        assert_eq!(limitted.len(), 2);
        assert_eq!(limitted, vec![4, 8]);
    }

    #[test]
    fn add() {
        let add: i32 = Stream::of(vec![1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
            .filter(|item| item % 2 == 0)
            .map(|item| item * 2)
            .limit(2)
            .add();

        assert_eq!(add, 12);
    }

    #[test]
    fn reduced() {
        let reduced: i32 = Stream::of(vec![1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
            .filter(|item| item % 2 == 0)
            .map(|item| item * 2)
            .limit(2)
            .reduce(0, |current, reduced| current + reduced);

        assert_eq!(reduced, 12);
    }

    #[test]
    fn for_each() {
        let mut items = Vec::new();

        Stream::of(vec![1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
            .filter(|item| item % 2 == 0)
            .map(|item| item * 2)
            .limit(2)
            .for_each(|current| items.push(*current));

        assert_eq!(items, vec![4, 8]);
    }

    #[test]
    fn max() {
        let max = Stream::of(vec![1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
            .max(|a, b| *a > *b);

        assert!(max.is_some());
        assert_eq!(max.unwrap(), 10);
    }

    #[test]
    fn min() {
        let min = Stream::of(vec![1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
            .min(|a, b| *a < *b);

        assert!(min.is_some());
        assert_eq!(min.unwrap(), 1);
    }
}