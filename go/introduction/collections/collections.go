package main

import "container/list"
import "container/ring"

func main()  {
	var list list.List
	list.PushBack(1)

	var ring ring.Ring
	ring.Move(1)
}
