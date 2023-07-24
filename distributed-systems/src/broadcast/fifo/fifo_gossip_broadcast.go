package fifo

import (
	"distributed-systems/src/broadcast"
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"
)

type FifoGossipBroadcaster struct {
	selfNodeId             uint32
	nodesToPickToBroadcast uint32
	initialTTL             int32
	seqNum                 uint32

	nodeConnectionsStore *broadcast.NodeConnectionsStore

	broadcastDataByNodeId map[uint32]*FifoGossipNodeBroadcastData
}

func CreateFifoGossip(nodesToPickToBroadcast uint32, initialTTL int32, selfNodeId uint32) *FifoGossipBroadcaster {
	return &FifoGossipBroadcaster{
		nodesToPickToBroadcast: nodesToPickToBroadcast,
		initialTTL:             initialTTL,
		selfNodeId:             selfNodeId,
		seqNum:                 0,
		broadcastDataByNodeId:  make(map[uint32]*FifoGossipNodeBroadcastData),
	}
}

func (broadcaster *FifoGossipBroadcaster) SetNodeConnectionsStore(store *broadcast.NodeConnectionsStore) broadcast.Broadcaster {
	broadcaster.nodeConnectionsStore = store
	return broadcaster
}

func (broadcaster *FifoGossipBroadcaster) OnBroadcastMessage(message *broadcast.Message, newMessageCallback func(newMessage *broadcast.Message)) {
	lastSeqNumDelivered := broadcaster.getLastSeqNumDelivered(message.NodeId)
	broadcastData := broadcaster.broadcastDataByNodeId[message.NodeId]
	msgSeqNumbReceived := message.SeqNum

	fmt.Printf("[%d] Recieved broadcast message from node %d with TTL %d and SeqNum %d (Prev: %d). Content: \"%s\"\n",
		broadcaster.selfNodeId, message.NodeId, message.TTL, message.SeqNum, lastSeqNumDelivered, message.Content)

	if msgSeqNumbReceived > lastSeqNumDelivered && message.TTL != 0 {
		broadcastData.AddToBuffer(message)

		for _, messageInBuffer := range broadcastData.GetDeliverableMessages(msgSeqNumbReceived) {
			broadcastData.lastSeqNumDelivered = messageInBuffer.SeqNum
			broadcaster.doBroadcast(messageInBuffer, false)

			newMessageCallback(messageInBuffer)
		}
	}
}

func (broadcaster *FifoGossipBroadcaster) getLastSeqNumDelivered(nodeId uint32) uint32 {
	if prevSeqNum, found := broadcaster.broadcastDataByNodeId[nodeId]; found {
		return prevSeqNum.GetLastSeqNumDelivered()
	} else {
		broadcaster.broadcastDataByNodeId[nodeId] = CreateFifoGossipNodeBroadcastData()
		return 0
	}
}

func (broadcaster *FifoGossipBroadcaster) Broadcast(message *broadcast.Message) {
	broadcaster.doBroadcast(message, true)
}

func (broadcaster *FifoGossipBroadcaster) doBroadcast(message *broadcast.Message, firstTime bool) {
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

func (broadcaster *FifoGossipBroadcaster) pickRandomConnections() []*broadcast.NodeConnection {
	randomNodesIdToBroadcast := broadcaster.pickRandomNodesId()
	randomNodes := make([]*broadcast.NodeConnection, 0)

	for i := 0; i < int(broadcaster.nodesToPickToBroadcast); i++ {
		randomNodeId := randomNodesIdToBroadcast[i]
		randomNode := broadcaster.nodeConnectionsStore.Get(randomNodeId)

		randomNodes = append(randomNodes, randomNode)
	}

	return randomNodes
}

func (broadcaster *FifoGossipBroadcaster) pickRandomNodesId() []uint32 {
	random := make([]uint32, 0)

	for uint32(len(random)) != broadcaster.nodesToPickToBroadcast {
		rand.Seed(time.Now().UnixNano())
		randomNodeId := uint32(rand.Intn(broadcaster.nodeConnectionsStore.Size()))

		if randomNodeId != broadcaster.selfNodeId && !contains(random, randomNodeId) {
			random = append(random, randomNodeId)
		}
	}

	return random
}

func contains(arr []uint32, toCheck uint32) bool {
	for _, value := range arr {
		if value == toCheck {
			return true
		}
	}

	return false
}
