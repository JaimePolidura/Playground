package elections

import (
	"time"
)

type RaftElection struct {
	term       uint64
	isFinished bool

	haveVoted      bool
	nodeIdVotedFor uint32

	//Implement listener
	electionTimeoutTimer *time.Timer
	heartbeatTimeoutMs   time.Duration
	onTimeoutCallback    func(term uint64)

	nodesVotedForMe map[uint32]uint32
}

func CreateRaftElection(electionTimeoutMs time.Duration, term uint64, onTimeoutCallback func(term uint64)) *RaftElection {
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

func (this *RaftElection) HaveNodeVotedForMe(otherNodeId uint32) bool {
	_, contains := this.nodesVotedForMe[otherNodeId]
	return contains
}

func (this *RaftElection) HaveIVoted() bool {
	return this.haveVoted
}

func (this *RaftElection) Finish() {
	this.isFinished = true

	if this.electionTimeoutTimer != nil {
		this.electionTimeoutTimer.Stop()
	}
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
