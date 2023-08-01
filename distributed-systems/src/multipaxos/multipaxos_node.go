package multipaxos

import (
	"distributed-systems/src/nodes"
	"distributed-systems/src/nodes/types"
	"distributed-systems/src/paxos"
	"fmt"
)

type MultiPaxosNode struct {
	Paxos *paxos.PaxosNode

	leaderNodeId uint32
	isLeader     bool
	leaderChosen bool

	//Leader
	valuesPendingAppendLog  chan uint32
	leaderLastAppendedIndex int32
	acceptedByIndexLog      map[int32]int32

	log                        []uint32
	lastAppendedCommittedIndex int32
	lastAppendedIndex          int32
}

func CreateMultiPaxosNode(nodeId uint32, port uint16, timeoutRequestMs uint64) *MultiPaxosNode {
	multipaxos := &MultiPaxosNode{
		Paxos:                  paxos.CreatePaxosNode(nodeId, port, timeoutRequestMs, func(value uint32) {}),
		valuesPendingAppendLog: make(chan uint32, 100),
		acceptedByIndexLog:     map[int32]int32{},
		log:                    make([]uint32, 10),

		leaderLastAppendedIndex:    -1,
		lastAppendedCommittedIndex: -1,
		lastAppendedIndex:          -1,
	}

	multipaxos.Paxos.SetOnConsensusReachedCallback(multipaxos.onLeaderConsensusReached)

	multipaxos.Paxos.AddMessageHandler(types.MESSAGE_MULTIPAXOS_REDIRECT_LEADER, multipaxos.handleRedirectLeaderMessage)
	multipaxos.Paxos.AddMessageHandler(types.MESSAGE_MULTIPAXOS_REPLAY_SUBMISSION, multipaxos.handleSubmissionReplay)
	multipaxos.Paxos.AddMessageHandler(types.MESSAGE_MULTIPAXOS_ACCEPTED, multipaxos.handleMultipaxosAcceptedMessage)
	multipaxos.Paxos.AddMessageHandler(types.MESSAGE_MULTIPAXOS_ACCEPT, multipaxos.handleMultipaxosAcceptMessage)
	multipaxos.Paxos.AddMessageHandler(types.MESSAGE_MULTIPAXOS_REPLAY_SUBMIT, multipaxos.handleReplaySubmit)

	return multipaxos
}

func (this *MultiPaxosNode) AppendLog(value uint32) {
	if !this.isLeader {
		fmt.Printf("[%d] Sending to leader append value: %d\n", this.Paxos.GetNodeId(), value)

		this.Paxos.GetConnectionManager().Send(this.leaderNodeId, nodes.CreateMessage(
			nodes.WithNodeId(this.Paxos.GetNodeId()),
			nodes.WithType(types.MESSAGE_MULTIPAXOS_REDIRECT_LEADER),
			nodes.WithContentUInt32(value)))
	}
	if this.isLeader {
		this.valuesPendingAppendLog <- value
	}
}

func (this *MultiPaxosNode) startListeningPendingAppendLog() {
	select {
	case valueToAppend := <-this.valuesPendingAppendLog:
		this.leaderLastAppendedIndex = this.leaderLastAppendedIndex + 1

		fmt.Printf("[%d] Received append value %d in index %d Sending MESSAGE_MULTIPAXOS_ACCEPT to followers\n",
			this.Paxos.GetNodeId(), valueToAppend, this.leaderLastAppendedIndex)

		this.sendToAcceptors(nodes.CreateMessage(
			nodes.WithNodeId(this.Paxos.GetNodeId()),
			nodes.WithContentInt32(this.leaderLastAppendedIndex),
			nodes.WithSeqNum(valueToAppend),
			nodes.WithType(types.MESSAGE_MULTIPAXOS_ACCEPT)))
	}
}

func (this *MultiPaxosNode) handleMultipaxosAcceptMessage(message *nodes.Message) {
	indexToAppend := message.GetContentToInt32()
	valueToAppend := message.SeqNum

	this.lastAppendedIndex = indexToAppend
	this.log[indexToAppend] = valueToAppend

	fmt.Printf("[%d] Received MESSAGE_MULTIPAXOS_ACCEPT new append log entry in index %d with value %d. The log %v. Last index commited %d\n",
		this.Paxos.GetNodeId(), indexToAppend, valueToAppend, this.log, this.lastAppendedCommittedIndex)

	if this.lastAppendedCommittedIndex+1 == indexToAppend {
		this.lastAppendedCommittedIndex = indexToAppend
	}
	if this.lastAppendedCommittedIndex+1 < indexToAppend {
		this.Paxos.GetConnectionManager().Send(this.leaderNodeId, nodes.CreateMessage(
			nodes.WithNodeId(this.Paxos.GetNodeId()),
			nodes.WithType(types.MESSAGE_MULTIPAXOS_REPLAY_SUBMIT),
			nodes.WithContentsInt32(this.lastAppendedCommittedIndex)))
	}

	this.Paxos.GetConnectionManager().Send(this.leaderNodeId, nodes.CreateMessage(
		nodes.WithNodeId(this.Paxos.GetNodeId()),
		nodes.WithType(types.MESSAGE_MULTIPAXOS_ACCEPTED),
		nodes.WithContentsInt32(indexToAppend, int32(valueToAppend))))
}

func (this *MultiPaxosNode) handleMultipaxosAcceptedMessage(message *nodes.Message) {
	logIndex := message.GetContentToInt32()
	valueToAppend := message.GetContentToUint32WithOffset(4)
	this.acceptedByIndexLog[logIndex] = this.acceptedByIndexLog[logIndex] + 1

	nodesAcceptedByIndex := this.acceptedByIndexLog[logIndex] + 1
	isQuorumSatisfied := nodesAcceptedByIndex >= int32(this.Paxos.GetConnectionManager().GetNumberConnections()/2+1)

	fmt.Printf("[%d] Received accepted message from node %d of index %d Nodes accepted in index %d Is quorum satisfied? %t\n",
		this.Paxos.GetNodeId(), message.NodeIdSender, logIndex, nodesAcceptedByIndex, isQuorumSatisfied)

	if isQuorumSatisfied {
		fmt.Printf("[%d] Setting last commited index %d The log %v\n",
			this.Paxos.GetNodeId(), logIndex, this.log)

		this.lastAppendedCommittedIndex = logIndex
		this.log[logIndex] = valueToAppend

		go this.startListeningPendingAppendLog()
	}
}

func (this *MultiPaxosNode) handleReplaySubmit(message *nodes.Message) {
	lastFollowerCommittedIndex := message.GetContentToUint32()
	followerNodeId := message.NodeIdSender

	for i := lastFollowerCommittedIndex + 1; i < uint32(len(this.log)); i++ {
		valueLog := this.log[i]

		this.Paxos.GetConnectionManager().Send(followerNodeId, nodes.CreateMessage(
			nodes.WithNodeId(this.Paxos.GetNodeId()),
			nodes.WithType(types.MESSAGE_MULTIPAXOS_REPLAY_SUBMISSION),
			nodes.WithContentsUInt32(i, valueLog)))
	}
}

func (this *MultiPaxosNode) handleSubmissionReplay(message *nodes.Message) {
	indexReplay := message.GetContentToInt32WithOffset(0)
	value := message.GetContentToUint32WithOffset(4)

	this.lastAppendedCommittedIndex = indexReplay
	this.log[indexReplay] = value
}

func (this *MultiPaxosNode) handleRedirectLeaderMessage(message *nodes.Message) {
	value := message.GetContentToUint32()
	this.valuesPendingAppendLog <- value
}

func (this *MultiPaxosNode) SetLeader() {
	fmt.Printf("[%d] Setting me as a leader\n", this.Paxos.GetNodeId())
	this.Paxos.Prepare(this.Paxos.GetNodeId())
}

func (this *MultiPaxosNode) onLeaderConsensusReached(leaderChosenNodeId uint32) {
	if this.leaderChosen {
		return
	}

	fmt.Printf("[%d] Leader with id %d has been chosen", this.Paxos.GetNodeId(), leaderChosenNodeId)
	fmt.Println("  ")

	this.leaderNodeId = leaderChosenNodeId
	this.leaderChosen = true

	if this.leaderNodeId == this.Paxos.GetNodeId() {
		this.isLeader = true
	}

	go this.startListeningPendingAppendLog()
}

func (this *MultiPaxosNode) sendToAcceptors(message *nodes.Message) {
	this.Paxos.GetConnectionManager().SendAllExcept(this.Paxos.GetNodeId(), message)
}
