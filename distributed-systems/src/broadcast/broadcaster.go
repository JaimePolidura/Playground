package broadcast

type Broadcaster interface {
	Broadcast(message *BroadcastMessage)
	OnBroadcastMessage(message []*BroadcastMessage, newMessageCallback func(newMessage *BroadcastMessage))

	SetNodeConnectionsStore(store *NodeConnectionsStore) Broadcaster
}
