package broadcast

import "distributed-systems/src/broadcast/algorithms"

func main() {
	nNodes := 12
	ttl := 6
	nodesToPick := 3
	initPort := 1000

	var broadcasterNode []*BroadcasterNode

	for i := 0; i < nNodes; i++ {
		broadcasterNode[i] = CreateBroadcasterNode(uint32(i), uint16(initPort+10), algorithms.CreateGossip(uint32(nodesToPick), uint32(ttl)))
	}
}
