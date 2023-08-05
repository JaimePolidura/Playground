package raft

import (
	"context"
	"distributed-systems/src/raft/elections"
	"distributed-systems/src/raft/log"
	"distributed-systems/src/raft/messages"
	"sync"
	"time"
)

type RaftNodeService interface {
	RequestVote(context.Context, *messages.RequestVoteRequest) *messages.RequestVoteResponse
	AppendEntries(context.Context, *messages.AppendEntriesRequest) *messages.AppendEntriesResponse
}

type Peer struct {
	RaftNodeService RaftNodeService
	NodeId          uint32
}

type RaftNode struct {
	NodeId uint32
	Port   uint16
	state  RaftState
	log    *log.RaftLog

	leaderNodeId    uint32
	currentTerm     uint64
	electionsByTerm map[uint64]*elections.RaftElection
	peers           []*Peer
	commitIndex     int32

	//nextIndex  map[uint32]int32
	//matchIndex map[uint32]int32

	nextIndex  *sync.Map
	matchIndex *sync.Map

	heartbeatTimeoutTimer *time.Timer
	heartbeatTimeoutMs    time.Duration

	heartbeatsTicker  *time.Ticker
	heartbeatTickerMs time.Duration

	electionTimeoutMs time.Duration
}

func CreateRaftNode(nodeId uint32, leaderNodeId uint32, port uint16, heartbeatTimeoutMs uint64, heartbeatTickerMs uint64, electionTimeoutMs uint64) *RaftNode {
	var state RaftState

	if leaderNodeId == nodeId {
		state = LEADER
	} else {
		state = FOLLOWER
	}

	return &RaftNode{
		heartbeatTimeoutMs: time.Duration(heartbeatTimeoutMs * uint64(time.Millisecond)),
		heartbeatTickerMs:  time.Duration(heartbeatTickerMs * uint64(time.Millisecond)),
		electionTimeoutMs:  time.Duration(electionTimeoutMs * uint64(time.Millisecond)),
		electionsByTerm:    map[uint64]*elections.RaftElection{},
		nextIndex:          new(sync.Map),
		matchIndex:         new(sync.Map),
		//matchIndex:         map[uint32]int32{},
		//nextIndex:          map[uint32]int32{},
		log:          log.CreateRaftLog(),
		peers:        make([]*Peer, 0),
		leaderNodeId: leaderNodeId,
		NodeId:       nodeId,
		state:        state,
		Port:         port,
	}
}

func (this *RaftNode) AddPeers(peers []*Peer) {
	for _, peer := range peers {
		if peer.NodeId != this.NodeId {
			this.peers = append(this.peers, peer)
		}
	}
}

func (this *RaftNode) Stop() {
	if this.heartbeatsTicker != nil {
		this.heartbeatsTicker.Stop()
	}
	if this.heartbeatTimeoutTimer != nil {
		this.heartbeatTimeoutTimer.Stop()
	}
	if this.heartbeatTimeoutTimer != nil {
		this.heartbeatTimeoutTimer.Stop()
	}
}

func (this *RaftNode) Start() {
	if this.IsLeader() {
		this.startLeader()
	} else {
		this.startFollower()
	}
}

func (this *RaftNode) IsLeader() bool {
	return this.NodeId == this.leaderNodeId
}

func (this *RaftNode) IsFollower() bool {
	return this.NodeId != this.leaderNodeId
}
