package raft

import (
	"distributed-systems/src/nodes"
	"distributed-systems/src/nodes/types"
	"fmt"
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

	fmt.Printf("[%d] Failure of leader %d detected. Starting election with term %d Proposing my self as the new leader\n",
		this.GetNodeId(), this.leaderNodeId, nextTerm)

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
		fmt.Printf("[%d] Ignoring REQUEST_ELECTION of node %d with term %d Outdated term %d \n",
			this.GetNodeId(), requestElectionCandidate, this.term, requestElectionTerm)

		this.GetConnectionManager().Send(requestElectionCandidate, nodes.CreateMessage(
			nodes.WithNodeId(this.GetNodeId()),
			nodes.WithType(types.MESSAGE_RAFT_OUTDATED_TERM),
			nodes.WithContentUInt32(this.term)))
		return
	}

	election := this.getElectionOrCreate(requestElectionTerm)

	if this.IsLeader() {
		this.stopLeader()
	}
	if election.HaveVoted() {
		fmt.Printf("[%d] Ignoring REQUEST_ELECTION of node %d with term %d Already voted for that term\n",
			this.GetNodeId(), requestElectionCandidate, requestElectionTerm, this.term)

		this.GetConnectionManager().Send(requestElectionCandidate, nodes.CreateMessage(
			nodes.WithNodeId(this.GetNodeId()),
			nodes.WithType(types.MESSAGE_RAFT_REQUEST_ELECTION_REJECTED_ALREADY_VOTED),
			nodes.WithContentUInt32(election.GetNodeIdVotedFor())))
		return
	}

	election.RegisterVoteFor(requestElectionCandidate)

	this.term = requestElectionTerm

	fmt.Printf("[%d] Received REQUEST_ELECTION Voting for node %d in term %d\n",
		this.GetNodeId(), requestElectionCandidate, requestElectionTerm)

	this.GetConnectionManager().Send(requestElectionCandidate, nodes.CreateMessage(
		nodes.WithNodeId(this.GetNodeId()),
		nodes.WithType(types.MESSAGE_RAFT_REQUEST_ELECTION_VOTED),
		nodes.WithContentUInt32(requestElectionTerm)))
}

func (this *RaftNode) handleOutdatedTerm(message *nodes.Message) {
	newTerm := message.GetContentToUint32()

	if this.IsLeader() {
		this.stopLeader()
	}

	fmt.Printf("[%d] Recived OUTDATED_TERM Update my term to %d Prev: %d\n",
		this.GetNodeId(), newTerm, this.term)

	this.term = newTerm
	this.state = FOLLOWER
}

func (this *RaftNode) handleRequestElectionVoted(message *nodes.Message) {
	termVote := message.GetContentToUint32()
	election := this.getElectionOrCreate(termVote)

	if termVote < this.term {
		fmt.Printf("[%d] Ignoring ELECTION_VOTED of node %d with term %d Outdated term %d \n",
			this.GetNodeId(), message.NodeIdSender, termVote, this.term)

		this.GetConnectionManager().Send(message.NodeIdSender, nodes.CreateMessage(
			nodes.WithNodeId(this.GetNodeId()),
			nodes.WithType(types.MESSAGE_RAFT_OUTDATED_TERM),
			nodes.WithContentUInt32(this.term)))
		return
	}

	election.RegisterVoteForMe(message.NodeIdSender)
	nNodesVotedForMe := election.GetNodesVotedForMe()
	quorumSatisfied := nNodesVotedForMe >= this.GetConnectionManager().GetNumberConnections()/2+1

	fmt.Printf("[%d] Recevied ELECTION_VOTED from node %d of term %d. Nº Nodes voted for me %d Is quorum satisfied? %t\n",
		this.GetNodeId(), message.NodeIdSender, termVote, nNodesVotedForMe, quorumSatisfied)

	if quorumSatisfied && election.IsOnGoing() {
		election.Finish()

		this.leaderNodeId = this.GetNodeId()
		this.startLeader()

		fmt.Printf("[%d] Leader election finished! New leader %d established in term %d Sending NODE_ELECTED to followers\n",
			this.GetNodeId(), this.GetNodeId(), this.term)

		this.GetConnectionManager().SendAll(nodes.CreateMessage(
			nodes.WithNodeId(this.GetNodeId()),
			nodes.WithType(types.MESSAGE_RAFT_REQUEST_ELECTION_NODE_ELECTED),
			nodes.WithContentsUInt32(this.term)))
	}
}

func (this *RaftNode) handleElectionNodeElected(message *nodes.Message) {
	term := message.GetContentToUint32()
	election := this.getElectionOrCreate(term)
	candidateChosen := message.NodeIdSender

	election.Finish()

	if term < this.term {
		fmt.Printf("[%d] Ignoring NODE_ELECTED of node %d with term %d Outdated term %d \n",
			this.GetNodeId(), message.NodeIdSender, term, this.term)

		this.GetConnectionManager().Send(message.NodeIdSender, nodes.CreateMessage(
			nodes.WithNodeId(this.GetNodeId()),
			nodes.WithType(types.MESSAGE_RAFT_OUTDATED_TERM),
			nodes.WithContentUInt32(this.term)))
		return
	}

	fmt.Printf("[%d] Established new leader %d in term %d\n",
		this.GetNodeId(), candidateChosen, term)

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
		this.heartbeatTimeoutTimer.Reset(this.heartbeatTimeoutMs)
	} else {
		this.heartbeatTimeoutTimer = time.NewTimer(this.heartbeatTimeoutMs)
		go this.handleHeartbeatTimeout()
	}

	this.state = FOLLOWER
}
