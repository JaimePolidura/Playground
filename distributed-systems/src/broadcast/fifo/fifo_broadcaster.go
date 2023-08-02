package fifo

import (
	"distributed-systems/src/broadcast"
	"distributed-systems/src/nodes"
	"distributed-systems/src/nodes/types"
	"fmt"
	"sync/atomic"
)

type FifoBroadcaster struct {
	selfNodeId             uint32
	nodesToPickToBroadcast uint32
	initialTTL             int32
	seqNum                 uint32

	doLogging   bool
	doGossiping bool

	nodeConnectionManager *nodes.ConnectionManager

	broadcastDataByNodeId map[uint32]*FifoNodeBroadcastData

	onBroadcastMessageCallback func(newMessage *nodes.Message)
}

func CreateFifoBroadcaster(nodesToPickToBroadcast uint32, initialTTL int32, selfNodeId uint32) *FifoBroadcaster {
	return &FifoBroadcaster{
		nodesToPickToBroadcast: nodesToPickToBroadcast,
		initialTTL:             initialTTL,
		selfNodeId:             selfNodeId,
		seqNum:                 0,
		doLogging:              true,
		doGossiping:            true,
		broadcastDataByNodeId:  make(map[uint32]*FifoNodeBroadcastData),
	}
}

func (this *FifoBroadcaster) EnableLogging() {
	this.doLogging = true
}

func (this *FifoBroadcaster) DisableLogging() {
	this.doLogging = false
}

func (this *FifoBroadcaster) EnableGossiping() {
	this.doGossiping = true
}

func (this *FifoBroadcaster) DisableGossiping() {
	this.doGossiping = false
}

func (this *FifoBroadcaster) OnStop() {
}

func (this *FifoBroadcaster) Broadcast(message *nodes.Message) {
	this.doBroadcast(message, true)
}

func (this *FifoBroadcaster) OnBroadcastMessage(message *nodes.Message) {
	lastSeqNumDelivered := this.getLastSeqNumDelivered(message.NodeIdOrigin)
	broadcastData := this.broadcastDataByNodeId[message.NodeIdOrigin]
	msgSeqNumbReceived := message.SeqNum

	this.log("[%d] Received broadcast message from node %d with TTL %d and SeqNum %d (Prev: %d). Content: \"%s\"\n",
		this.selfNodeId, message.NodeIdOrigin, message.TTL, message.SeqNum, lastSeqNumDelivered, message.Content)

	if msgSeqNumbReceived > lastSeqNumDelivered && message.NodeIdOrigin != this.selfNodeId {
		broadcastData.AddToBuffer(message)

		for _, messageInBuffer := range broadcastData.RetrieveDeliverableMessages(msgSeqNumbReceived) {
			messageInBuffer = messageInBuffer.RemoveFlag(types.FLAG_BROADCAST)

			this.onBroadcastMessageCallback(messageInBuffer)
		}
	}
	if message.TTL != 0 && this.doGossiping {
		this.doBroadcast(message, false)
	}
}

func (this *FifoBroadcaster) SetNodeConnectionsManager(manager *nodes.ConnectionManager) broadcast.Broadcaster {
	this.nodeConnectionManager = manager
	return this
}

func (this *FifoBroadcaster) SetOnBroadcastMessageCallback(callback func(newMessage *nodes.Message)) broadcast.Broadcaster {
	this.onBroadcastMessageCallback = callback
	return this
}

func (this *FifoBroadcaster) getLastSeqNumDelivered(nodeId uint32) uint32 {
	if prevSeqNum, found := this.broadcastDataByNodeId[nodeId]; found {
		return prevSeqNum.GetLastSeqNumDelivered()
	} else {
		this.broadcastDataByNodeId[nodeId] = CreateFifoNodeBroadcastData()
		return 0
	}
}

func (this *FifoBroadcaster) doBroadcast(message *nodes.Message, firstTime bool) {
	atomic.AddUint32(&this.seqNum, 1)

	message = message.AddFlag(types.FLAG_BROADCAST)

	if firstTime {
		message.SeqNum = this.seqNum
		message.TTL = this.initialTTL
	} else if this.doGossiping {
		message.TTL = message.TTL - 1
	}

	randomNodesId := this.pickRandomNodesId()

	this.log("[%d] Broadcasting \"%s\" with TTL %d and SeqNum %d to %v\n", this.selfNodeId, message.Content,
		message.TTL, message.SeqNum, randomNodesId)

	for _, nodeId := range randomNodesId {
		this.nodeConnectionManager.Send(nodeId, message)
	}
}

func (this *FifoBroadcaster) pickRandomNodesId() []uint32 {
	if this.doGossiping {
		return broadcast.PickRandomNodesId(this.selfNodeId, this.nodesToPickToBroadcast, this.nodeConnectionManager.GetNumberConnections())
	} else {
		return this.nodeConnectionManager.GetNodesId()
	}

}

func (this *FifoBroadcaster) log(format string, other ...any) {
	if this.doLogging {
		fmt.Printf(format, other...)
	}
}
