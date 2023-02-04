package main

import (
	"io/ioutil"
	"net/http"
)

type HomePageSize struct {
	url string
	size int
}

func main()  {
	results := make(chan HomePageSize)

	urls := []string{
		"http://www.apple.com",
		"http://www.amazon.com",
		"http://www.google.com",
		"http://www.microsoft.com",
	}

	for _, url := range urls {
		go func(url string) {
			res, err := http.Get(url)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()

			bs, _ := ioutil.ReadAll(res.Body)

			results <- HomePageSize{
				url: url,
				size: len(bs),
			}
		}(url)
	}
}


