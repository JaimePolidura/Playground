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
			1000,
			zab.CreateZabBroadcaster(nodeId, 0, 1500, func(newMessage *nodes.Message) { onMessage(copyOfNodeId, newMessage) }))

		for otherNodeId := uint32(0); otherNodeId < nNodes; otherNodeId++ {
			zabNodes[nodeId].GetNode().AddOtherNodeConnection(otherNodeId, otherNodeId+1000)
		}
	}
	for nodeId := uint32(0); nodeId < nNodes; nodeId++ {
		zabNodes[nodeId].GetNode().StartListening()
	}
	for nodeId := uint32(0); nodeId < nNodes; nodeId++ {
		for _, zabNode := range zabNodes {
			zabNodes[nodeId].GetNode().OpenConnectionToNode(zabNode.GetNode())
		}
	}
	for _, zabNode := range zabNodes {
		zabNode.SetStateToBroadcast()
	}

	zabNodes[1].GetNode().BroadcastString("Running on zab 1º!")
	//zabNodes[1].GetNode().BroadcastString("Running on zab 2º!")
	time.Sleep(time.Second * 5)
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
	}
	for i := uint32(0); i < nNodes; i++ {
		broadcasterNodes[i].StartListening()
	}
	for i := uint32(0); i < nNodes; i++ {
		broadcasterNodes[i].OpenConnectionsToNodes(broadcasterNodes)
	}

	broadcasterNodes[0].BroadcastString("Running on fifo :D")
	time.Sleep(time.Second * 5000)
}

func onMessage(receivedNodeId uint32, message *nodes.Message) {
	fmt.Printf("[%d] %s\n", receivedNodeId, message.Content)
}
