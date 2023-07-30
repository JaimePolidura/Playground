package paxos

import (
	"distributed-systems/src/broadcast"
	"distributed-systems/src/broadcast/fifo"
	"distributed-systems/src/nodes"
	"distributed-systems/src/nodes/types"
	"sync/atomic"
	"time"
)

//The value of the paxos round is set in the Message property SeqNum

type PaxosNode struct {
	broadcast.Node

	proposalCounter uint64

	//Proposer
	numberAcceptorsPromised      uint32
	numberAcceptorsAccepted      uint32
	timeoutRequest               *time.Timer
	timeoutRequestMs             uint64
	valueBeingProposed           uint32
	proposalIdValueBeingProposed uint64
	promiseQuorumSatisfied       bool
	acceptQuorumSatisfied        bool

	//Acceptor
	proposalIdPromised   uint64
	proposalIdAccepted   uint64
	someProposalAccepted bool
	lastValueAccepted    uint32

	onConsensusReachedCallback func(value uint32)
}

func CreatePaxosNode(nodeId uint32, port uint16, timeoutRequestMs uint64, onConsensusReachedCallback func(value uint32)) *PaxosNode {
	paxosNode := &PaxosNode{
		Node:                       *broadcast.CreateNode(nodeId, port, fifo.CreateFifoBroadcaster(3, 6, nodeId)),
		onConsensusReachedCallback: onConsensusReachedCallback,
		timeoutRequestMs:           timeoutRequestMs,
	}

	paxosNode.AddMessageHandler(types.MESSAGE_PAXOS_PROMISE_ACCEPT, paxosNode.handlePromiseAcceptMessage)
	paxosNode.AddMessageHandler(types.MESSAGE_PAXOS_ACCEPTED, paxosNode.handleAcceptedMessage)
	paxosNode.AddMessageHandler(types.MESSAGE_PAXOS_PREPARE, paxosNode.handlePrepareMessage)
	paxosNode.AddMessageHandler(types.MESSAGE_PAXOS_PROMISE, paxosNode.handlePromiseMessage)
	paxosNode.AddMessageHandler(types.MESSAGE_PAXOS_ACCEPT, paxosNode.handleAcceptMessage)

	return paxosNode
}

func (this *PaxosNode) Propose(value uint32) {
	this.prepare(value)
}

func (this *PaxosNode) getNextProposalId() uint64 {
	newCounter := atomic.AddUint64(&this.proposalCounter, 1)
	return (newCounter << 32) | uint64(this.GetNodeId())
}

func (this *PaxosNode) sendToLearners(message *nodes.Message) {
}

func (this *PaxosNode) sendToAcceptors(message *nodes.Message) {
	this.GetConnectionManager().SendAll(message)
}
