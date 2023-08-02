package raft

import (
	"distributed-systems/src/broadcast"
	"distributed-systems/src/broadcast/fifo"
	"distributed-systems/src/nodes/types"
	"distributed-systems/src/raft/elections"
	"distributed-systems/src/raft/log"
	"time"
)

type RaftNode struct {
	broadcast.Node

	leaderNodeId uint32
	currentTerm  uint32
	state        RaftState

	log *log.RaftLog

	electionsByTerm map[uint32]*elections.RaftElection

	heartbeatTimeoutTimer *time.Timer
	heartbeatTimeoutMs    time.Duration

	heartbeatsTicker  *time.Ticker
	heartbeatTickerMs time.Duration

	electionTimeoutMs time.Duration

	onConsensusCallback func(entries []*log.RaftLogEntry)
}

func CreateRaftNode(heartbeatTimeoutMs uint64, heartbeatTickerMs uint64, electionTimeoutMs uint64, leaderNodeId uint32, nodeId uint32, port uint16) *RaftNode {
	var state RaftState

	if leaderNodeId == nodeId {
		state = LEADER
	} else {
		state = FOLLOWER
	}

	fifoBroadcaster := fifo.CreateFifoBroadcaster(3, 6, nodeId)
	fifoBroadcaster.DisableGossiping()
	fifoBroadcaster.DisableLogging()

	raftNode := &RaftNode{
		Node:               *broadcast.CreateNode(nodeId, port, fifoBroadcaster),
		leaderNodeId:       leaderNodeId,
		electionsByTerm:    map[uint32]*elections.RaftElection{},
		heartbeatTimeoutMs: time.Duration(heartbeatTimeoutMs * uint64(time.Millisecond)),
		heartbeatTickerMs:  time.Duration(heartbeatTickerMs * uint64(time.Millisecond)),
		electionTimeoutMs:  time.Duration(electionTimeoutMs * uint64(time.Millisecond)),
		state:              state,
		log:                log.CreateRaftLog(),
	}

	raftNode.AddMessageHandler(types.MESSAGE_RAFT_REQUEST_ELECTION, raftNode.handleRequestElection)
	raftNode.AddMessageHandler(types.MESSAGE_RAFT_OUTDATED_TERM, raftNode.handleRequestElectionOutdatedTerm)
	raftNode.AddMessageHandler(types.MESSAGE_RAFT_REQUEST_ELECTION_VOTED, raftNode.handleRequestElectionVoted)
	raftNode.AddMessageHandler(types.MESSAGE_RAFT_REQUEST_ELECTION_NODE_ELECTED, raftNode.handleElectionNodeElected)
	raftNode.AddMessageHandler(types.MESSAGE_HEARTBEAT, raftNode.handleHeartbeatMessage)

	raftNode.AddMessageHandler(types.MESSAGE_RAFT_LOG_APPEND_ENTRIES, raftNode.handleAppendEntries)
	raftNode.AddMessageHandler(types.MESSAGE_RAFT_LOG_APPENDED_ENTRY, raftNode.handleAppendedEntry)
	raftNode.AddMessageHandler(types.MESSAGE_RAFT_LOG_DO_COMMIT, raftNode.handleDoCommit)
	raftNode.AddMessageHandler(types.MESSAGE_RAFT_LOG_OUTDATED_ENTRIES, raftNode.handleOutdatedEntries)

	return raftNode
}

func (this *RaftNode) SetOnConsensusCallback(callback func(entries []*log.RaftLogEntry)) {
	this.onConsensusCallback = callback
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

func (this *RaftNode) IsFollower() bool {
	return this.GetNodeId() != this.leaderNodeId
}

func (this *RaftNode) getElectionOrCreate(newTerm uint32) *elections.RaftElection {
	if _, contained := this.electionsByTerm[newTerm]; !contained {
		this.electionsByTerm[newTerm] = elections.CreateRaftElection(this.electionTimeoutMs, newTerm, this.onElectionTimeout)
	}

	election, _ := this.electionsByTerm[newTerm]
	return election
}
