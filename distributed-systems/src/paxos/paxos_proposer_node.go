package paxos

import (
	"distributed-systems/src/nodes"
	"distributed-systems/src/nodes/types"
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
	this.doHandlePromiseMessage(proposalIdPromise)
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

	this.doHandlePromiseMessage(acceptedProposalId)
}

func (this *PaxosNode) doHandlePromiseMessage(proposalId uint64) {
	this.stopTimerRequestTimeout()

	this.numberAcceptorsPromised++
	minNodesQuorum := this.Node.GetConnectionManager().GetNumberConnections()/2 + 1
	quorumSatisfied := this.numberAcceptorsPromised >= minNodesQuorum

	if quorumSatisfied {
		this.setupTimerRequestTimeout()
		this.sendToAcceptors(nodes.CreateMessage(
			nodes.WithNodeId(this.GetNodeId()),
			nodes.WithContentUInt64(proposalId),
			nodes.WithType(types.MESSAGE_PAXOS_ACCEPT)))
	}
}

func (this *PaxosNode) handleAcceptedMessage(message *nodes.Message) {
	this.stopTimerRequestTimeout()

	this.numberAcceptorsPromised++
	minNodesQuorum := this.Node.GetConnectionManager().GetNumberConnections()/2 + 1
	quorumSatisfied := this.numberAcceptorsPromised >= minNodesQuorum

	if quorumSatisfied {
		this.onConsensusReachedCallback(message.SeqNum)
	}
}

func (this *PaxosNode) startTimerRequestTimeout() {
	timeout := time.Duration(this.timeoutPrepareMs * uint64(time.Millisecond))

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
