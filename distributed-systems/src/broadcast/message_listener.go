package broadcast

import (
	"fmt"
	"net"
	"strconv"
)

type MessageListener struct {
	selfNodeId uint32
	selfPort   uint16
	buffer     []byte

	connectionsByNodeId map[uint32]*NodeConnection
}

func CreateMessageListener(selfNodeId uint32, selfPort uint16) *MessageListener {
	return &MessageListener{
		selfNodeId:          selfNodeId,
		selfPort:            selfPort,
		buffer:              make([]byte, 1024),
		connectionsByNodeId: make(map[uint32]*NodeConnection),
	}
}

func (listener *MessageListener) ListenAsync(onReadCallback func(message *Message)) {
	conn, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(int(listener.selfPort)))

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	go func() {
		for {
			conn, _ := conn.Accept()
			go listener.handleConnection(conn, onReadCallback)
		}
	}()
}

func (listener *MessageListener) handleConnection(conn net.Conn, onReadCallback func(message *Message)) {
	for {
		conn.Read(listener.buffer)

		message, err := Deserialize(listener.buffer)

		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		listener.connectionsByNodeId[message.NodeId] = CreateNodeConnection(message.NodeId, conn)
		onReadCallback(message)
		ZeroArray(&listener.buffer)
	}
}

func (listener *MessageListener) Write(nodeId uint32, message *Message) {
	listener.connectionsByNodeId[nodeId].Write(message)
}
