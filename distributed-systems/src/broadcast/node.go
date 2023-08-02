package broadcast

import (
	"distributed-systems/src/nodes"
	"distributed-systems/src/nodes/types"
)

type Node struct {
	selfNodeId uint32

	connectionManager *nodes.ConnectionManager
	broadcasterNode   *BroadcasterNode

	messageHandlers map[uint8]func(message *nodes.Message)
}

func CreateNode(nodeId uint32, port uint16, broadcaster Broadcaster) *Node {
	nodeConnectionManager := nodes.CreateConnectionManager(nodeId, port)
	broadcaster.SetNodeConnectionsManager(nodeConnectionManager)

	node := &Node{
		selfNodeId:        nodeId,
		broadcasterNode:   CreateBroadcasterNode(nodeId, port, broadcaster, nodeConnectionManager),
		messageHandlers:   make(map[uint8]func(message *nodes.Message)),
		connectionManager: nodeConnectionManager,
	}

	node.broadcasterNode.GetBroadcaster().SetOnBroadcastMessageCallback(node.executeHandlerForMessage)
	node.AddMessageHandler(types.MESSAGE_NODE_STOPPED, node.nodeStoppedMessageHandler)

	return node
}

func (this *Node) AddMessagesTypesToListenBroadcast(messageTypes ...uint8) {
	for _, messageType := range messageTypes {
		this.AddMessageHandler(messageType, this.broadcastMessageHandler)
	}
}

func (this *Node) Stop() {
	this.connectionManager.Stop()
	this.broadcasterNode.Stop()
	this.broadcasterNode.Broadcast(nodes.CreateMessage(
		nodes.WithNodeId(this.selfNodeId),
		nodes.WithType(types.MESSAGE_NODE_STOPPED),
		nodes.WithFlags(types.FLAG_BYPASS_LEADER, types.FLAG_BYPASS_ORDERING)))
}

func (this *Node) StartListeningAsync() {
	go func() {
		this.connectionManager.StartListeningAsync()

		for message := range this.connectionManager.NewMessage {
			if message.HasFlag(types.FLAG_BROADCAST) {
				this.broadcasterNode.GetBroadcaster().OnBroadcastMessage(message)
			} else {
				this.executeHandlerForMessage(message)
			}
		}
	}()
}

func (this *Node) broadcastMessageHandler(message *nodes.Message) {
	this.broadcasterNode.GetBroadcaster().OnBroadcastMessage(message)
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
	this.AddMessageHandler(types.MESSAGE_BROADCAST, callback)
}

func (this *Node) BroadcastString(content string, messageType uint8) {
	this.broadcasterNode.Broadcast(nodes.CreateMessage(
		nodes.WithNodeId(this.selfNodeId),
		nodes.WithType(messageType),
		nodes.WithFlags(types.FLAG_BROADCAST),
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

func (this *Node) AddMessageHandler(typeMessage uint8, handler func(message *nodes.Message)) {
	this.messageHandlers[typeMessage] = handler
}

func (this *Node) executeHandlerForMessage(message *nodes.Message) {
	callback, contained := this.messageHandlers[message.Type]
	if contained {
		callback(message)
	}
}

func (this *Node) GetConnectionManager() *nodes.ConnectionManager {
	return this.connectionManager
}

func (this *Node) nodeStoppedMessageHandler(message *nodes.Message) {
	nodeRemovedId := message.NodeIdSender
	this.connectionManager.StopByNodeId(nodeRemovedId)
}
