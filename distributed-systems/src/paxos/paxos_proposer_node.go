package paxos

import (
	"distributed-systems/src/nodes"
	"distributed-systems/src/nodes/types"
	"fmt"
	"time"
)

func (this *PaxosNode) setupTimerRequestTimeout() {
	this.stopTimerRequestTimeout()
	this.startTimerRequestTimeout()

	go func() {
		select {
		case <-this.timeoutRequest.C:
			this.prepare(this.valueBeingProposed) //Retry with a higher ID
		}
	}()
}

func (this *PaxosNode) prepare(value uint32) {
	nextProposalId := this.getNextProposalId()

	fmt.Printf("[%d] Sending PREPARE(%d) to acceptors with value %d\n",
		this.GetNodeId(), nextProposalId, value)

	this.setupTimerRequestTimeout()
	this.sendToAcceptors(nodes.CreateMessage(
		nodes.WithNodeId(this.GetNodeId()),
		nodes.WithContentUInt64(nextProposalId),
		nodes.WithType(types.MESSAGE_PAXOS_PREPARE),
		nodes.WithSeqNum(value)))
}

func (this *PaxosNode) handlePromiseMessage(message *nodes.Message) {
	this.stopTimerRequestTimeout()

	proposalIdPromise := message.GetContentToUint64()
	this.doHandlePromiseMessage(proposalIdPromise, message.NodeIdSender, message.SeqNum)
}

func (this *PaxosNode) handlePromiseAcceptMessage(message *nodes.Message) {
	this.stopTimerRequestTimeout()

	promiseProposalId := message.GetContentToUint64()
	acceptedProposalId := message.GetContentToUint64WithOffset(8)
	acceptedValue := message.SeqNum

	this.valueBeingProposed = acceptedValue

	if promiseProposalId < acceptedProposalId {
		this.proposalIdValueBeingProposed = acceptedProposalId
	}

	fmt.Printf("[%d] Received PROMISE_ACCEPT(%d, %d) from acceptor node %d of already accepted value %d\n",
		this.GetNodeId(), promiseProposalId, acceptedProposalId, message.NodeIdSender, acceptedValue)

	this.doHandlePromiseMessage(acceptedProposalId, message.NodeIdSender, message.SeqNum)
}

func (this *PaxosNode) doHandlePromiseMessage(proposalId uint64, nodeIdSender uint32, value uint32) {
	this.stopTimerRequestTimeout()

	if this.promiseQuorumSatisfied {
		return
	}

	this.numberAcceptorsPromised++
	minNodesQuorum := this.Node.GetConnectionManager().GetNumberConnections()/2 + 1
	quorumSatisfied := this.numberAcceptorsPromised >= minNodesQuorum

	fmt.Printf("[%d] Received PROMISE(%d) from acceptor node %d Nº Acceptor nodes promised: %d Is quorum satisfied? %t\n",
		this.GetNodeId(), proposalId, nodeIdSender, this.numberAcceptorsPromised, quorumSatisfied)

	if quorumSatisfied {
		this.promiseQuorumSatisfied = true

		fmt.Printf("[%d] Promise quorum satisfied. Sending to proposers ACCEPT(%d)\n",
			this.GetNodeId(), proposalId)

		this.setupTimerRequestTimeout()
		this.sendToAcceptors(nodes.CreateMessage(
			nodes.WithNodeId(this.GetNodeId()),
			nodes.WithContentUInt64(proposalId),
			nodes.WithSeqNum(value),
			nodes.WithType(types.MESSAGE_PAXOS_ACCEPT)))
	}
}

func (this *PaxosNode) handleAcceptedMessage(message *nodes.Message) {
	this.stopTimerRequestTimeout()

	if this.acceptQuorumSatisfied {
		return
	}

	this.numberAcceptorsAccepted++
	minNodesQuorum := this.Node.GetConnectionManager().GetNumberConnections()/2 + 1
	quorumSatisfied := this.numberAcceptorsAccepted >= minNodesQuorum
	proposalId := message.GetContentToUint64()

	fmt.Printf("[%d] Received ACCEPTED(%d) from acceptor node %d Nº Acceptor nodes accepted: %d Is quorum satisfied? %t\n",
		this.GetNodeId(), proposalId, message.NodeIdSender, this.numberAcceptorsAccepted, quorumSatisfied)

	if quorumSatisfied {
		this.acceptQuorumSatisfied = true
		this.onConsensusReachedCallback(message.SeqNum)
	}
}

func (this *PaxosNode) startTimerRequestTimeout() {
	timeout := time.Duration(this.timeoutRequestMs * uint64(time.Millisecond))

	if this.timeoutRequest != nil {
		this.timeoutRequest.Reset(timeout)
	} else {
		this.timeoutRequest = time.NewTimer(timeout)
	}
}

func (this *PaxosNode) stopTimerRequestTimeout() {
	if this.timeoutRequest != nil {
		this.timeoutRequest.Stop()
	}
}
