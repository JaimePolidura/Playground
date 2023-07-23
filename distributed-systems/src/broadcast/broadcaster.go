package broadcast

type Broadcaster interface {
	Broadcast(message *Message)
	OnBroadcastMessage(message *Message, newMessageCallback func(newMessage *Message))

	SetNodeConnectionsStore(store *NodeConnectionsStore) Broadcaster
}
