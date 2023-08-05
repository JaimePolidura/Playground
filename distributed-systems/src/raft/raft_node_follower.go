package raft

import (
	"context"
	"distributed-systems/src/raft/messages"
	"distributed-systems/src/utils"
	"time"
)

func (this *RaftNode) AppendEntries(context context.Context, request *messages.AppendEntriesRequest) *messages.AppendEntriesResponse {
	if this.currentTerm > request.Term {
		this.updateOutdatedTerm(request.Term)
	}
	if request.Term < this.currentTerm {
		return &messages.AppendEntriesResponse{Term: this.currentTerm, Success: false}
	}
	if this.IsLeader() {
		this.stopLeader()
	}
	if request.PrevLogIndex != -1 && !this.log.HasIndex(request.PrevLogIndex) || this.log.GetTermByIndex(request.PrevLogIndex) != request.PrevLogTerm {
		return &messages.AppendEntriesResponse{Term: this.currentTerm, Success: false}
	}
	if this.heartbeatTimeoutTimer != nil {
		this.heartbeatTimeoutTimer.Reset(this.heartbeatTimeoutMs)
	}

	logInsertIndex := request.PrevLogIndex + 1
	newEntriesIndex := int32(0)

	for {
		if logInsertIndex >= this.log.Size() || newEntriesIndex >= this.log.Size() {
			break
		}
		if this.log.GetTermByIndex(logInsertIndex) != this.log.GetTermByIndex(newEntriesIndex) {
			break
		}
		logInsertIndex++
		newEntriesIndex++
	}

	if newEntriesIndex < this.log.Size() {
		this.log.Entries = append(this.log.Entries[:logInsertIndex], request.Entries[newEntriesIndex:]...)
	}

	if request.LeaderCommit > this.commitIndex {
		this.commitIndex = utils.MinInt32(request.LeaderCommit, this.log.GetLastIndex())
	}

	return &messages.AppendEntriesResponse{
		Term:    this.currentTerm,
		Success: true,
	}
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

func (this *RaftNode) stopFollower() {
	this.heartbeatTimeoutTimer.Stop()
}
