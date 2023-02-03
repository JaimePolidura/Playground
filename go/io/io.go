package main

import (
	"fmt"
	"os"
)

func main()  {
	file, err := os.Open("test.txt")

	if err != nil {
		fmt.Println("Fatal error")
		return
	}

	defer file.Close()

	stat, err := file.Stat()

	buffer := make([]byte, stat.Size())
	_, err = file.Read(buffer)

	if err != nil {
		fmt.Println("Fatal error")
		return
	}

	fileContentString := string(buffer)

	fmt.Println(fileContentString)
}


