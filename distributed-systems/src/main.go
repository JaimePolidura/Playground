package main

import (
	"distributed-systems/src/broadcast"
	"distributed-systems/src/broadcast/fifo"
	"distributed-systems/src/broadcast/zab"
	"time"
)

func main() {
	startZab()
	//startFifo()
}

func startZab() {
	nNodes := uint32(4)
	initPort := uint16(1000)
	broadcasterNodes := make([]*broadcast.BroadcasterNode, nNodes)

	for i := uint32(0); i < nNodes; i++ {
		broadcasterNodes[i] = broadcast.CreateBroadcasterNode(i, initPort+uint16(i), zab.CreateZabBroadcaster(i, 0))
		for j := uint32(0); j < nNodes; j++ {
			broadcasterNodes[i].AddOtherNode(j, j+1000)
		}
	}
	for i := uint32(0); i < nNodes; i++ {
		broadcasterNodes[i].StartListening()
	}
	for i := uint32(0); i < nNodes; i++ {
		broadcasterNodes[i].OpenConnectionsToNodes(broadcasterNodes)
	}
	broadcasterNodes[1].Broadcast("Running on zab 1º!")
	broadcasterNodes[1].Broadcast("Running on zab 2º!")
	time.Sleep(time.Second * 5)
}

func startFifo() {
	nNodes := uint32(6)
	ttl := int32(3)
	nodesToPick := uint32(2)
	initPort := uint16(1000)

	broadcasterNodes := make([]*broadcast.BroadcasterNode, nNodes)

	for i := uint32(0); i < nNodes; i++ {
		broadcasterNodes[i] = broadcast.CreateBroadcasterNode(i, initPort+uint16(i), fifo.CreateFifoBroadcaster(nodesToPick, ttl, i))

		for j := uint32(0); j < nNodes; j++ {
			broadcasterNodes[i].AddOtherNode(j, j+1000)
		}
	}
	for i := uint32(0); i < nNodes; i++ {
		broadcasterNodes[i].StartListening()
	}
	for i := uint32(0); i < nNodes; i++ {
		broadcasterNodes[i].OpenConnectionsToNodes(broadcasterNodes)
	}

	broadcasterNodes[0].Broadcast("Running on fifo :D")
	time.Sleep(time.Second * 5)

}
