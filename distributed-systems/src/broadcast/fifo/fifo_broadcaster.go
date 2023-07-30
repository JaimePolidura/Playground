package fifo

import (
	"distributed-systems/src/broadcast"
	"distributed-systems/src/nodes"
	"fmt"
	"sync/atomic"
)

type FifoBroadcaster struct {
	selfNodeId             uint32
	nodesToPickToBroadcast uint32
	initialTTL             int32
	seqNum                 uint32

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
		broadcastDataByNodeId:  make(map[uint32]*FifoNodeBroadcastData),
	}
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

	fmt.Printf("[%d] Received broadcast message from node %d with TTL %d and SeqNum %d (Prev: %d). Content: \"%s\"\n",
		this.selfNodeId, message.NodeIdOrigin, message.TTL, message.SeqNum, lastSeqNumDelivered, message.Content)

	if msgSeqNumbReceived > lastSeqNumDelivered && message.NodeIdOrigin != this.selfNodeId {
		broadcastData.AddToBuffer(message)

		for _, messageInBuffer := range broadcastData.RetrieveDeliverableMessages(msgSeqNumbReceived) {
			this.onBroadcastMessageCallback(messageInBuffer)
		}
	}
	if message.TTL != 0 {
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

	if firstTime {
		message.SeqNum = this.seqNum
		message.TTL = this.initialTTL
	} else {
		message.TTL = message.TTL - 1
	}

	randomNodesId := this.pickRandomNodesId()

	fmt.Printf("[%d] Broadcasting \"%s\" with TTL %d and SeqNum %d to %v\n", this.selfNodeId, message.Content,
		message.TTL, message.SeqNum, randomNodesId)

	for _, nodeId := range randomNodesId {
		this.nodeConnectionManager.Send(nodeId, message)
	}
}

func (this *FifoBroadcaster) pickRandomNodesId() []uint32 {
	return broadcast.PickRandomNodesId(this.selfNodeId, this.nodesToPickToBroadcast, this.nodeConnectionManager.GetNumberConnections())
}
