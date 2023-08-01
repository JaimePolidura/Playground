package raft

import (
	"distributed-systems/src/nodes"
	"distributed-systems/src/nodes/types"
	"time"
)

func (this *RaftNode) handleHeartbeatTimeout() {
	for {
		select {
		case <-this.heartbeatTimeoutTimer.C:
			if this.state == FOLLOWER {
				this.startElection()
			}
		}
	}
}

func (this *RaftNode) onElectionTimeout(termElectionTimeout uint32) {
	if termElectionTimeout == this.term {
		this.startElection()
	}
}

func (this *RaftNode) startElection() {
	this.stopFollower()

	this.state = CANDIDATE
	nextTerm := this.term + 1
	this.term = nextTerm

	this.startElectionTimeout(nextTerm)
	this.GetConnectionManager().SendAllExcept(this.leaderNodeId, nodes.CreateMessage(
		nodes.WithNodeId(this.GetNodeId()),
		nodes.WithType(types.MESSAGE_RAFT_REQUEST_ELECTION),
		nodes.WithContentUInt32(nextTerm)))
}

func (this *RaftNode) handleRequestElection(message *nodes.Message) {
	requestElectionTerm := message.GetContentToUint32()
	requestElectionCandidate := message.NodeIdSender

	if requestElectionTerm <= this.term {
		this.GetConnectionManager().Send(requestElectionCandidate, nodes.CreateMessage(
			nodes.WithNodeId(this.GetNodeId()),
			nodes.WithType(types.MESSAGE_RAFT_REQUEST_ELECTION_REJECTED_OUTDATED_TERM),
			nodes.WithContentUInt32(this.term)))
		return
	}

	election := this.getElectionOrCreate(requestElectionTerm)

	if this.IsLeader() {
		this.stopLeader()
	}
	if election.HaveVoted() {
		this.GetConnectionManager().Send(requestElectionCandidate, nodes.CreateMessage(
			nodes.WithNodeId(this.GetNodeId()),
			nodes.WithType(types.MESSAGE_RAFT_REQUEST_ELECTION_REJECTED_ALREADY_VOTED),
			nodes.WithContentUInt32(election.GetNodeIdVotedFor())))
		return
	}

	election.RegisterVoteFor(requestElectionCandidate)

	this.term = requestElectionTerm

	this.GetConnectionManager().Send(requestElectionCandidate, nodes.CreateMessage(
		nodes.WithNodeId(this.GetNodeId()),
		nodes.WithType(types.MESSAGE_RAFT_REQUEST_ELECTION_VOTED),
		nodes.WithContentUInt32(requestElectionTerm)))
}

func (this *RaftNode) handleRequestElectionRejectedOutdatedTerm(message *nodes.Message) {
	newTerm := message.GetContentToUint32()

	this.term = newTerm
	this.state = FOLLOWER
}

func (this *RaftNode) handleRequestElectionVoted(message *nodes.Message) {
	term := message.GetContentToUint32()
	election := this.getElectionOrCreate(term)

	if term <= this.term {
		this.GetConnectionManager().Send(message.NodeIdSender, nodes.CreateMessage(
			nodes.WithNodeId(this.GetNodeId()),
			nodes.WithType(types.MESSAGE_RAFT_REQUEST_ELECTION_REJECTED_OUTDATED_TERM),
			nodes.WithContentUInt32(this.term)))
		return
	}

	election.RegisterVoteForMe(message.NodeIdSender)
	nNodesVotedForMe := election.GetNodesVotedForMe()
	quorumSatisfied := nNodesVotedForMe >= this.GetConnectionManager().GetNumberConnections()/2+1

	if quorumSatisfied && election.IsOnGoing() {
		election.Finish()

		this.leaderNodeId = this.GetNodeId()
		this.startLeader()

		this.GetConnectionManager().SendAll(nodes.CreateMessage(
			nodes.WithNodeId(this.GetNodeId()),
			nodes.WithType(types.MESSAGE_RAFT_REQUEST_ELECTION_NODE_ELECTED),
			nodes.WithContentsUInt32(this.GetNodeId(), term)))
	}
}

func (this *RaftNode) handleElectionNodeElected(message *nodes.Message) {
	term := message.GetContentToUint32()
	election := this.getElectionOrCreate(term)
	candidateChosen := message.NodeIdSender

	election.Finish()

	if term <= this.term {
		this.GetConnectionManager().Send(message.NodeIdSender, nodes.CreateMessage(
			nodes.WithNodeId(this.GetNodeId()),
			nodes.WithType(types.MESSAGE_RAFT_REQUEST_ELECTION_REJECTED_OUTDATED_TERM),
			nodes.WithContentUInt32(this.term)))
		return
	}

	this.leaderNodeId = candidateChosen

	this.startFollower()
}

func (this *RaftNode) handleRequestElectionOutdatedTerm(message *nodes.Message) {
	this.term = message.GetContentToUint32()
}

func (this *RaftNode) handleHeartbeatMessage(message *nodes.Message) {
	this.heartbeatTimeoutTimer.Reset(this.heartbeatTimeoutMs)
}

func (this *RaftNode) startElectionTimeout(newTerm uint32) {
	this.getElectionOrCreate(newTerm) //Will initialize the timer
}

func (this *RaftNode) stopFollower() {
	this.heartbeatTimeoutTimer.Stop()
}

func (this *RaftNode) startFollower() {
	if this.heartbeatTimeoutTimer != nil {
		this.heartbeatTimeoutTimer.Reset(this.heartbeatTickerMs)
	} else {
		this.heartbeatTimeoutTimer = time.NewTimer(this.heartbeatTickerMs)
		go this.handleHeartbeatTimeout()
	}

	this.state = FOLLOWER
}
