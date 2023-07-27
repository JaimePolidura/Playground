package zab

import (
	"distributed-systems/src/broadcast/zab"
	"distributed-systems/src/nodes"
	"time"
)

func (this *ZabNode) startHeartbeatTimer() {
	select {
	case <-this.heartbeatTimerTimeout.C:
		if this.IsFollower() && this.state == BROADCAST {
			this.node.Broadcast(nodes.CreateMessage(
				nodes.WithNodeId(this.GetNodeId()),
				nodes.WithType(zab.MESSAGE_ELECTION_FAILURE_DETECTED),
				nodes.WithFlags(nodes.FLAG_BYPASS_LEADER)))
		}
	}
}

func (this *ZabNode) handleNodeFailureMessage(message []*nodes.Message) {
	this.node.DisableBroadcast()
	this.state = ELECTION

	distance := this.getRingDistanceClockwise(this.leaderNodeId)

	if distance == 1 {
		this.node.Broadcast(nodes.CreateMessage(
			nodes.WithNodeId(this.GetNodeId()),
			nodes.WithType(zab.MESSAGE_ELECTION_PROPOSAL),
			nodes.WithFlags(nodes.FLAG_BYPASS_LEADER)))
	}
}

func (this *ZabNode) handleElectionProposalMessage(messages []*nodes.Message) {
	proposerNodeId := messages[0].NodeIdSender
	nodeConnection := this.node.GetNodeConnectionsStore().Get(proposerNodeId)

	nodeConnection.Write(nodes.CreateMessage(
		nodes.WithNodeId(this.GetNodeId()),
		nodes.WithType(zab.MESSAGE_ELECTION_PROPOSAL),
		nodes.WithFlags(nodes.FLAG_BYPASS_LEADER)))
}

func (this *ZabNode) handleElectionAckProposalMessage(message []*nodes.Message) {
	this.nNodesThatHaveAckElectionProposal++
	nNodesQuorum := this.node.GetNodeConnectionsStore().Size()/2 + 1
	isQuorumSatisfied := this.nNodesThatHaveAckElectionProposal >= nNodesQuorum

	if isQuorumSatisfied {
		this.node.Broadcast(nodes.CreateMessage(
			nodes.WithNodeId(this.GetNodeId()),
			nodes.WithType(zab.MESSAGE_ELECTION_COMMIT),
			nodes.WithFlags(nodes.FLAG_BYPASS_LEADER)))

		this.saveNewLeader(this.selfNodeIdRingIndex)
	}
}

func (this *ZabNode) handleHeartbeatMessage(message []*nodes.Message) {
	this.heartbeatTimerTimeout.Reset(time.Duration(this.heartbeatTimeoutMs))
}

func (this *ZabNode) handleElectionCommitMessage(messages []*nodes.Message) {
	this.saveNewLeader(messages[0].NodeIdSender)
}

func (this *ZabNode) saveNewLeader(newLeaderNodeId uint32) {
	this.leaderNodeId = newLeaderNodeId
	this.state = BROADCAST
	this.node.EnableBroadcast()
}
