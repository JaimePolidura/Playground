package main

import (
	"fmt"
	"strconv"
)

func main() {
	for i := 1; i < 6; i++ {
		if i & 1 == 0 {
			fmt.Println(strconv.Itoa(i) + " is even")
		} else {
			fmt.Println(strconv.Itoa(i) + " is odd")
		}
	}
}

func printRange() {
	i := 1

	for i <= 3 {
		fmt.Println(i)
		i++
	}
}
