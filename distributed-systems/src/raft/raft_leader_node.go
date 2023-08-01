package raft

import (
	"distributed-systems/src/nodes"
	"distributed-systems/src/nodes/types"
	"time"
)

func (this *RaftNode) startLeader() {
	this.state = LEADER

	this.setupHeartbeatsTickerLeader()
}

func (this *RaftNode) stopLeader() {
	this.state = FOLLOWER

	this.heartbeatsTicker.Stop()
}

func (this *RaftNode) setupHeartbeatsTickerLeader() {
	if this.heartbeatsTicker == nil {
		this.heartbeatsTicker = time.NewTicker(this.heartbeatTickerMs)
		go this.startSendingHeartbeats()
	} else {
		this.heartbeatsTicker.Reset(this.heartbeatTickerMs)
	}
}

func (this *RaftNode) startSendingHeartbeats() {
	for {
		select {
		case <-this.heartbeatsTicker.C:
			if this.state == LEADER {
				this.GetConnectionManager().SendAllExcept(this.GetNodeId(), nodes.CreateMessage(
					nodes.WithNodeId(this.GetNodeId()),
					nodes.WithType(types.MESSAGE_HEARTBEAT)))
			}
		}
	}
}
