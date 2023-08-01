package paxos

import (
	"distributed-systems/src/nodes"
	"distributed-systems/src/nodes/types"
	"fmt"
)

func (this *PaxosNode) handlePrepareMessage(message *nodes.Message) {
	proposalId := message.GetContentToUint64()
	proposerNodeId := message.NodeIdSender
	valuePrepared := message.SeqNum

	if proposalId < this.proposalIdPromised {
		fmt.Printf("[%d] Ignoring received PREPARE(%d) from proposer node %d of valuePrepared %d Promised to ignore ids lower than %d\n",
			this.GetNodeId(), proposalId, message.NodeIdSender, valuePrepared, this.proposalIdPromised)
		return
	}

	if !this.someProposalAccepted {
		fmt.Printf("[%d] Received PREPARE(%d) from proposer node %d of valuePrepared %d Sending back PROMISE(%d)\n",
			this.GetNodeId(), proposalId, message.NodeIdSender, valuePrepared, proposalId)

		this.proposalIdPromised = proposalId

		this.GetConnectionManager().Send(proposerNodeId, nodes.CreateMessage(
			nodes.WithNodeId(this.GetNodeId()),
			nodes.WithContentUInt64(proposalId),
			nodes.WithType(types.MESSAGE_PAXOS_PROMISE),
			nodes.WithSeqNum(valuePrepared)))
	} else {
		fmt.Printf("[%d] Received PREPARE(%d) from proposer node %d of valuePrepared %d But already accepted proposal id with %d and valuePrepared %d. Sending back PROMISE_ACCPET(%d, %d)\n",
			this.GetNodeId(), proposalId, message.NodeIdSender, this.proposalIdAccepted, this.lastValueAccepted, this.proposalIdAccepted, proposalId, this.proposalIdAccepted)

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
		fmt.Printf("[%d] Received ACCEPT(%d) from proposer node %d of value %d. Sending back ACCEPTED(%d)\n",
			this.GetNodeId(), proposalId, proposerNodeId, value, proposalId)

		this.proposalIdAccepted = proposalId
		this.someProposalAccepted = true
		this.lastValueAccepted = value

		this.GetConnectionManager().Send(proposerNodeId, nodes.CreateMessage(
			nodes.WithNodeId(this.GetNodeId()),
			nodes.WithContentUInt64(proposalId),
			nodes.WithType(types.MESSAGE_PAXOS_ACCEPTED),
			nodes.WithSeqNum(value)))
	} else {
		fmt.Printf("[%d] Ignoring received ACCEPT(%d) from node %d Promised to ignore ids lower than %d\n",
			this.GetNodeId(), proposalId, message.NodeIdSender, this.proposalIdPromised)
	}
}
