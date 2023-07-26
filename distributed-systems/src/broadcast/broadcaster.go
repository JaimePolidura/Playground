package broadcast

import (
	"distributed-systems/src/nodes"
)

type Broadcaster interface {
	Broadcast(message *nodes.Message)
	OnBroadcastMessage(message []*nodes.Message, newMessageCallback func(newMessage *nodes.Message))

	SetNodeConnectionsStore(store *nodes.NodeConnectionsStore) Broadcaster
}
