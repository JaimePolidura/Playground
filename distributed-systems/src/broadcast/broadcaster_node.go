package broadcast

import (
	"distributed-systems/src/nodes"
)

type BroadcasterNode struct {
	selfNodeId uint32
	port       uint16

	canBroadcast bool

	nodeConnectionsStore *nodes.NodeConnectionsStore
	messageListener      *nodes.MessageListener
	broadcaster          Broadcaster

	pendingToBroadcast []*nodes.Message
}

func CreateBroadcasterNode(nodeId uint32, port uint16, broadcaster Broadcaster, store *nodes.NodeConnectionsStore) *BroadcasterNode {
	broadcaster.SetNodeConnectionsStore(store)

	return &BroadcasterNode{
		selfNodeId:           nodeId,
		port:                 port,
		broadcaster:          broadcaster,
		messageListener:      nodes.CreateMessageListener(nodeId, port),
		nodeConnectionsStore: store,
		canBroadcast:         true,
		pendingToBroadcast:   make([]*nodes.Message, 0),
	}
}

func (this *BroadcasterNode) Broadcast(message *nodes.Message) {
	if !this.canBroadcast {
		this.pendingToBroadcast = append(this.pendingToBroadcast, message)
		return
	}

	this.broadcaster.Broadcast(message)
}

func (this *BroadcasterNode) GetBroadcaster() Broadcaster {
	return this.broadcaster
}

func (this *BroadcasterNode) OpenConnectionsToNodes(nodes []*BroadcasterNode) {
	for _, node := range nodes {
		if node.selfNodeId != this.selfNodeId {
			this.nodeConnectionsStore.Open(node.selfNodeId)
		}
	}
}

func (this *BroadcasterNode) DisableBroadcast() {
	this.canBroadcast = false
}

func (this *BroadcasterNode) EnableBroadcast() {
	this.canBroadcast = true
	for _, message := range this.pendingToBroadcast {
		this.broadcaster.Broadcast(message)
	}

	this.pendingToBroadcast = []*nodes.Message{}
}
