package main

import (
	"bufio"
	"distributed-systems/src/broadcast"
	"distributed-systems/src/broadcast/fifo"
	"distributed-systems/src/broadcast/zab"
	"distributed-systems/src/multipaxos"
	"distributed-systems/src/nodes"
	"distributed-systems/src/nodes/types"
	"distributed-systems/src/paxos"
	"distributed-systems/src/raft"
	"fmt"
	"os"
	"time"
)

func main() {
	//startFifo()
	//startZab()
	//startPaxos()
	//startMultipaxos()
	startRaft()
}

func startRaft() {
	nNodes := uint32(6)
	raftNodes := make([]*raft.RaftNode, nNodes)

	for nodeId := uint32(0); nodeId < nNodes; nodeId++ {
		raftNodes[nodeId] = raft.CreateRaftNode(1000, 1000, 1000, 0, nodeId, uint16(nodeId)+1000)

		for otherNodeId := uint32(0); otherNodeId < nNodes; otherNodeId++ {
			raftNodes[nodeId].Node.AddOtherNodeConnection(otherNodeId, otherNodeId+1000)
		}
	}

	for nodeId := uint32(0); nodeId < nNodes; nodeId++ {
		raftNodes[nodeId].Node.StartListeningAsync()
	}
	for nodeId := uint32(0); nodeId < nNodes; nodeId++ {
		raftNodes[nodeId].Node.GetConnectionManager().OpenAllConnections()
	}
	for nodeId := uint32(0); nodeId < nNodes; nodeId++ {
		raftNodes[nodeId].Start()
	}

	blockMainThread()
}

func startMultipaxos() {
	nNodes := uint32(6)
	multiPaxosNodes := make([]*multipaxos.MultiPaxosNode, nNodes)

	for nodeId := uint32(0); nodeId < nNodes; nodeId++ {
		multiPaxosNodes[nodeId] = multipaxos.CreateMultiPaxosNode(nodeId, uint16(nodeId)+1000, 2000)

		for otherNodeId := uint32(0); otherNodeId < nNodes; otherNodeId++ {
			multiPaxosNodes[nodeId].Paxos.AddOtherNodeConnection(otherNodeId, otherNodeId+1000)
		}
	}
	for nodeId := uint32(0); nodeId < nNodes; nodeId++ {
		multiPaxosNodes[nodeId].Paxos.StartListeningAsync()
	}
	for nodeId := uint32(0); nodeId < nNodes; nodeId++ {
		multiPaxosNodes[nodeId].Paxos.GetConnectionManager().OpenAllConnections()
	}

	multiPaxosNodes[0].SetLeader()
	time.Sleep(1 * time.Second)
	multiPaxosNodes[1].AppendLog(1)
	multiPaxosNodes[2].AppendLog(13)

	blockMainThread()
}

func startPaxos() {
	nNodes := uint32(6)
	paxosNodes := make([]*paxos.PaxosNode, nNodes)

	for nodeId := uint32(0); nodeId < nNodes; nodeId++ {
		paxosNodes[nodeId] = paxos.CreatePaxosNode(nodeId, uint16(nodeId)+1000, 2000, onPaxosConsensus)

		for otherNodeId := uint32(0); otherNodeId < nNodes; otherNodeId++ {
			paxosNodes[nodeId].AddOtherNodeConnection(otherNodeId, otherNodeId+1000)
		}
	}
	for nodeId := uint32(0); nodeId < nNodes; nodeId++ {
		paxosNodes[nodeId].StartListeningAsync()
	}
	for nodeId := uint32(0); nodeId < nNodes; nodeId++ {
		paxosNodes[nodeId].GetConnectionManager().OpenAllConnections()
	}

	paxosNodes[0].Prepare(11)
	paxosNodes[1].Prepare(12)
	paxosNodes[2].Prepare(13)
	paxosNodes[3].Prepare(14)

	blockMainThread()
}

func onPaxosConsensus(value uint32) {
	fmt.Println("REACHED CONSENSUS ON VALUE ", value)
}

func startZab() {
	nNodes := uint32(4)
	initPort := uint16(1000)
	zabNodes := make([]*zab.ZabNode, nNodes)

	for nodeId := uint32(0); nodeId < nNodes; nodeId++ {
		copyOfNodeId := nodeId
		prevNodeId := nodeId - 1

		if prevNodeId < 0 {
			prevNodeId = nNodes - 1
		}

		zabNodes[nodeId] = zab.CreateZabNode(nodeId,
			initPort+uint16(nodeId),
			0,
			250,
			1000,
			prevNodeId,
			[]uint32{0, 1, 2, 3},
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

	zabNodes[1].GetNode().BroadcastString("Running on zab 1º!", types.MESSAGE_DO_BROADCAST)
	time.Sleep(time.Second * 2)
	fmt.Println("    ")
	zabNodes[0].Stop()
	zabNodes[1].Stop()
	time.Sleep(time.Second * 5)
	fmt.Println("    ")
	zabNodes[2].GetNode().BroadcastString("Joder!", types.MESSAGE_DO_BROADCAST)
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

	broadcasterNodes[1].BroadcastString("Running on fifo :D", zab.BROADCAST)
	time.Sleep(time.Second * 5000)
}

func onMessage(receivedNodeId uint32, message *nodes.Message) {
	fmt.Printf("[%d] %s\n", receivedNodeId, message.Content)
}

func blockMainThread() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		scanner.Text()
	}
}
