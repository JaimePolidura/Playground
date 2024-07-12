use std::ops::Add;

pub struct Stream<Input> {
    values: Option<Vec<Input>>
}

impl<Input: 'static> Stream<Input>
where
    Input: Add<Output = Input> + Default + Copy
{
    pub fn add(&mut self) -> Input {
        let mut add_result = Input::default();

        for &value in self.values.as_ref().unwrap() {
            add_result = add_result + value;
        }

        add_result
    }
}

impl<Input: 'static> Stream<Input> {
    pub fn of<I>(iterator: I) -> Stream<Input>
    where
        I: IntoIterator<Item =Input>, <I as IntoIterator>::IntoIter: 'static
    {
        Stream { values: Some(iterator.into_iter().collect()) }
    }

    pub fn filter<Predicate>(&mut self, predicate: Predicate) -> Stream<Input>
    where
        Predicate: Fn(&Input) -> bool
    {
        let mut filtered: Vec<Input> = Vec::new();

        while self.values.as_ref().unwrap().len() > 0 {
            let item: Input = self.values.as_mut().unwrap().remove(0);
            if predicate(&item) {
                filtered.push(item);
            }
        }

        Stream { values: Some(filtered) }
    }

    pub fn map<Mapper, Output: 'static>(&mut self, mapper: Mapper) -> Stream<Output>
    where
        Mapper: Fn(&Input) -> Output
    {
        let mut mapped_values: Vec<Output> = Vec::new();

        for item in self.values.as_ref().unwrap() {
            let mapped: Output = mapper(&item);
            mapped_values.push(mapped);
        }

        Stream {values: Some(mapped_values)}
    }

    pub fn limit(&mut self, limit: usize) -> Stream<Input> {
        let mut values: Vec<Input> = Vec::with_capacity(limit);

        for i in 0..limit {
            if self.values.as_ref().unwrap().len() > 0 {
                let value = self.values.as_mut().unwrap().remove(0);
                values.insert(i, value);
            } else {
                break
            }
        }

        Stream { values: Some(values) }
    }

    pub fn for_each<Consumer>(&mut self, mut consumer: Consumer) -> Stream<Input>
    where
        Consumer: FnMut(&Input)
    {
        for item in self.values.as_ref().unwrap() {
            consumer(item);
        }

        Stream {values: Some(self.values.take().unwrap())}
    }

    //Comparator(a, b) returns true if a > b
    pub fn max<Comparator>(&mut self, comparator: Comparator) -> Option<Input>
    where
        Comparator: Fn(&Input, &Input) -> bool
    {
        self.min_max(comparator)
    }

    //Comparator(a, b) returns true if a < b
    pub fn min<Comparator>(&mut self, comparator: Comparator) -> Option<Input>
    where
        Comparator: Fn(&Input, &Input) -> bool
    {
        self.min_max(comparator)
    }

    fn min_max<Comparator>(&mut self, comparator: Comparator) -> Option<Input>
    where
        Comparator: Fn(&Input, &Input) -> bool
    {
        let mut min_max_value_seen_option: Option<Input> = None;

        while let Some(current_value) = self.values.as_mut().unwrap().pop() {
            match min_max_value_seen_option.as_ref() {
                Some(max_value_seen) => {
                    if comparator(&current_value, max_value_seen) {
                        min_max_value_seen_option = Some(current_value);
                    }
                },
                None => min_max_value_seen_option = Some(current_value),
            }
        }

        min_max_value_seen_option
    }

    pub fn reduce<Output, Reducer>(&mut self, default: Output, reducer: Reducer) -> Output
    where
        Reducer: Fn(&Input, &Output) -> Output
    {
        let mut reduced = default;

        for item in self.values.as_ref().unwrap() {
            reduced = reducer(&item, &reduced);
        }

        reduced
    }

    pub fn all_match<Predicate>(&mut self, predicate: Predicate) -> bool
    where
        Predicate: Fn(&Input) -> bool
    {
        let prev_size = self.values.as_ref().unwrap().len();
        let new_size = self.filter(predicate).values.as_ref().unwrap().len();

        prev_size == new_size
    }

    pub fn none_match<Predicate>(&mut self, predicate: Predicate) -> bool
    where
        Predicate: Fn(&Input) -> bool
    {
        self.filter(predicate).values.as_ref().unwrap().is_empty()
    }

    pub fn any_match<Predicate>(&mut self, predicate: Predicate) -> bool
    where
        Predicate: Fn(&Input) -> bool
    {
        !self.filter(predicate).values.as_ref().unwrap().is_empty()
    }

    pub fn last(&mut self) -> Option<Input> {
        self.values.as_mut().unwrap().pop()
    }

    pub fn first(&mut self) -> Option<Input> {
        if !self.values.as_ref().unwrap().is_empty() {
            Some(self.values.as_mut().unwrap().remove(0))
        } else {
            None
        }
    }

    pub fn to_vec(self) -> Vec<Input> {
        return self.values.unwrap();
    }

    pub fn count(&self) -> usize {
        self.values.as_ref().unwrap().len()
    }
}