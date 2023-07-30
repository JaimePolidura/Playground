package broadcast

import (
	"distributed-systems/src/nodes"
	"distributed-systems/src/nodes/types"
)

type BroadcasterNode struct {
	selfNodeId uint32
	port       uint16

	canBroadcast bool

	nodesConnectionManager *nodes.ConnectionManager
	messageListener        *nodes.MessageListener
	broadcaster            Broadcaster

	pendingToBroadcast []*nodes.Message
}

func CreateBroadcasterNode(nodeId uint32, port uint16, broadcaster Broadcaster, nodesConnectionManager *nodes.ConnectionManager) *BroadcasterNode {
	broadcaster.SetNodeConnectionsManager(nodesConnectionManager)

	return &BroadcasterNode{
		selfNodeId:             nodeId,
		port:                   port,
		broadcaster:            broadcaster,
		messageListener:        nodes.CreateMessageListener(nodeId, port),
		nodesConnectionManager: nodesConnectionManager,
		canBroadcast:           true,
		pendingToBroadcast:     make([]*nodes.Message, 0),
	}
}

func (this *BroadcasterNode) Broadcast(message *nodes.Message) {
	if !this.canBroadcast && !message.HasFlag(types.FLAG_BYPASS_ORDERING) {
		this.pendingToBroadcast = append(this.pendingToBroadcast, message)
		return
	}

	this.broadcaster.Broadcast(message)
}

func (this *BroadcasterNode) GetBroadcaster() Broadcaster {
	return this.broadcaster
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

func (this *BroadcasterNode) Stop() {
	this.broadcaster.OnStop()
}
