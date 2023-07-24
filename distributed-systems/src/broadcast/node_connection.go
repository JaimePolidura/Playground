package broadcast

import (
	"fmt"
	"net"
	"strconv"
)

type NodeConnection struct {
	nativeConnection net.Conn
	selfNodeId       uint32
	port             uint32
}

func CreateNodeConnection(nodeId uint32, port uint32) *NodeConnection {
	return &NodeConnection{
		selfNodeId: nodeId,
		port:       port,
	}
}

func (nodeConnection *NodeConnection) Open() {
	conn, err := net.Dial("tcp", "127.0.0.1:"+strconv.Itoa(int(nodeConnection.port)))

	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	nodeConnection.nativeConnection = conn
}

func (nodeConnection *NodeConnection) WriteAll(messages []*Message) {
	for _, message := range messages {
		message.NodeIdSender = nodeConnection.selfNodeId
	}

	serialized := SerializeAll(messages)
	nodeConnection.nativeConnection.Write(serialized)
}

func (nodeConnection *NodeConnection) Write(message *Message) {
	message.NodeIdSender = nodeConnection.selfNodeId

	serialized := Serialize(message)
	nodeConnection.nativeConnection.Write(serialized)
}

func (nodeConnection *NodeConnection) GetNodeId() uint32 {
	return nodeConnection.selfNodeId
}

func ToString(connection []*NodeConnection) string {
	var toReturn string

	for _, connection := range connection {
		toReturn += strconv.Itoa(int(connection.selfNodeId)) + ", "
	}

	return toReturn
}
