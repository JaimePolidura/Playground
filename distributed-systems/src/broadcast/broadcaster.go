package broadcast

type Broadcaster interface {
	Broadcast(message *Message)
	OnBroadcastMessage(message *Message)

	SetNodeConnectionsStore(store *NodeConnectionsStore) Broadcaster
}
