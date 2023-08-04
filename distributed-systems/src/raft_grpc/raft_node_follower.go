package raft_grpc

import (
	"distributed-systems/src/raft"
	"time"
)

func (this *RaftNode) startFollower() {
	if this.heartbeatTimeoutTimer != nil {
		this.heartbeatTimeoutTimer.Reset(this.heartbeatTimeoutMs)
	} else {
		this.heartbeatTimeoutTimer = time.NewTimer(this.heartbeatTimeoutMs)
		go this.handleHeartbeatTimeout()
	}

	this.state = raft.FOLLOWER
}

func (this *RaftNode) stopFollower() {
	this.heartbeatTimeoutTimer.Stop()
}
