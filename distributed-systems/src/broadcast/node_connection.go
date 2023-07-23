package broadcast

import (
	"fmt"
	"net"
	"strconv"
)

type NodeConnection struct {
	nativeConnection net.Conn
	nodeId           uint32
	port             uint32
}

func CreateNodeConnection(nodeId uint32, port uint32) *NodeConnection {
	return &NodeConnection{
		nodeId: nodeId,
		port:   port,
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

func (nodeConnection *NodeConnection) Write(message *Message) {
	serialized := Serialize(message)
	nodeConnection.nativeConnection.Write(serialized)
}

func ToString(connection []*NodeConnection) string {
	var toReturn string

	for _, connection := range connection {
		toReturn += strconv.Itoa(int(connection.nodeId)) + ", "
	}

	return toReturn
}
