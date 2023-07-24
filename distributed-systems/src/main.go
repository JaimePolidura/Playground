package main

import (
	"distributed-systems/src/broadcast"
	"distributed-systems/src/broadcast/fifo"
	"fmt"
	"time"
)

func main() {
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

	broadcasterNodes[0].Broadcast("First broadcast :D")

	fmt.Println()
	fmt.Println()
	//
	time.Sleep(time.Second * 5)
	//
	//broadcasterNodes[1].Broadcast("Second broadcast :D")
	//
	//fmt.Println()
	//fmt.Println()
	//
	//time.Sleep(time.Second * 5)
	//
	//broadcasterNodes[7].Broadcast("Third broadcast :D")
	//
	//fmt.Println()
	//fmt.Println()
	//
	//time.Sleep(time.Second * 60)
}
