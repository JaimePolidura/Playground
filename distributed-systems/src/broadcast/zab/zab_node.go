package zab

import (
	"distributed-systems/src/broadcast"
	"distributed-systems/src/nodes"
	"distributed-systems/src/nodes/types"
	"sync"
)

import (
	"time"
)

type ZabNode struct {
	broadcast.Node

	state NodeState

	prevNodeRing uint32
	leaderNodeId uint32
	nodesIdRing  []uint32

	//Leader
	heartbeatSenderTicker *time.Ticker

	//Follower
	heartbeatTimerTimeout *time.Timer
	heartbeatTimeoutMs    uint64

	//Follower Election
	nNodesThatHaveAckElectionProposal uint32
	largestSeqNumSeenFromFollowers    uint32
	heartbeatCandidateTimerTimeout    *time.Timer
	nHeartbeatCandidateTimersTimeout  uint32
	nodesVotesRegistry                map[uint32]map[uint32]uint32 //Node id to vote -> nodeId voted
	lastFailedNodeId                  uint32
	lastFailedNodeLock                sync.Mutex
}

func CreateZabNode(selfNodeId uint32, port uint16, leaderNodeId uint32, heartbeatTimeMs uint64, heartbeatTimeoutMs uint64, prevNodeRing uint32, nodesIdRing []uint32, broadcasterNode *ZabBroadcaster, onBroadcastMessage func(message *nodes.Message)) *ZabNode {
	zabNode := &ZabNode{
		Node:                             *broadcast.CreateNode(selfNodeId, port, broadcasterNode),
		heartbeatSenderTicker:            time.NewTicker(time.Duration(heartbeatTimeMs * uint64(time.Millisecond))),
		heartbeatTimerTimeout:            time.NewTimer(time.Duration(heartbeatTimeoutMs * uint64(time.Millisecond))),
		nodesVotesRegistry:               make(map[uint32]map[uint32]uint32),
		heartbeatTimeoutMs:               heartbeatTimeoutMs,
		leaderNodeId:                     leaderNodeId,
		nodesIdRing:                      nodesIdRing,
		prevNodeRing:                     prevNodeRing,
		nHeartbeatCandidateTimersTimeout: 2,
		lastFailedNodeId:                 0xFFFFFFFF,
		state:                            STARTING,
	}

	if zabNode.IsLeader() {
		go zabNode.startSendingHeartbeats()
	}
	if zabNode.IsFollower() {
		go zabNode.startLeaderHeartbeatTimerTimeout()
	}

	zabNode.AddMessageHandler(types.MESSAGE_BROADCAST, onBroadcastMessage)

	zabNode.AddMessageHandler(types.MESSAGE_DO_BROADCAST, broadcasterNode.HandleDoBroadcast)
	zabNode.AddMessageHandler(types.MESSAGE_ACK, broadcasterNode.HandleAckMessage)

	zabNode.AddMessageHandler(types.MESSAGE_ZAB_ELECTION_FAILURE_DETECTED, zabNode.handleNodeFailureMessage)
	zabNode.AddMessageHandler(types.MESSAGE_HEARTBEAT, zabNode.handleHeartbeatMessage)
	zabNode.AddMessageHandler(types.MESSAGE_ZAB_ELECTION_COMMIT, zabNode.handleElectionCommitMessage)
	zabNode.AddMessageHandler(types.MESSAGE_ZAB_ELECTION_ACK_PROPOSAL, zabNode.handleElectionAckProposalMessage)
	zabNode.AddMessageHandler(types.MESSAGE_ZAB_ELECTION_PROPOSAL, zabNode.handleElectionProposalMessage)

	return zabNode
}

func (this *ZabNode) BroadcastString(toBroadcast string) {
	this.Broadcast(nodes.CreateMessage(
		nodes.WithNodeId(this.GetNodeId()),
		nodes.WithContentString(toBroadcast),
		nodes.WithType(types.MESSAGE_DO_BROADCAST)))
}

func (this *ZabNode) IsLeader() bool {
	return this.leaderNodeId == this.GetNodeId()
}

func (this *ZabNode) IsFollower() bool {
	return this.leaderNodeId != this.GetNodeId()
}

func (this *ZabNode) Stop() {
	this.state = STOPPED
	if this.heartbeatCandidateTimerTimeout != nil {
		this.heartbeatCandidateTimerTimeout.Stop()
	}
	if this.heartbeatTimerTimeout != nil {
		this.heartbeatTimerTimeout.Stop()
	}
	this.Node.Stop()
}

func (this *ZabNode) SetOnBroadcastMessageCallback(callback func(message *nodes.Message)) {
	//TODO
}

func (this *ZabNode) SetStateToBroadcast() {
	this.state = BROADCAST
}

func (this *ZabNode) getNextIndexInRingByIndex(prevIndex uint32) uint32 { //SenderNodeId <- Self
	if prevIndex+1 >= uint32(len(this.nodesIdRing)) {
		return 0
	} else {
		return prevIndex + 1
	}
}

func (this *ZabNode) getBackIndexInRingByIndex(nextIndex uint32) uint32 { //SenderNodeId <- Self
	if nextIndex-1 < 0 {
		return this.nodesIdRing[len(this.nodesIdRing)-1]
	} else {
		return nextIndex - 1
	}
}

func (this *ZabNode) getRingIndexByNodeId(nodeIdToSearch uint32) uint32 { //SenderNodeId <- Self
	for actualIndex, actualNodeId := range this.nodesIdRing {
		if actualNodeId == nodeIdToSearch {
			return uint32(actualIndex)
		}
	}

	panic("wtf")
}

func (this *ZabNode) getRingClockwiseDistanceByIndex(a uint32, b uint32) uint32 { //a -> b
	if b >= a {
		return b - a
	} else {
		return (uint32(len(this.nodesIdRing)) - (a + 1)) + b
	}
}
