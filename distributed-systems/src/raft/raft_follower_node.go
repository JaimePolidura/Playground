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
	if termElectionTimeout == this.currentTerm {
		this.startElection()
	}
}

func (this *RaftNode) startElection() {
	this.stopFollower()

	this.state = CANDIDATE
	nextTerm := this.currentTerm + 1
	this.currentTerm = nextTerm

	fmt.Printf("[%d] Failure of leader %d detected. Starting election with currentTerm %d Proposing my self as the new leader\n",
		this.GetNodeId(), this.leaderNodeId, nextTerm)

	this.startElectionTimeout(nextTerm)
	this.Broadcast(nodes.CreateMessage(
		nodes.WithNodeId(this.GetNodeId()),
		nodes.WithType(types.MESSAGE_RAFT_REQUEST_ELECTION),
		nodes.WithFlags(types.FLAG_BROADCAST),
		nodes.WithContentUInt32(nextTerm)))
}

func (this *RaftNode) handleRequestElection(message *nodes.Message) {
	requestElectionTerm := message.GetContentToUint32()
	requestElectionCandidate := message.NodeIdSender

	if requestElectionTerm <= this.currentTerm {
		fmt.Printf("[%d] Ignoring REQUEST_ELECTION of node %d with currentTerm %d Outdated currentTerm %d \n",
			this.GetNodeId(), requestElectionCandidate, this.currentTerm, requestElectionTerm)

		this.GetConnectionManager().Send(requestElectionCandidate, nodes.CreateMessage(
			nodes.WithNodeId(this.GetNodeId()),
			nodes.WithType(types.MESSAGE_RAFT_OUTDATED_TERM),
			nodes.WithContentUInt32(this.currentTerm)))
		return
	}

	election := this.getElectionOrCreate(requestElectionTerm)

	if this.IsLeader() {
		this.stopLeader()
	}
	if election.HaveVoted() {
		fmt.Printf("[%d] Ignoring REQUEST_ELECTION of node %d with currentTerm %d Already voted for that currentTerm\n",
			this.GetNodeId(), requestElectionCandidate, requestElectionTerm, this.currentTerm)

		this.GetConnectionManager().Send(requestElectionCandidate, nodes.CreateMessage(
			nodes.WithNodeId(this.GetNodeId()),
			nodes.WithType(types.MESSAGE_RAFT_REQUEST_ELECTION_REJECTED_ALREADY_VOTED),
			nodes.WithContentUInt32(election.GetNodeIdVotedFor())))
		return
	}

	election.RegisterVoteFor(requestElectionCandidate)

	this.currentTerm = requestElectionTerm

	fmt.Printf("[%d] Received REQUEST_ELECTION Voting for node %d in currentTerm %d\n",
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

	fmt.Printf("[%d] Recived OUTDATED_TERM Update my currentTerm to %d Prev: %d\n",
		this.GetNodeId(), newTerm, this.currentTerm)

	this.currentTerm = newTerm
	this.state = FOLLOWER
}

func (this *RaftNode) handleRequestElectionVoted(message *nodes.Message) {
	termVote := message.GetContentToUint32()
	election := this.getElectionOrCreate(termVote)

	if termVote < this.currentTerm {
		fmt.Printf("[%d] Ignoring ELECTION_VOTED of node %d with currentTerm %d Outdated currentTerm %d \n",
			this.GetNodeId(), message.NodeIdSender, termVote, this.currentTerm)

		this.GetConnectionManager().Send(message.NodeIdSender, nodes.CreateMessage(
			nodes.WithNodeId(this.GetNodeId()),
			nodes.WithType(types.MESSAGE_RAFT_OUTDATED_TERM),
			nodes.WithContentUInt32(this.currentTerm)))
		return
	}

	election.RegisterVoteForMe(message.NodeIdSender)
	nNodesVotedForMe := election.GetNodesVotedForMe()
	quorumSatisfied := nNodesVotedForMe >= this.GetConnectionManager().GetNumberConnections()/2+1

	fmt.Printf("[%d] Recevied ELECTION_VOTED from node %d of currentTerm %d. Nº Nodes voted for me %d Is quorum satisfied? %t\n",
		this.GetNodeId(), message.NodeIdSender, termVote, nNodesVotedForMe, quorumSatisfied)

	if quorumSatisfied && election.IsOnGoing() {
		election.Finish()

		this.leaderNodeId = this.GetNodeId()
		this.startLeader()

		fmt.Printf("[%d] Leader election finished! New leader %d established in currentTerm %d Sending NODE_ELECTED to followers\n",
			this.GetNodeId(), this.GetNodeId(), this.currentTerm)

		this.Broadcast(nodes.CreateMessage(
			nodes.WithNodeId(this.GetNodeId()),
			nodes.WithType(types.MESSAGE_RAFT_REQUEST_ELECTION_NODE_ELECTED),
			nodes.WithFlags(types.FLAG_BROADCAST),
			nodes.WithContentsUInt32(this.currentTerm)))
	}
}

func (this *RaftNode) handleElectionNodeElected(message *nodes.Message) {
	term := message.GetContentToUint32()
	election := this.getElectionOrCreate(term)
	candidateChosen := message.NodeIdSender

	election.Finish()

	if term < this.currentTerm {
		fmt.Printf("[%d] Ignoring NODE_ELECTED of node %d with currentTerm %d Outdated currentTerm %d \n",
			this.GetNodeId(), message.NodeIdSender, term, this.currentTerm)

		this.GetConnectionManager().Send(message.NodeIdSender, nodes.CreateMessage(
			nodes.WithNodeId(this.GetNodeId()),
			nodes.WithType(types.MESSAGE_RAFT_OUTDATED_TERM),
			nodes.WithContentUInt32(this.currentTerm)))
		return
	}

	fmt.Printf("[%d] Established new leader %d in currentTerm %d\n",
		this.GetNodeId(), candidateChosen, term)

	this.leaderNodeId = candidateChosen

	this.startFollower()
}

func (this *RaftNode) handleRequestElectionOutdatedTerm(message *nodes.Message) {
	this.currentTerm = message.GetContentToUint32()
	if this.IsLeader() {
		this.stopLeader()
	}
}

func (this *RaftNode) handleAppendEntries(message *nodes.Message) {
	termLeader := message.GetContentToUint32WithOffset(0)
	prevIndexLeader := message.GetContentToUint32WithOffset(4)
	prevTermLeader := message.GetContentToUint32WithOffset(8)
	entryLeaderToAppend := message.GetContentToUint32WithOffset(12)

	if termLeader < this.currentTerm {
		fmt.Printf("[%d] Ignoring APPEND_ENTRIES of node %d with currentTerm %d Outdated currentTerm %d \n",
			this.GetNodeId(), message.NodeIdSender, termLeader, this.currentTerm)

		this.GetConnectionManager().Send(message.NodeIdSender, nodes.CreateMessage(
			nodes.WithNodeId(this.GetNodeId()),
			nodes.WithType(types.MESSAGE_RAFT_OUTDATED_TERM),
			nodes.WithContentUInt32(this.currentTerm)))
		return
	}

	if this.log.HasIndex(prevIndexLeader) && this.log.GetTermByIndex(prevIndexLeader) != prevTermLeader {
		fmt.Printf("[%d] Ignoring APPEND_ENTRIES from node %d with term %d Log terms at index %d doest match, my term at that index: %d leader term %d\n",
			this.GetNodeId(), message.NodeIdSender, termLeader, prevIndexLeader, this.log.GetTermByIndex(prevIndexLeader), prevTermLeader)

		this.GetConnectionManager().Send(message.NodeIdSender, nodes.CreateMessage(
			nodes.WithNodeId(this.GetNodeId()),
			nodes.WithType(types.MESSAGE_RAFT_LOG_OUTDATED_ENTRIES),
			nodes.WithContentUInt32(this.log.GetLastCommittedIndex())))
		return
	}

	fmt.Printf("[%d] Received APPEND_ENTRIES from node %d with prev term %d, prev index %d, new index %d and value %d Sending back MESSAGE_RAFT_LOG_APPENDED_ENTRY\n",
		this.GetNodeId(), message.NodeIdSender, prevTermLeader, prevIndexLeader, prevIndexLeader+1, entryLeaderToAppend)

	this.log.AddUncommittedEntry(entryLeaderToAppend, this.currentTerm, prevIndexLeader+1)

	this.GetConnectionManager().Send(message.NodeIdSender, nodes.CreateMessage(
		nodes.WithNodeId(this.GetNodeId()),
		nodes.WithType(types.MESSAGE_RAFT_LOG_APPENDED_ENTRY)))
}

func (this *RaftNode) handleDoCommit(message *nodes.Message) {
	committedEntries := this.log.Commit()

	fmt.Printf("[%d] Recevied DO_COMMIT Commiting entries %d The actual log %v\n",
		this.GetNodeId(), len(committedEntries), this.log.GetCommittedLogEntries())

	if this.onConsensusCallback != nil {
		this.onConsensusCallback(committedEntries)
	}
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
