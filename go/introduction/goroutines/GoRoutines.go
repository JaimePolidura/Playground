package main

import (
	"fmt"
	"sync"
	"time"
)

func main()  {
	var waitGroup sync.WaitGroup

	for i := 0; i < 5; i++ {
		waitGroup.Add(1)
		go func(i int) {
			defer waitGroup.Done()

			time.Sleep(time.Duration(i) * time.Second)

			fmt.Println("Finished", i)
		}(i)
	}

	fmt.Println("Waiting")
	waitGroup.Wait()
	fmt.Println("Finished")
}
