package main

import (
	"distributed-systems/src/broadcast"
	"distributed-systems/src/broadcast/fifo"
	"distributed-systems/src/broadcast/zab"
	zab2 "distributed-systems/src/consensus/zab"
	"distributed-systems/src/nodes"
	"fmt"
	"time"
)

func main() {
	startZab()
	//startFifo()
}

func startZab() {
	nNodes := uint32(4)
	initPort := uint16(1000)
	zabNodes := make([]*zab2.ZabNode, nNodes)

	for nodeId := uint32(0); nodeId < nNodes; nodeId++ {
		copyOfNodeId := nodeId

		zabNodes[nodeId] = zab2.CreateZabNode(nodeId,
			initPort+uint16(nodeId),
			0,
			100,
			2000,
			zab.CreateZabBroadcaster(nodeId, 0, 1500, func(newMessage *nodes.Message) { onMessage(copyOfNodeId, newMessage) }))

		for otherNodeId := uint32(0); otherNodeId < nNodes; otherNodeId++ {
			zabNodes[nodeId].GetNode().AddOtherNodeConnection(otherNodeId, otherNodeId+1000)
		}
	}
	for nodeId := uint32(0); nodeId < nNodes; nodeId++ {
		zabNodes[nodeId].GetNode().StartListeningAsync()
	}
	for nodeId := uint32(0); nodeId < nNodes; nodeId++ {
		zabNodes[nodeId].GetNode().GetConnectionManager().OpenAllConnections()
	}
	for _, zabNode := range zabNodes {
		zabNode.SetStateToBroadcast()
	}

	zabNodes[1].GetNode().BroadcastString("Running on zab 1º!", zab.MESSAGE_DO_BROADCAST)
	time.Sleep(time.Second * 2)
	fmt.Println("    ")
	zabNodes[0].Stop()
	time.Sleep(time.Second * 500)
}

func startFifo() {
	nNodes := uint32(6)
	ttl := int32(3)
	nodesToPick := uint32(2)
	initPort := uint16(1000)

	broadcasterNodes := make([]*broadcast.Node, nNodes)

	for i := uint32(0); i < nNodes; i++ {
		nodeId := i

		broadcasterNodes[i] = broadcast.CreateNode(i, initPort+uint16(i), fifo.CreateFifoBroadcaster(nodesToPick, ttl, i))
		broadcasterNodes[i].OnBroadcastMessage(func(message *nodes.Message) {
			onMessage(nodeId, message)
		})

		for j := uint32(0); j < nNodes; j++ {
			broadcasterNodes[i].AddOtherNodeConnection(j, j+1000)
		}

		broadcasterNodes[i].StartListeningAsync()
	}
	for i := uint32(0); i < nNodes; i++ {
		broadcasterNodes[i].GetConnectionManager().OpenAllConnections()
	}

	broadcasterNodes[1].BroadcastString("Running on fifo :D", zab2.BROADCAST)
	time.Sleep(time.Second * 5000)
}

func onMessage(receivedNodeId uint32, message *nodes.Message) {
	fmt.Printf("[%d] %s\n", receivedNodeId, message.Content)
}
