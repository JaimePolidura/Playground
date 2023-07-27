package broadcast

import (
	"distributed-systems/src/nodes"
)

type Broadcaster interface {
	Broadcast(message *nodes.Message)
	OnBroadcastMessage(newMessage *nodes.Message)

	SetNodeConnectionsStore(store *nodes.NodeConnectionsStore) Broadcaster
	SetOnBroadcastMessageCallback(callback func(newMessage *nodes.Message)) Broadcaster
}
