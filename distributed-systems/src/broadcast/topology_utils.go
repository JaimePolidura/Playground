package broadcast

import (
	"math/rand"
	"time"
)

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
