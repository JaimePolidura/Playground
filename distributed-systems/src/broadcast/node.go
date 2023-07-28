package broadcast

import (
	"distributed-systems/src/nodes"
)

type Node struct {
	selfNodeId uint32
	port       uint16

	connectionManager *nodes.ConnectionManager
	broadcasterNode   *BroadcasterNode

	messageHandlers map[uint8]func(message []*nodes.Message)

	onBroadcastMessageCallback func(message *nodes.Message)
}

func CreateNode(nodeId uint32, port uint16, broadcaster Broadcaster) *Node {
	nodeConnectionManager := nodes.CreateConnectionManager(nodeId, port)
	broadcaster.SetNodeConnectionsManager(nodeConnectionManager)

	node := &Node{
		selfNodeId:        nodeId,
		port:              port,
		broadcasterNode:   CreateBroadcasterNode(nodeId, port, broadcaster, nodeConnectionManager),
		messageHandlers:   make(map[uint8]func(message []*nodes.Message)),
		connectionManager: nodeConnectionManager,
	}

	node.AddMessageHandler(nodes.MESSAGE_BROADCAST, node.broadcastMessageHandler)

	return node
}

func (this *Node) Stop() {
	this.connectionManager.Stop()
}

func (this *Node) StartListeningAsync() {
	go func() {
		this.connectionManager.StartListeningAsync()

		for message := range this.connectionManager.NewMessage {
			this.executeHandlerForMessage(message)
		}
	}()
}

func (this *Node) broadcastMessageHandler(messages []*nodes.Message) {
	for _, message := range messages {
		this.broadcasterNode.GetBroadcaster().OnBroadcastMessage(message)
	}
}

func (this *Node) AddOtherNodeConnection(otherNodeId uint32, port uint32) {
	if otherNodeId != this.selfNodeId {
		this.connectionManager.Add(otherNodeId, port)
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

func (this *Node) BroadcastString(content string, messageType uint8) {
	this.broadcasterNode.Broadcast(nodes.CreateMessage(
		nodes.WithNodeId(this.selfNodeId),
		nodes.WithType(messageType),
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

func (this *Node) AddMessageHandler(typeMessage uint8, handler func(message []*nodes.Message)) {
	this.messageHandlers[typeMessage] = handler
}

func (this *Node) executeHandlerForMessage(message *nodes.Message) {
	callback, contained := this.messageHandlers[message.Type]
	if contained {
		callback([]*nodes.Message{message})
	}
}

func (this *Node) executeHandlerForMessages(messages []*nodes.Message) {
	callback, contained := this.messageHandlers[messages[0].Type]
	if len(messages) > 0 && contained {
		callback(messages)
	}
}

func (this *Node) GetConnectionManager() *nodes.ConnectionManager {
	return this.connectionManager
}
