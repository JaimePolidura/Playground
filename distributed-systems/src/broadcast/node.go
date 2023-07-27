package broadcast

import (
	"distributed-systems/src/nodes"
)

type Node struct {
	selfNodeId uint32
	port       uint16

	nodeConnectionsStore *nodes.NodeConnectionsStore
	messageListener      *nodes.MessageListener
	broadcasterNode      *BroadcasterNode

	onBroadcastMessageCallback func(message *nodes.Message)

	messageHandlers map[uint8]func(message []*nodes.Message)
}

func CreateNode(nodeId uint32, port uint16, broadcaster Broadcaster) *Node {
	nodeConnectionsStore := nodes.CreateNodeConnectionStore()
	broadcaster.SetNodeConnectionsStore(nodeConnectionsStore)

	node := &Node{
		selfNodeId:           nodeId,
		port:                 port,
		nodeConnectionsStore: nodeConnectionsStore,
		messageListener:      nodes.CreateMessageListener(nodeId, port),
		broadcasterNode:      CreateBroadcasterNode(nodeId, port, broadcaster, nodeConnectionsStore),
		messageHandlers:      make(map[uint8]func(message []*nodes.Message)),
	}

	node.AddMessageHandler(nodes.MESSAGE_BROADCAST, node.broadcastMessageHandler)

	return node
}

func (this *Node) StartListening() {
	this.messageListener.ListenAsync(func(messages []*nodes.Message) {
		this.executeHandlerForMessages(messages)
	})
}

func (this *Node) broadcastMessageHandler(messages []*nodes.Message) {
	for _, message := range messages {
		this.broadcasterNode.GetBroadcaster().OnBroadcastMessage(message)
	}
}

func (this *Node) AddOtherNodeConnection(otherNodeId uint32, port uint32) {
	if otherNodeId != this.selfNodeId {
		this.nodeConnectionsStore.Add(otherNodeId, port, this.selfNodeId)
	}
}

func (this *Node) DisableBroadcast() {
	this.broadcasterNode.DisableBroadcast()
}

func (this *Node) EnableBroadcast() {
	this.broadcasterNode.EnableBroadcast()
}

func (this *Node) OnBroadcastMessage(callback func(message *nodes.Message)) {
	this.broadcasterNode.broadcaster.SetOnBroadcastMessageCallback(callback)
	this.onBroadcastMessageCallback = callback
}

func (this *Node) BroadcastString(content string) {
	this.broadcasterNode.Broadcast(nodes.CreateMessage(
		nodes.WithNodeId(this.selfNodeId),
		nodes.WithType(nodes.MESSAGE_DO_BROADCAST),
		nodes.WithContentString(content)))
}

func (this *Node) Broadcast(message *nodes.Message) {
	this.broadcasterNode.Broadcast(message)
}

func (this *Node) GetBroadcaster() Broadcaster {
	return this.broadcasterNode.GetBroadcaster()
}

func (this *Node) GetNodeId() uint32 {
	return this.selfNodeId
}

func (this *Node) OpenConnectionsToNodes(nodes []*Node) {
	for _, node := range nodes {
		if node.selfNodeId != this.selfNodeId {
			this.nodeConnectionsStore.Open(node.selfNodeId)
		}
	}
}

func (this *Node) AddMessageHandler(typeMessage uint8, handler func(message []*nodes.Message)) {
	this.messageHandlers[typeMessage] = handler
}

func (this *Node) OpenConnectionToNode(node *Node) {
	if node.selfNodeId != this.selfNodeId {
		this.nodeConnectionsStore.Open(node.selfNodeId)
	}
}

func (this *Node) GetNodeConnectionsStore() *nodes.NodeConnectionsStore {
	return this.nodeConnectionsStore
}

func (this *Node) executeHandlerForMessage(message *nodes.Message) {
	this.messageHandlers[message.Type]([]*nodes.Message{message})
}

func (this *Node) executeHandlerForMessages(messages []*nodes.Message) {
	if len(messages) > 0 {
		this.messageHandlers[messages[0].Type](messages)
	}
}

func (this *Node) GetNodeConnections() []*nodes.NodeConnection {
	return this.nodeConnectionsStore.ToArrayNodeConnections()
}
