package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {
	var c = make(chan string)
	go pinger(c)
	go printer(c)
	var input string
	fmt.Scanln(&input)
}

func pinger(c chan<- string) {
	for i := 0; ; i++ {
		c <- "ping " + strconv.Itoa(i)
	}
}

func printer(c <-chan string) { //O en los argumentos podemos tener: func printer(c <-chan string) {
	for {
		msg := <- c
		fmt.Println(msg)
		time.Sleep(time.Second * 1)
	}
}