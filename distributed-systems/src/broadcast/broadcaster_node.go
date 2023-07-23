package broadcast

type BroadcasterNode struct {
	nodeId uint32
	port   uint16

	nodeConnectionsStore *NodeConnectionsStore
	messageListener      *MessageListener
	broadcaster          Broadcaster
}

func CreateBroadcasterNode(nodeId uint32, port uint16, broadcaster Broadcaster) *BroadcasterNode {
	nodeConnectionsStore := CreateNodeConnectionStore()
	broadcaster.SetNodeConnectionsStore(nodeConnectionsStore)

	return &BroadcasterNode{
		nodeId:               nodeId,
		port:                 port,
		broadcaster:          broadcaster,
		messageListener:      CreateMessageListener(nodeId, port),
		nodeConnectionsStore: nodeConnectionsStore,
	}
}

func (broadcasterNode BroadcasterNode) StartListening() {
	broadcasterNode.messageListener.ListenAsync(func(message *Message) {
		broadcasterNode.broadcaster.OnBroadcastMessage(message)
	})
}
