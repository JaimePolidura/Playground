package zab

import (
	"distributed-systems/src/broadcast"
	"distributed-systems/src/broadcast/zab"
	"distributed-systems/src/nodes"
	"sync"
)

import (
	"time"
)

type ZabNode struct {
	node *broadcast.Node

	state NodeState

	prevNodeRing uint32
	leaderNodeId uint32

	//Follower
	heartbeatTimerTimeout *time.Timer
	heartbeatTimeoutMs    uint64

	//Follower Election
	nNodesThatHaveAckElectionProposal uint32
	largestSeqNumSeenFromFollowers    uint32
	electionLock                      sync.Mutex

	//Leader
	heartbeatSenderTicker *time.Ticker

	epoch uint32
}

func CreateZabNode(selfNodeId uint32, port uint16, leaderNodeId uint32, heartbeatTimeMs uint64, heartbeatTimeoutMs uint64, prevNodeRing uint32, broadcasterNode *zab.ZabBroadcaster) *ZabNode {
	zabNode := &ZabNode{
		node:                  broadcast.CreateNode(selfNodeId, port, broadcasterNode),
		heartbeatSenderTicker: time.NewTicker(time.Duration(heartbeatTimeMs * uint64(time.Millisecond))),
		heartbeatTimerTimeout: time.NewTimer(time.Duration(heartbeatTimeoutMs * uint64(time.Millisecond))),
		heartbeatTimeoutMs:    heartbeatTimeoutMs,
		leaderNodeId:          leaderNodeId,
		prevNodeRing:          prevNodeRing,
		epoch:                 0,
		state:                 STARTING,
	}

	if zabNode.IsLeader() {
		go zabNode.startSendingHeartbeats()
	}
	if zabNode.IsFollower() {
		go zabNode.startHeartbeatTimer()
	}

	zabNode.node.AddMessageHandler(zab.MESSAGE_DO_BROADCAST, broadcasterNode.HandleDoBroadcast)
	zabNode.node.AddMessageHandler(zab.MESSAGE_ACK, broadcasterNode.HandleAckMessage)

	zabNode.node.AddMessageHandler(zab.MESSAGE_ELECTION_FAILURE_DETECTED, zabNode.handleNodeFailureMessage)
	zabNode.node.AddMessageHandler(zab.MESSAGE_HEARTBEAT, zabNode.handleHeartbeatMessage)
	zabNode.node.AddMessageHandler(zab.MESSAGE_ELECTION_COMMIT, zabNode.handleElectionCommitMessage)
	zabNode.node.AddMessageHandler(zab.MESSAGE_ELECTION_ACK_PROPOSAL, zabNode.handleElectionAckProposalMessage)
	zabNode.node.AddMessageHandler(zab.MESSAGE_ELECTION_PROPOSAL, zabNode.handleElectionProposalMessage)

	return zabNode
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

func (this *ZabNode) Stop() {
	this.node.Stop()
	this.node.GetBroadcaster().(*zab.ZabBroadcaster).Stop()
}

func (this *ZabNode) GetConnectionManager() *nodes.ConnectionManager {
	return this.node.GetConnectionManager()
}

func (this *ZabNode) SetStateToBroadcast() {
	this.state = BROADCAST
}
