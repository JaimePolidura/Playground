package nodes

import (
	"distributed-systems/src/utils"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"sync"
)

type NewMessage struct {
	from    net.Conn
	message *Message
}

type MessageListener struct {
	selfNodeId        uint32
	selfPort          uint16
	bufferMessageSize sync.Pool

	listener            net.Listener
	nodeConnectionStore *NodeConnectionsStore

	newMessage chan NewMessage
}

func CreateMessageListener(selfNodeId uint32, selfPort uint16) *MessageListener {
	return &MessageListener{
		selfNodeId:          selfNodeId,
		selfPort:            selfPort,
		newMessage:          make(chan NewMessage, 100),
		nodeConnectionStore: CreateNodeConnectionStore(),
		bufferMessageSize: sync.Pool{New: func() interface{} {
			return make([]byte, 4)
		}},
	}
}

func (this *MessageListener) Stop() {
	this.listener.Close()
}

func (this *MessageListener) ListenAsync() {
	lis, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(int(this.selfPort)))

	if err != nil {
		fmt.Printf("[%d] ERROR %s %s\n", this.selfNodeId, "message_listener.go:ListenAsync", err.Error())
		return
	}

	this.listener = lis

	fmt.Printf("[%d] Started listening on port %d\n", this.selfNodeId, this.selfPort)

	go func() {
		for {
			conn, _ := lis.Accept()

			go this.handleNewConnection(conn)
		}
	}()
}

func (this *MessageListener) handleNewConnection(conn net.Conn) {
	defer func() {
		recover() //Interrupted by close
	}()

	for {
		bufferSize := this.bufferMessageSize.Get().([]byte)

		messages, err := this.deserializeMessages(conn, bufferSize)

		if err != nil {
			fmt.Printf("[%d] ERROR %s %s\n", this.selfNodeId, "message_listener.go:handleNewConnection", err.Error())
			continue
		}

		for _, message := range messages {
			this.newMessage <- NewMessage{
				from:    conn,
				message: message,
			}
		}

		utils.ZeroArray(&bufferSize)

		this.bufferMessageSize.Put(bufferSize)
	}
}

func (this *MessageListener) deserializeMessages(conn net.Conn, bufferSize []byte) ([]*Message, error) {
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
