package paxos

import (
	"distributed-systems/src/nodes"
	"distributed-systems/src/nodes/types"
)

func (this *PaxosNode) handlePrepareMessage(message *nodes.Message) {
	proposalId := message.GetContentToUint64()
	proposerNodeId := message.NodeIdSender
	value := message.SeqNum

	if proposalId < this.proposalIdPromised {
		return
	}

	if !this.someProposalAccepted {
		this.proposalIdPromised = proposalId

		this.GetConnectionManager().Send(proposerNodeId, nodes.CreateMessage(
			nodes.WithNodeId(this.GetNodeId()),
			nodes.WithContentUInt64(proposalId),
			nodes.WithType(types.MESSAGE_PAXOS_PROMISE),
			nodes.WithSeqNum(value)))
	} else {
		this.GetConnectionManager().Send(proposerNodeId, nodes.CreateMessage(
			nodes.WithNodeId(this.GetNodeId()),
			nodes.WithContentsUInt64(proposalId, this.proposalIdAccepted),
			nodes.WithType(types.MESSAGE_PAXOS_PROMISE_ACCEPT),
			nodes.WithSeqNum(this.lastValueAccepted)))
	}
}

func (this *PaxosNode) handleAcceptMessage(message *nodes.Message) {
	proposalId := message.GetContentToUint64()
	proposerNodeId := message.NodeIdSender
	value := message.SeqNum

	if proposalId >= this.proposalIdPromised {
		this.proposalIdAccepted = proposalId
		this.someProposalAccepted = true
		this.lastValueAccepted = value

		this.GetConnectionManager().Send(proposerNodeId, nodes.CreateMessage(
			nodes.WithNodeId(this.GetNodeId()),
			nodes.WithContentUInt64(proposalId),
			nodes.WithType(types.MESSAGE_PAXOS_ACCEPTED),
			nodes.WithSeqNum(value)))
	}
}
