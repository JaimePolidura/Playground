package broadcast

import "distributed-systems/src/nodes"

type Node struct {
	selfNodeId uint32
	port       uint16

	nodeConnectionsStore *nodes.NodeConnectionsStore
	messageListener      *nodes.MessageListener
	broadcasterNode      *BroadcasterNode

	onBroadcastMessageCallback func(message *nodes.Message)
	onSingleMessageCallback    func(message *nodes.Message)
}

func CreateNode(nodeId uint32, port uint16, broadcaster Broadcaster) *Node {
	nodeConnectionsStore := nodes.CreateNodeConnectionStore()
	broadcaster.SetNodeConnectionsStore(nodeConnectionsStore)

	return &Node{
		selfNodeId:           nodeId,
		port:                 port,
		nodeConnectionsStore: nodeConnectionsStore,
		messageListener:      nodes.CreateMessageListener(nodeId, port),
		broadcasterNode:      CreateBroadcasterNode(nodeId, port, broadcaster, nodeConnectionsStore),
	}
}

func (this *Node) StartListening() {
	this.messageListener.ListenAsync(func(messages []*nodes.Message) {
		for _, message := range messages {
			if message.HasFlag(nodes.BROADCAST_FLAG) {
				this.broadcasterNode.GetBroadcaster().OnBroadcastMessage(messages, this.onBroadcastMessageCallback)
				break
			} else {
				this.onSingleMessageCallback(message)
			}
		}
	})
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
	this.onBroadcastMessageCallback = callback
}

func (this *Node) OnSingleMessage(callback func(message *nodes.Message)) {
	this.onSingleMessageCallback = callback
}

func (this *Node) BroadcastString(content string) {
	this.broadcasterNode.Broadcast(nodes.CreateMessageBroadcast(this.selfNodeId, this.selfNodeId, content))
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

func (this *Node) OpenConnectionToNode(node *Node) {
	if node.selfNodeId != this.selfNodeId {
		this.nodeConnectionsStore.Open(node.selfNodeId)
	}
}

func (this *Node) GetNodeConnectionsStore() *nodes.NodeConnectionsStore {
	return this.nodeConnectionsStore
}

func (this *Node) GetNodeConnections() []*nodes.NodeConnection {
	return this.nodeConnectionsStore.ToArrayNodeConnections()
}
