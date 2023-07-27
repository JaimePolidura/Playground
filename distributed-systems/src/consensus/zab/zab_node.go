package zab

import (
	"distributed-systems/src/broadcast"
	"distributed-systems/src/broadcast/zab"
)

import (
	"time"
)

type ZabNode struct {
	node *broadcast.Node

	state NodeState

	heartbeatTimeMs       uint64
	heartbeatTimeoutMs    uint64
	heartbeatSenderTicker *time.Ticker
	heartbeatTimerTimeout *time.Timer

	nodesIdRing                       []uint32
	selfNodeIdRingIndex               uint32
	nNodesThatHaveAckElectionProposal uint32

	epoch uint32

	leaderNodeId uint32
}

func CreateZabNode(selfNodeId uint32, port uint16, leaderNodeId uint32, heartbeatTimeMs uint64, heartbeatTimeoutMs uint64, broadcasterNode *zab.ZabBroadcaster) *ZabNode {
	zabNode := &ZabNode{
		node:                  broadcast.CreateNode(selfNodeId, port, broadcasterNode),
		heartbeatSenderTicker: time.NewTicker(time.Duration(heartbeatTimeMs * uint64(time.Millisecond))),
		heartbeatTimerTimeout: time.NewTimer(time.Duration(heartbeatTimeoutMs * uint64(time.Millisecond))),
		heartbeatTimeMs:       heartbeatTimeMs,
		leaderNodeId:          leaderNodeId,
		epoch:                 0,
		state:                 STARTING,
		nodesIdRing:           make([]uint32, 0),
	}

	if zabNode.IsLeader() {
		go zabNode.startSendingHeartbeats()
	}
	if zabNode.IsFollower() {
		go zabNode.startHeartbeatTimer()
	}

	zabNode.node.AddMessageHandler(zab.MESSAGE_ACK_SUBMIT_RETRANSMISSION, broadcasterNode.HandleAckSubmitRetransmissionMessage)
	zabNode.node.AddMessageHandler(zab.MESSAGE_DO_BROADCAST, broadcasterNode.HandleDoBroadcast)
	zabNode.node.AddMessageHandler(zab.MESSAGE_ACK, broadcasterNode.HandleAckMessage)

	zabNode.node.AddMessageHandler(zab.MESSAGE_ELECTION_FAILURE_DETECTED, zabNode.handleNodeFailureMessage)
	zabNode.node.AddMessageHandler(zab.MESSAGE_HEARTBEAT, zabNode.handleHeartbeatMessage)
	zabNode.node.AddMessageHandler(zab.MESSAGE_ELECTION_COMMIT, zabNode.handleElectionCommitMessage)
	zabNode.node.AddMessageHandler(zab.MESSAGE_ELECTION_ACK_PROPOSAL, zabNode.handleElectionAckProposalMessage)
	zabNode.node.AddMessageHandler(zab.MESSAGE_ELECTION_PROPOSAL, zabNode.handleElectionProposalMessage)

	return zabNode
}

func (this *ZabNode) GetDurationHeartbeatTimeout(heartbeatSenderTicker uint64) time.Duration {
	return time.Duration(heartbeatSenderTicker*2 + heartbeatSenderTicker)
}

func (this *ZabNode) GetNodeId() uint32 {
	return this.node.GetNodeId()
}

func (this *ZabNode) GetNode() *broadcast.Node {
	return this.node
}

func (this *ZabNode) IsLeader() bool {
	return this.leaderNodeId == this.node.GetNodeId()
}

func (this *ZabNode) IsFollower() bool {
	return this.leaderNodeId != this.node.GetNodeId()
}

func (this *ZabNode) SetStateToBroadcast() {
	this.state = BROADCAST
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
