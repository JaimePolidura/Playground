package broadcast

import (
	"math/rand"
	"time"
)

func PickRandomNodesConnections(nodeConnectionsStore *NodeConnectionsStore, randomNodesId []uint32) []*NodeConnection {
	randomNodes := make([]*NodeConnection, 0)

	for i := 0; i < len(randomNodesId); i++ {
		randomNodeId := randomNodesId[i]
		randomNode := nodeConnectionsStore.Get(randomNodeId)

		randomNodes = append(randomNodes, randomNode)
	}

	return randomNodes
}

func PickRandomNodesId(selfNodeId uint32, nodesToPickToBroadcast uint32, numberNodes uint32) []uint32 {
	random := make([]uint32, 0)

	for uint32(len(random)) != nodesToPickToBroadcast {
		rand.Seed(time.Now().UnixNano())
		randomNodeId := uint32(rand.Intn(int(numberNodes)))

		if randomNodeId != selfNodeId && !contains(random, randomNodeId) {
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
