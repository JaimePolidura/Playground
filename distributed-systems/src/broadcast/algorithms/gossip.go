package algorithms

import (
	"distributed-systems/src/broadcast"
	"math/rand"
	"sync/atomic"
	"time"
)

type GossipBroadcaster struct {
	nodesToPickToBroadcast uint32
	initialTTL             uint32
	seqNum                 uint32

	nodeConnectionsStore *broadcast.NodeConnectionsStore

	seqNumsByNodesId map[uint32]uint32
}

func CreateGossip(nodesToPickToBroadcast uint32, initialTTL uint32) GossipBroadcaster {
	return GossipBroadcaster{
		nodesToPickToBroadcast: nodesToPickToBroadcast,
		initialTTL:             initialTTL,
		seqNum:                 0,
		seqNumsByNodesId:       make(map[uint32]uint32),
	}
}

func (broadcaster GossipBroadcaster) SetNodeConnectionsStore(store *broadcast.NodeConnectionsStore) broadcast.Broadcaster {
	broadcaster.nodeConnectionsStore = store
	return broadcaster
}

func (broadcaster GossipBroadcaster) Broadcast(message *broadcast.Message) {
	atomic.AddUint32(&broadcaster.seqNum, 1)

	message.TTL = message.TTL - 1

	randomNodesConnections := broadcaster.pickRandomConnections()

	for i := 0; i < len(randomNodesConnections); i++ {
		nodeConnection := randomNodesConnections[i]
		nodeConnection.Write(message)
	}
}

func (broadcaster GossipBroadcaster) OnBroadcastMessage(message *broadcast.Message) {
	prevSeqNum := broadcaster.getPrevSeqNum(message.NodeId)
	msgSeqNum := message.SeqNum

	if msgSeqNum > prevSeqNum {
		broadcaster.seqNumsByNodesId[message.NodeId] = message.SeqNum
		broadcaster.Broadcast(message)
	}
}

func (broadcaster GossipBroadcaster) getPrevSeqNum(nodeId uint32) uint32 {
	if prevSeqNum, found := broadcaster.seqNumsByNodesId[nodeId]; found {
		return prevSeqNum
	} else {
		return 0
	}
}

func (broadcaster GossipBroadcaster) pickRandomConnections() []broadcast.NodeConnection {
	randomNodesIndexToBroadcast := pickRandomIndexes(broadcaster.nodesToPickToBroadcast)
	randomNodesConnectionArray := broadcaster.nodeConnectionsStore.ToArrayNodeConnections()

	for i := 0; i < int(broadcaster.nodesToPickToBroadcast); i++ {
		randomNodeIdIndex := randomNodesIndexToBroadcast[i]
		randomNodeConnection := randomNodesConnectionArray[randomNodeIdIndex]

		randomNodesConnectionArray = append(randomNodesConnectionArray, randomNodeConnection)
	}

	return randomNodesConnectionArray
}

func pickRandomIndexes(toPickRandomBroadcast uint32) []int {
	random := make([]int, toPickRandomBroadcast)
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < int(toPickRandomBroadcast); i++ {
		random[i] = rand.Intn(int(toPickRandomBroadcast))
	}

	return random
}
