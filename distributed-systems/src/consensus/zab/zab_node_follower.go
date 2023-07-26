package zab

import "distributed-systems/src/nodes"

func (this *ZabNode) startHeartbeatTimer() {
	select {
	case <-this.heartbeatTimerTimeout.C:
		this.node.Broadcast(nodes.CreateMessageBroadcast(this.node.GetNodeId(), this.node.GetNodeId(), "").WithType(MESSAGE_ELECTION_FAILURE_DETECTED))
	}
}

func (this *ZabNode) handleNodeFailureMessage() {
	this.node.DisableBroadcast()
	this.state = ELECTION

	distance := this.getRingDistanceClockwise(this.leaderNodeId)

	if distance == 1 {
		this.node.Broadcast(nodes.CreateMessageBroadcast(this.node.GetNodeId(), this.node.GetNodeId(), "").WithType(MESSAGE_ELECTION_PROPOSAL))
	}
}

func (this *ZabNode) ackProposal(message *nodes.Message) {
	proposerNodeId := message.NodeIdSender
	nodeConnection := this.node.GetNodeConnectionsStore().Get(proposerNodeId)

	nodeConnection.Write(nodes.CreateMessageWithType(this.GetNodeId(), this.GetNodeId(), "", MESSAGE_ELECTION_ACK_PROPOSAL))
}

func (this *ZabNode) collectAckProposal() {
	this.nNodesThatHaveAckElectionProposal++
	nNodesQuorum := this.node.GetNodeConnectionsStore().Size()/2 + 1
	isQuorumSatisfied := this.nNodesThatHaveAckElectionProposal >= nNodesQuorum

	if isQuorumSatisfied {
		this.node.Broadcast(nodes.CreateMessageBroadcast(this.node.GetNodeId(), this.node.GetNodeId(), "").WithType(MESSAGE_ELECTION_COMMIT))
		this.saveNewLeader(this.selfNodeIdRingIndex)
	}
}

func (this *ZabNode) saveNewLeader(newLeaderNodeId uint32) {
	this.leaderNodeId = newLeaderNodeId
	this.state = BROADCAST
	this.node.EnableBroadcast()
}

func (this *ZabNode) getRingDistanceClockwise(otherNodeId uint32) uint32 {
	indexOfOtherNode := uint32(0)

	for index, nodeId := range this.nodesIdRing {
		if nodeId == otherNodeId {
			indexOfOtherNode = uint32(index)
			break
		}
	}

	if indexOfOtherNode > this.selfNodeIdRingIndex {
		return indexOfOtherNode - this.selfNodeIdRingIndex
	} else {
		return this.selfNodeIdRingIndex + (uint32(len(this.nodesIdRing)) - indexOfOtherNode) + 1
	}
}
