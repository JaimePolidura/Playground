package main

import (
	"fmt"
	"go/linkedlist/code"
)

func main()  {
	list := &code.LinkedList[int]{}

	list.Add(1)
	list.Add(2)
	list.Add(3)

	list.Stream().Filter(func(it int) bool {
		return it % 2 == 0
	}).Foreach(func(it int) {
		fmt.Println(it)
	})
}
