package nodes

import (
	"fmt"
	"net"
	"strconv"
)

type NodeConnection struct {
	nativeConnection net.Conn

	selfNodeId uint32 //References nodeId of the node that holds the connection

	nodeId uint32 //References nodeId of nativeConnections
	port   uint32
}

func CreateNodeConnection(nodeId uint32, port uint32, selfNodeId uint32) *NodeConnection {
	return &NodeConnection{
		selfNodeId: selfNodeId,
		nodeId:     nodeId,
		port:       port,
	}
}

func (this *NodeConnection) Open() {
	conn, err := net.Dial("tcp", "127.0.0.1:"+strconv.Itoa(int(this.port)))

	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	this.nativeConnection = conn
}

func (this *NodeConnection) WriteAll(messages []*Message) {
	for _, message := range messages {
		message.NodeIdSender = this.selfNodeId
	}

	serialized := SerializeAll(messages)
	this.nativeConnection.Write(serialized)
}

func (this *NodeConnection) Write(message *Message) {
	message.NodeIdSender = this.selfNodeId

	serialized := Serialize(message)
	this.nativeConnection.Write(serialized)
}

func (this *NodeConnection) GetNodeId() uint32 {
	return this.nodeId
}

func ToString(connection []*NodeConnection) string {
	var toReturn string

	for _, connection := range connection {
		toReturn += strconv.Itoa(int(connection.nodeId)) + ", "
	}

	return toReturn
}
