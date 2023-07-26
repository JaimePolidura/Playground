package nodes

import (
	"distributed-systems/src/utils"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"sync"
)

type MessageListener struct {
	selfNodeId        uint32
	selfPort          uint16
	bufferMessageSize sync.Pool
}

func CreateMessageListener(selfNodeId uint32, selfPort uint16) *MessageListener {
	return &MessageListener{
		selfNodeId: selfNodeId,
		selfPort:   selfPort,
		bufferMessageSize: sync.Pool{New: func() interface{} {
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
		bufferSize := listener.bufferMessageSize.Get().([]byte)

		messages, err := listener.deserializeMessages(conn, bufferSize)

		if err != nil {
			fmt.Printf("[%d] ERROR %s %s\n", listener.selfNodeId, "message_listener.go:handleNewConnection", err.Error())
			continue
		}

		onReadCallback(messages)
		utils.ZeroArray(&bufferSize)

		listener.bufferMessageSize.Put(bufferSize)
	}
}

func (listener *MessageListener) deserializeMessages(conn net.Conn, bufferSize []byte) ([]*Message, error) {
	messages := make([]*Message, 0)

	conn.Read(bufferSize)
	messageSize := binary.BigEndian.Uint32(bufferSize)
	messageBuffer := make([]byte, messageSize)
	start := uint32(0)

	for start+1 < messageSize {
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
