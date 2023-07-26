package main

import (
	"distributed-systems/src/broadcast"
	"distributed-systems/src/broadcast/fifo"
	"distributed-systems/src/nodes"
	"fmt"
	"time"
)

func main() {
	//startZab()
	startFifo()
}

func startZab() {
	//nNodes := uint32(4)
	//initPort := uint16(1000)
	//broadcasterNodes := make([]*broadcast.BroadcasterNode, nNodes)
	//
	//for nodeId := uint32(0); nodeId < nNodes; nodeId++ {
	//	broadcasterNodes[nodeId] = broadcast.CreateBroadcasterNode(nodeId, initPort+uint16(nodeId), zab.CreateZabBroadcaster(nodeId, 0))
	//	for otherNodeId := uint32(0); otherNodeId < nNodes; otherNodeId++ {
	//		broadcasterNodes[nodeId].AddOtherNode(otherNodeId, otherNodeId+1000)
	//	}
	//}
	//for nodeId := uint32(0); nodeId < nNodes; nodeId++ {
	//	broadcasterNodes[nodeId].StartListening()
	//}
	//for nodeId := uint32(0); nodeId < nNodes; nodeId++ {
	//	broadcasterNodes[nodeId].OpenConnectionsToNodes(broadcasterNodes)
	//}
	//broadcasterNodes[1].BroadcastString("Running on zab 1º!")
	//broadcasterNodes[1].BroadcastString("Running on zab 2º!")
	//time.Sleep(time.Second * 5)
}

func startFifo() {
	nNodes := uint32(6)
	ttl := int32(3)
	nodesToPick := uint32(2)
	initPort := uint16(1000)

	broadcasterNodes := make([]*broadcast.Node, nNodes)

	for i := uint32(0); i < nNodes; i++ {
		broadcasterNodes[i] = broadcast.CreateNode(i, initPort+uint16(i), fifo.CreateFifoBroadcaster(nodesToPick, ttl, i))
		broadcasterNodes[i].OnBroadcastMessage(func(message *nodes.Message) {
			onMessage(i, message)
		})

		for j := uint32(0); j < nNodes; j++ {
			broadcasterNodes[i].AddOtherNodeConnection(j, j+1000)
		}
	}
	for i := uint32(0); i < nNodes; i++ {
		broadcasterNodes[i].StartListening()
	}
	for i := uint32(0); i < nNodes; i++ {
		broadcasterNodes[i].OpenConnectionsToNodes(broadcasterNodes)
	}

	broadcasterNodes[0].BroadcastString("Running on fifo :D")
	time.Sleep(time.Second * 5)

}

func onMessage(receivedNodeId uint32, message *nodes.Message) {
	fmt.Printf("[%d] RECEIVED BROADCASTE MESSAGE\n", receivedNodeId)
}
