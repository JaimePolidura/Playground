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

	iterator := list.Iterate()

	for iterator.HastNext() {
		fmt.Println("Value ", iterator.Next())
	}
}
