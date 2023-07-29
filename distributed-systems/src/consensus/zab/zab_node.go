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
	nodesIdRing  []uint32

	//Follower
	heartbeatTimerTimeout *time.Timer
	heartbeatTimeoutMs    uint64

	//Follower Election
	nNodesThatHaveAckElectionProposal uint32
	largestSeqNumSeenFromFollowers    uint32
	heartbeatCandidateTimerTimeout    *time.Timer
	nHeartbeatCandidateTimersTimeout  uint32
	nodesVotesRegistry                map[uint32]map[uint32]uint32 //Node id to vote -> nodeId voted

	lastFailedNodeId   uint32
	lastFailedNodeLock sync.Mutex

	//Leader
	heartbeatSenderTicker *time.Ticker

	epoch uint32
}

func CreateZabNode(selfNodeId uint32, port uint16, leaderNodeId uint32, heartbeatTimeMs uint64, heartbeatTimeoutMs uint64, prevNodeRing uint32, nodesIdRing []uint32, broadcasterNode *zab.ZabBroadcaster) *ZabNode {
	zabNode := &ZabNode{
		node:                             broadcast.CreateNode(selfNodeId, port, broadcasterNode),
		heartbeatSenderTicker:            time.NewTicker(time.Duration(heartbeatTimeMs * uint64(time.Millisecond))),
		heartbeatTimerTimeout:            time.NewTimer(time.Duration(heartbeatTimeoutMs * uint64(time.Millisecond))),
		nodesVotesRegistry:               make(map[uint32]map[uint32]uint32),
		heartbeatTimeoutMs:               heartbeatTimeoutMs,
		leaderNodeId:                     leaderNodeId,
		nodesIdRing:                      nodesIdRing,
		prevNodeRing:                     prevNodeRing,
		nHeartbeatCandidateTimersTimeout: 2,
		epoch:                            0,
		lastFailedNodeId:                 0xFFFFFFFF,
		state:                            STARTING,
	}

	if zabNode.IsLeader() {
		go zabNode.startSendingHeartbeats()
	}
	if zabNode.IsFollower() {
		go zabNode.startLeaderHeartbeatTimerTimeout()
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
	this.state = STOPPED
	if this.heartbeatCandidateTimerTimeout != nil {
		this.heartbeatCandidateTimerTimeout.Stop()
	}
	if this.heartbeatTimerTimeout != nil {
		this.heartbeatTimerTimeout.Stop()
	}
	this.node.Stop()
	this.node.GetBroadcaster().(*zab.ZabBroadcaster).Stop()
}

func (this *ZabNode) GetConnectionManager() *nodes.ConnectionManager {
	return this.node.GetConnectionManager()
}

func (this *ZabNode) SetStateToBroadcast() {
	this.state = BROADCAST
}

func (this *ZabNode) setupHeartbeatCandidateTimerTimeout(failedNode uint32) { //NodeId <- Self
	if this.heartbeatCandidateTimerTimeout != nil {
		this.heartbeatCandidateTimerTimeout.Stop()
	}

	indexFailedNode := this.getRingIndexByNodeId(failedNode)
	indexSelfNode := this.getRingIndexByNodeId(this.GetNodeId())
	indexCandidate := this.getNextIndexInRingByIndex(indexFailedNode)
	distanceFromFailedCandidate := this.getRingClockwiseDistanceByIndex(indexCandidate, indexSelfNode)

	if distanceFromFailedCandidate == 0 {
		return
	}
	if distanceFromFailedCandidate <= this.nHeartbeatCandidateTimersTimeout {
		timeout := time.Duration(uint64(distanceFromFailedCandidate) * (this.heartbeatTimeoutMs + 1) * uint64(time.Millisecond))
		this.heartbeatCandidateTimerTimeout = time.NewTimer(timeout)
		go this.startCandidateHeartbeatTimerTimeout()
	}
}

func (this *ZabNode) restartHeartbeatCandidateTimerTimeout() {
	indexFailedNode := this.getRingIndexByNodeId(this.lastFailedNodeId)
	indexSelfNode := this.getRingIndexByNodeId(this.GetNodeId())
	indexCandidate := this.getNextIndexInRingByIndex(indexFailedNode)

	distanceFromFailedCandidate := this.getRingClockwiseDistanceByIndex(indexCandidate, indexSelfNode)
	timeout := time.Duration(uint64(distanceFromFailedCandidate) * (this.heartbeatTimeoutMs + 1) * uint64(time.Millisecond))
	this.heartbeatCandidateTimerTimeout.Reset(timeout)
}

func (this *ZabNode) getNextIndexInRingByIndex(prevIndex uint32) uint32 { //NodeId <- Self
	if prevIndex+1 >= uint32(len(this.nodesIdRing)) {
		return 0
	} else {
		return prevIndex + 1
	}
}

func (this *ZabNode) getBackIndexInRingByIndex(nextIndex uint32) uint32 { //NodeId <- Self
	if nextIndex-1 < 0 {
		return this.nodesIdRing[len(this.nodesIdRing)-1]
	} else {
		return nextIndex - 1
	}
}

func (this *ZabNode) getRingIndexByNodeId(nodeIdToSearch uint32) uint32 { //NodeId <- Self
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
