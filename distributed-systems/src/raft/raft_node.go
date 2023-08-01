package raft

import (
	"distributed-systems/src/broadcast"
	"distributed-systems/src/broadcast/fifo"
	"distributed-systems/src/nodes/types"
	"time"
)

type RaftNode struct {
	broadcast.Node

	leaderNodeId uint32

	state RaftState

	electionsByTerm map[uint32]*RaftElection

	heartbeatTimeoutTimer *time.Timer
	heartbeatTimeoutMs    time.Duration

	heartbeatsTicker  *time.Ticker
	heartbeatTickerMs time.Duration

	electionTimeoutMs time.Duration

	term uint32
}

func CreateRaftNode(heartbeatTimeoutMs uint64, heartbeatTickerMs uint64, electionTimeoutMs uint64, leaderNodeId uint32, nodeId uint32, port uint16) *RaftNode {
	var state RaftState

	if leaderNodeId == nodeId {
		state = LEADER
	} else {
		state = FOLLOWER
	}

	raftNode := &RaftNode{
		Node:               *broadcast.CreateNode(nodeId, port, fifo.CreateFifoBroadcaster(3, 6, nodeId)),
		leaderNodeId:       leaderNodeId,
		electionsByTerm:    map[uint32]*RaftElection{},
		heartbeatTimeoutMs: time.Duration(heartbeatTimeoutMs * uint64(time.Millisecond)),
		heartbeatTickerMs:  time.Duration(heartbeatTickerMs * uint64(time.Millisecond)),
		electionTimeoutMs:  time.Duration(electionTimeoutMs * uint64(time.Millisecond)),
		state:              state,
	}

	raftNode.AddMessageHandler(types.MESSAGE_RAFT_REQUEST_ELECTION, raftNode.handleRequestElection)
	raftNode.AddMessageHandler(types.MESSAGE_RAFT_REQUEST_ELECTION_REJECTED_OUTDATED_TERM, raftNode.handleRequestElectionOutdatedTerm)
	raftNode.AddMessageHandler(types.MESSAGE_RAFT_REQUEST_ELECTION_VOTED, raftNode.handleRequestElectionVoted)
	raftNode.AddMessageHandler(types.MESSAGE_RAFT_REQUEST_ELECTION_NODE_ELECTED, raftNode.handleElectionNodeElected)

	return raftNode
}

func (this *RaftNode) Start() {
	if this.IsLeader() {
		this.startLeader()
	} else {
		this.startFollower()
	}
}

func (this *RaftNode) IsLeader() bool {
	return this.GetNodeId() == this.leaderNodeId
}

func (this *RaftNode) getElectionOrCreate(newTerm uint32) *RaftElection {
	if _, contained := this.electionsByTerm[newTerm]; !contained {
		this.electionsByTerm[newTerm] = CreateRaftElection(this.electionTimeoutMs, newTerm, this.onElectionTimeout)
	}

	election, _ := this.electionsByTerm[newTerm]
	return election
}
