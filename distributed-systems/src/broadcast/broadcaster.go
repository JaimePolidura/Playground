package broadcast

import (
	"distributed-systems/src/nodes"
)

type Broadcaster interface {
	Broadcast(message *nodes.Message)
	OnBroadcastMessage(newMessage *nodes.Message)

	SetNodeConnectionsManager(nodesConnectionManager *nodes.ConnectionManager) Broadcaster
	SetOnBroadcastMessageCallback(callback func(newMessage *nodes.Message)) Broadcaster
}
