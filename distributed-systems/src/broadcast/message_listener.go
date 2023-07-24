package broadcast

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"sync"
)

type MessageListener struct {
	selfNodeId          uint32
	selfPort            uint16
	bufferMessageLength sync.Pool
}

func CreateMessageListener(selfNodeId uint32, selfPort uint16) *MessageListener {
	return &MessageListener{
		selfNodeId: selfNodeId,
		selfPort:   selfPort,
		bufferMessageLength: sync.Pool{New: func() interface{} {
			return make([]byte, 4)
		}},
	}
}

func (listener *MessageListener) ListenAsync(onReadCallback func(message []*Message)) {
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

func (listener *MessageListener) handleNewConnection(conn net.Conn, onReadCallback func(message []*Message)) {
	for {
		bufferLength := listener.bufferMessageLength.Get().([]byte)

		messages, err := listener.deserializeMessages(conn, bufferLength)

		if err != nil {
			fmt.Printf("[%d] ERROR %s %s\n", listener.selfNodeId, "message_listener.go:handleNewConnection", err.Error())
			continue
		}

		onReadCallback(messages)
		ZeroArray(&bufferLength)

		listener.bufferMessageLength.Put(bufferLength)
	}
}

func (listener *MessageListener) deserializeMessages(conn net.Conn, bufferLength []byte) ([]*Message, error) {
	messages := make([]*Message, 0)

	conn.Read(bufferLength)
	messageLength := binary.BigEndian.Uint32(bufferLength)
	messageBuffer := make([]byte, messageLength)
	start := uint32(0)
	
	for start < messageLength {
		conn.Read(messageBuffer)
		message, nextStart, err := Deserialize(messageBuffer, start)
		start += nextStart

		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}
