package raft_grpc

import (
	"context"
	"distributed-systems/src/raft"
	"distributed-systems/src/raft_grpc/messages"
	"time"
)

func (this *RaftNode) startLeader() {
	this.state = raft.LEADER
	if this.heartbeatsTicker == nil {
		this.heartbeatsTicker = time.NewTicker(this.heartbeatTickerMs)
	}

	go this.startSendingHeartbeats()
}

func (this *RaftNode) stopLeader() {
	this.state = raft.FOLLOWER
	this.heartbeatsTicker.Stop()
}

func (this *RaftNode) startSendingHeartbeats() {
	for {
		<-this.heartbeatsTicker.C

		if this.state == raft.LEADER {
			for _, peer := range this.peers {
				peer.RaftNodeService.ReceiveLeaderHeartbeat(context.Background(), &messages.HeartbeatRequest{
					SenderNodeId: this.NodeId,
					Term:         this.currentTerm,
				})
			}
		}
	}
}
