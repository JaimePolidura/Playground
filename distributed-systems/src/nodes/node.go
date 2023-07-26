package nodes

import (
	"distributed-systems/src/broadcast"
)

type Node struct {
	selfNodeId uint32
	port       uint16

	nodeConnectionsStore *NodeConnectionsStore
	messageListener      *MessageListener
	broadcasterNode      *broadcast.BroadcasterNode

	onBroadcastMessageCallback func(message *Message)
	onSingleMessageCallback    func(message *Message)
}

func CreateNode(nodeId uint32, port uint16, broadcaster broadcast.Broadcaster) *Node {
	nodeConnectionsStore := CreateNodeConnectionStore()
	broadcaster.SetNodeConnectionsStore(nodeConnectionsStore)

	return &Node{
		selfNodeId:           nodeId,
		port:                 port,
		nodeConnectionsStore: nodeConnectionsStore,
		messageListener:      CreateMessageListener(nodeId, port),
		broadcasterNode:      broadcast.CreateBroadcasterNode(nodeId, port, broadcaster),
	}
}

func (this *Node) StartListening() {
	this.messageListener.ListenAsync(func(messages []*Message) {
		for _, message := range messages {
			if message.HasFlag(BROADCAST_FLAG) {
				this.broadcasterNode.GetBroadcaster().OnBroadcastMessage(messages, this.onBroadcastMessageCallback)
				break
			} else {
				this.onSingleMessageCallback(message)
			}
		}
	})
}

func (this *Node) DisableBroadcast() {
	this.broadcasterNode.DisableBroadcast()
}

func (this *Node) EnableBroadcast() {
	this.broadcasterNode.EnableBroadcast()
}

func (this *Node) OnBroadcastMessage(callback func(message *Message)) {
	this.onBroadcastMessageCallback = callback
}

func (this *Node) OnSingleMessage(callback func(message *Message)) {
	this.onSingleMessageCallback = callback
}

func (this *Node) BroadcastString(content string) {
	this.broadcasterNode.Broadcast(CreateMessageBroadcast(this.selfNodeId, this.selfNodeId, content))
}

func (this *Node) Broadcast(message *Message) {
	this.broadcasterNode.Broadcast(message)
}

func (this *Node) GetBroadcaster() broadcast.Broadcaster {
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

func (this *Node) GetNodeConnectionsStore() *NodeConnectionsStore {
	return this.nodeConnectionsStore
}

func (this *Node) GetNodeConnections() []*NodeConnection {
	return this.nodeConnectionsStore.ToArrayNodeConnections()
}
