package broadcast

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"sync"
)

type NodeConnection struct {
	nativeConnection net.Conn
	nativeWriter     *bufio.Writer

	selfNodeId uint32
	port       uint32
}

func CreateNodeConnection(nodeId uint32, port uint32) *NodeConnection {
	return &NodeConnection{
		selfNodeId: nodeId,
		port:       port,
	}
}

func (this *NodeConnection) Open() {
	conn, err := net.Dial("tcp", "127.0.0.1:"+strconv.Itoa(int(this.port)))

	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	this.nativeWriter = bufio.NewWriter(conn)
	this.nativeConnection = conn
}

func (this *NodeConnection) WriteAll(messages []*BroadcastMessage) {
	for _, message := range messages {
		message.NodeIdSender = this.selfNodeId
	}

	serialized := SerializeAll(messages)
	this.nativeConnection.Write(serialized)
}

func (this *NodeConnection) Write(message *BroadcastMessage) {
	message.NodeIdSender = this.selfNodeId

	serialized := Serialize(message)
	this.nativeConnection.Write(serialized)
}

func (this *NodeConnection) WriteBuffered(message *BroadcastMessage) {
	message.NodeIdSender = this.selfNodeId

	serialized := Serialize(message)
	this.nativeWriter.Write(serialized)
}

func (this *NodeConnection) FlushAsync(wait *sync.WaitGroup) {
	go func() {
		this.nativeWriter.Flush()
		wait.Done()
	}()
}

func (this *NodeConnection) GetNodeId() uint32 {
	return this.selfNodeId
}

func ToString(connection []*NodeConnection) string {
	var toReturn string

	for _, connection := range connection {
		toReturn += strconv.Itoa(int(connection.selfNodeId)) + ", "
	}

	return toReturn
}
