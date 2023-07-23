package broadcast

import (
	"net"
)

type NodeConnection struct {
	nativeConnection net.Conn
	nodeId           uint32
}

func CreateNodeConnection(nodeId uint32, conn net.Conn) *NodeConnection {
	return &NodeConnection{
		nativeConnection: conn,
		nodeId:           nodeId,
	}
}

func (nodeConnection *NodeConnection) Write(message *Message) {
	serialized := Serialize(message)
	nodeConnection.nativeConnection.Write(serialized)
}
