package broadcast

import (
	"fmt"
	"net"
	"strconv"
	"sync"
)

type MessageListener struct {
	selfNodeId uint32
	selfPort   uint16
	buffer     sync.Pool
}

func CreateMessageListener(selfNodeId uint32, selfPort uint16) *MessageListener {
	return &MessageListener{
		selfNodeId: selfNodeId,
		selfPort:   selfPort,
		buffer: sync.Pool{New: func() interface{} {
			return make([]byte, 1024)
		}},
	}
}

func (listener *MessageListener) ListenAsync(onReadCallback func(message *Message)) {
	conn, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(int(listener.selfPort)))

	if err != nil {
		fmt.Printf("[%d] ERROR %s %s\n", listener.selfNodeId, "message_listener.go:ListenAsync", err.Error())
		return
	}

	fmt.Printf("[%d] Started listening on port %d\n", listener.selfNodeId, listener.selfPort)

	go func() {
		for {
			conn, _ := conn.Accept()

			go listener.handleNewConnection(conn, onReadCallback)
		}
	}()
}

func (listener *MessageListener) handleNewConnection(conn net.Conn, onReadCallback func(message *Message)) {
	for {
		buffer := listener.buffer.Get().([]byte)

		conn.Read(buffer)

		message, err := Deserialize(buffer)

		if err != nil {
			fmt.Printf("[%d] ERROR %s %s\n", listener.selfNodeId, "message_listener.go:handleNewConnection", err.Error())
			continue
		}

		onReadCallback(message)
		ZeroArray(&buffer)

		listener.buffer.Put(buffer)
	}
}
