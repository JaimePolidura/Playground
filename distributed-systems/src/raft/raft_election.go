package raft

import (
	"time"
)

type RaftElection struct {
	term       uint32
	isFinished bool

	haveVoted      bool
	nodeIdVotedFor uint32

	electionTimeoutTimer *time.Timer
	heartbeatTimeoutMs   time.Duration
	onTimeoutCallback    func(term uint32)

	nodesVotedForMe map[uint32]uint32
}

func CreateRaftElection(electionTimeoutMs time.Duration, term uint32, onTimeoutCallback func(term uint32)) *RaftElection {
	return &RaftElection{
		electionTimeoutTimer: time.NewTimer(electionTimeoutMs),
		heartbeatTimeoutMs:   electionTimeoutMs,
		onTimeoutCallback:    onTimeoutCallback,
		nodesVotedForMe:      make(map[uint32]uint32),
		term:                 term,
	}
}

func (this *RaftElection) GetNodeIdVotedFor() uint32 {
	return this.nodeIdVotedFor
}

func (this *RaftElection) IsFinished() bool {
	return this.isFinished
}

func (this *RaftElection) IsOnGoing() bool {
	return !this.isFinished
}

func (this *RaftElection) HaveVoted() bool {
	return this.haveVoted
}

func (this *RaftElection) Finish() {
	if this.electionTimeoutTimer != nil {
		this.electionTimeoutTimer.Stop()
	}

	this.isFinished = true
}

func (this *RaftElection) RegisterVoteForMe(nodeIdWhoVotedMe uint32) {
	if _, nodeAlreadyVoted := this.nodesVotedForMe[nodeIdWhoVotedMe]; !nodeAlreadyVoted {
		this.nodesVotedForMe[nodeIdWhoVotedMe] = nodeIdWhoVotedMe
	}
}

func (this *RaftElection) RegisterVoteFor(nodeIdToVote uint32) {
	if !this.haveVoted {
		this.nodeIdVotedFor = nodeIdToVote
		this.haveVoted = true
	}
}

func (this *RaftElection) GetNodesVotedForMe() uint32 {
	return uint32(len(this.nodesVotedForMe) + 1)
}
