package raft_grpc

import (
	"context"
	"distributed-systems/src/raft"
	"distributed-systems/src/raft/elections"
	"distributed-systems/src/raft_grpc/messages"
	"fmt"
)

func (this *RaftNode) startElection() {
	this.stopFollower()

	this.state = raft.CANDIDATE
	nextTerm := this.currentTerm + 1
	this.currentTerm = nextTerm

	fmt.Printf("[%d] Failure of leader %d detected. Starting election with currentTerm %d Proposing my self as the new leader\n",
		this.NodeId, this.leaderNodeId, nextTerm)

	this.startElectionTimeout(nextTerm)

	for _, peer := range this.peers {
		go func(_peer *Peer) {
			this.handlePeerRequestVoteResponse(_peer, _peer.RaftNodeService.RequestVote(context.Background(), &messages.RequestVoteRequest{
				Term:         nextTerm,
				CandidateId:  this.NodeId,
				LastLogIndex: 0,
				LastLogTerm:  0,
			}))
		}(peer)
	}
}

func (this *RaftNode) RequestVote(context context.Context, request *messages.RequestVoteRequest) *messages.RequestVoteResponse {
	if this.currentTerm < request.Term {
		this.updateOutdatedTerm(request.Term)
	}
	if this.currentTerm > request.Term {
		return &messages.RequestVoteResponse{VoteGranted: false, Term: this.currentTerm}
	}
	if this.IsLeader() {
		this.stopLeader()
	}

	election := this.getElectionOrCreate(request.Term)

	if election.HaveIVoted() {
		fmt.Printf("[%d] Ignoring REQUEST_ELECTION of node %d with currentTerm %d Already voted for that currentTerm\n",
			this.NodeId, request.CandidateId, request.Term, this.currentTerm)
		return &messages.RequestVoteResponse{VoteGranted: false, Term: this.currentTerm}
	}

	election.RegisterVoteFor(request.CandidateId)

	this.currentTerm = request.Term

	fmt.Printf("[%d] Received REQUEST_ELECTION Voting for node %d in currentTerm %d\n",
		this.NodeId, request.CandidateId, request.Term)

	return &messages.RequestVoteResponse{VoteGranted: true, Term: request.Term}
}

func (this *RaftNode) handlePeerRequestVoteResponse(peer *Peer, response *messages.RequestVoteResponse) {
	if response.Term < this.currentTerm {
		fmt.Printf("[%d] Ignoring ELECTION_VOTED of node %d with currentTerm %d Outdated currentTerm %d \n",
			this.NodeId, peer.NodeId, response.Term, this.currentTerm)
		return
	}

	if response.VoteGranted {
		this.handlePeerVote(peer, response)
	} else if !response.VoteGranted && response.Term > this.currentTerm {
		fmt.Printf("[%d] Recived OUTDATED_TERM Update my currentTerm to %d Prev: %d\n",
			this.NodeId, response.Term, this.currentTerm)
		this.updateOutdatedTerm(response.Term)
	}
}

func (this *RaftNode) handlePeerVote(peer *Peer, response *messages.RequestVoteResponse) {
	election := this.getElectionOrCreate(this.currentTerm)

	if election.HaveNodeVotedForMe(peer.NodeId) {
		return
	}

	election.RegisterVoteForMe(peer.NodeId)
	nNodesVotedForMe := election.GetNodesVotedForMe()
	quorumSatisfied := nNodesVotedForMe >= uint32(len(this.peers))/2+1

	fmt.Printf("[%d] Recevied vote from node %d of currentTerm %d. Nº Nodes voted for me %d Is quorum satisfied? %t\n",
		this.NodeId, peer.NodeId, response.Term, nNodesVotedForMe, quorumSatisfied)

	if quorumSatisfied && election.IsOnGoing() {
		fmt.Printf("[%d] Leader election finished! New leader %d established in term %d\n",
			this.NodeId, this.NodeId, this.currentTerm)

		election.Finish()
		this.leaderNodeId = this.NodeId
		this.startLeader()

		fmt.Printf("[%d] Leader election finished! New leader %d established in currentTerm %d Sending NODE_ELECTED to followers\n",
			this.NodeId, this.NodeId, this.currentTerm)
	}
}

func (this *RaftNode) ReceiveLeaderHeartbeat(context context.Context, request *messages.HeartbeatRequest) {
	if this.currentTerm > request.Term {
		return
	}
	if this.IsLeader() {
		this.stopLeader()
	}
	if this.currentTerm < request.Term {
		this.currentTerm = request.Term
	}
	if this.leaderNodeId != request.SenderNodeId {
		this.leaderNodeId = request.SenderNodeId
	}
	if this.heartbeatTimeoutTimer != nil {
		this.heartbeatTimeoutTimer.Reset(this.heartbeatTimeoutMs)
	}
}

func (this *RaftNode) handleHeartbeatTimeout() {
	for {
		select {
		case <-this.heartbeatTimeoutTimer.C:
			if this.state == raft.FOLLOWER {
				this.startElection()
			}
		}
	}
}

func (this *RaftNode) onElectionTimeout(termElectionTimeout uint64) {
	if termElectionTimeout == this.currentTerm {
		this.startElection()
	}
}

func (this *RaftNode) updateOutdatedTerm(newTerm uint64) {
	if this.IsLeader() {
		this.stopLeader()
	}

	this.currentTerm = newTerm
	this.state = raft.FOLLOWER
}

func (this *RaftNode) startElectionTimeout(newTerm uint64) {
	this.getElectionOrCreate(newTerm) //Will initialize the timer
}

func (this *RaftNode) getElectionOrCreate(newTerm uint64) *elections.RaftElection {
	if _, contained := this.electionsByTerm[newTerm]; !contained {
		this.electionsByTerm[newTerm] = elections.CreateRaftElection(this.electionTimeoutMs, newTerm, this.onElectionTimeout)
	}

	election, _ := this.electionsByTerm[newTerm]
	return election
}
