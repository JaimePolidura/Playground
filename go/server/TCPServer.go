package main

import (
	"encoding/gob"
	"fmt"
	"go/server/useless"
	"net"
)

func main()  {
	useless.ImUseless()

	go server()
	go client()

	var input string
	fmt.Scanln(&input)
}

func server()  {
	listener, err := net.Listen("tcp", ":9999")

	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		connection, err := listener.Accept()

		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleConnection(connection)
	}
}

func handleConnection(connection net.Conn)  {
	defer connection.Close()

	var msg string
	err := gob.NewDecoder(connection).Decode(&msg)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("[SERVER] Recieved", msg)
}

func client()  {
	connection, err := net.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer connection.Close()

	messageToSend := "Hello from client"
	fmt.Println("[CLIENT] Sending ", messageToSend)

	err = gob.NewEncoder(connection).Encode(messageToSend)
	if err != nil {
		fmt.Println(err)
	}

}
