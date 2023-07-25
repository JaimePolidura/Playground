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

func (broadcaster *FifoBroadcaster) SetNodeConnectionsStore(store *broadcast.NodeConnectionsStore) broadcast.Broadcaster {
	broadcaster.nodeConnectionsStore = store
	return broadcaster
}

func (broadcaster *FifoBroadcaster) OnBroadcastMessage(messages []*broadcast.BroadcastMessage, newMessageCallback func(newMessage *broadcast.BroadcastMessage)) {
	message := messages[0]
	lastSeqNumDelivered := broadcaster.getLastSeqNumDelivered(message.NodeIdOrigin)
	broadcastData := broadcaster.broadcastDataByNodeId[message.NodeIdOrigin]
	msgSeqNumbReceived := message.SeqNum

	fmt.Printf("[%d] Recieved broadcast message from node %d with TTL %d and SeqNum %d (Prev: %d). Content: \"%s\"\n",
		broadcaster.selfNodeId, message.NodeIdOrigin, message.TTL, message.SeqNum, lastSeqNumDelivered, message.Content)

	if msgSeqNumbReceived > lastSeqNumDelivered && message.TTL != 0 {
		broadcastData.AddToBuffer(message)

		for _, messageInBuffer := range broadcastData.GetDeliverableMessages(msgSeqNumbReceived) {
			broadcastData.lastSeqNumDelivered = messageInBuffer.SeqNum
			broadcaster.doBroadcast(messageInBuffer, false)

			newMessageCallback(messageInBuffer)
		}
	}
}

func (broadcaster *FifoBroadcaster) getLastSeqNumDelivered(nodeId uint32) uint32 {
	if prevSeqNum, found := broadcaster.broadcastDataByNodeId[nodeId]; found {
		return prevSeqNum.GetLastSeqNumDelivered()
	} else {
		broadcaster.broadcastDataByNodeId[nodeId] = CreateFifoNodeBroadcastData()
		return 0
	}
}

func (broadcaster *FifoBroadcaster) Broadcast(message *broadcast.BroadcastMessage) {
	broadcaster.doBroadcast(message, true)
}

func (broadcaster *FifoBroadcaster) doBroadcast(message *broadcast.BroadcastMessage, firstTime bool) {
	atomic.AddUint32(&broadcaster.seqNum, 1)

	if firstTime {
		message.SeqNum = broadcaster.seqNum
		message.TTL = broadcaster.initialTTL
	} else {
		message.TTL = message.TTL - 1
	}

	randomNodesConnections := broadcaster.pickRandomConnections()

	fmt.Printf("[%d] Broadcasting \"%s\" with TTL %d and SeqNum %d to %s\n", broadcaster.selfNodeId, message.Content,
		message.TTL, message.SeqNum, broadcast.ToString(randomNodesConnections))

	for i := 0; i < len(randomNodesConnections); i++ {
		nodeConnection := randomNodesConnections[i]
		nodeConnection.Write(message)
	}
}

func (broadcaster *FifoBroadcaster) pickRandomConnections() []*broadcast.NodeConnection {
	return broadcast.PickRandomNodesConnections(broadcaster.nodeConnectionsStore,
		broadcast.PickRandomNodesId(broadcaster.selfNodeId, broadcaster.nodesToPickToBroadcast, broadcaster.nodeConnectionsStore.Size()))
}
