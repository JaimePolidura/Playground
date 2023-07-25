package fifo

import (
	"distributed-systems/src/broadcast"
	"fmt"
	"sync/atomic"
)

type FifoBroadcaster struct {
	selfNodeId             uint32
	nodesToPickToBroadcast uint32
	initialTTL             int32
	seqNum                 uint32

	nodeConnectionsStore *broadcast.NodeConnectionsStore

	broadcastDataByNodeId map[uint32]*FifoNodeBroadcastData
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

func (this *FifoBroadcaster) SetNodeConnectionsStore(store *broadcast.NodeConnectionsStore) broadcast.Broadcaster {
	this.nodeConnectionsStore = store
	return this
}

func (this *FifoBroadcaster) OnBroadcastMessage(messages []*broadcast.BroadcastMessage, newMessageCallback func(newMessage *broadcast.BroadcastMessage)) {
	message := messages[0]
	lastSeqNumDelivered := this.getLastSeqNumDelivered(message.NodeIdOrigin)
	broadcastData := this.broadcastDataByNodeId[message.NodeIdOrigin]
	msgSeqNumbReceived := message.SeqNum

	fmt.Printf("[%d] Recieved broadcast message from node %d with TTL %d and SeqNum %d (Prev: %d). Content: \"%s\"\n",
		this.selfNodeId, message.NodeIdOrigin, message.TTL, message.SeqNum, lastSeqNumDelivered, message.Content)

	if msgSeqNumbReceived > lastSeqNumDelivered && message.TTL != 0 && message.NodeIdOrigin != this.selfNodeId {
		broadcastData.AddToBuffer(message)

		for _, messageInBuffer := range broadcastData.GetDeliverableMessages(msgSeqNumbReceived) {
			broadcastData.lastSeqNumDelivered = messageInBuffer.SeqNum
			this.doBroadcast(messageInBuffer, false)

			newMessageCallback(messageInBuffer)
		}
	}
}

func (this *FifoBroadcaster) getLastSeqNumDelivered(nodeId uint32) uint32 {
	if prevSeqNum, found := this.broadcastDataByNodeId[nodeId]; found {
		return prevSeqNum.GetLastSeqNumDelivered()
	} else {
		this.broadcastDataByNodeId[nodeId] = CreateFifoNodeBroadcastData()
		return 0
	}
}

func (this *FifoBroadcaster) Broadcast(message *broadcast.BroadcastMessage) {
	this.doBroadcast(message, true)
}

func (this *FifoBroadcaster) doBroadcast(message *broadcast.BroadcastMessage, firstTime bool) {
	atomic.AddUint32(&this.seqNum, 1)

	if firstTime {
		message.SeqNum = this.seqNum
		message.TTL = this.initialTTL
	} else {
		message.TTL = message.TTL - 1
	}

	randomNodesConnections := this.pickRandomConnections()

	fmt.Printf("[%d] Broadcasting \"%s\" with TTL %d and SeqNum %d to %s\n", this.selfNodeId, message.Content,
		message.TTL, message.SeqNum, broadcast.ToString(randomNodesConnections))

	for i := 0; i < len(randomNodesConnections); i++ {
		nodeConnection := randomNodesConnections[i]
		nodeConnection.Write(message)
	}
}

func (this *FifoBroadcaster) pickRandomConnections() []*broadcast.NodeConnection {
	return broadcast.PickRandomNodesConnections(this.nodeConnectionsStore,
		broadcast.PickRandomNodesId(this.selfNodeId, this.nodesToPickToBroadcast, this.nodeConnectionsStore.Size()))
}
