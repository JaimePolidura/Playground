package raft

import (
	"distributed-systems/src/nodes"
	"distributed-systems/src/nodes/types"
	"fmt"
	"time"
)

func (this *RaftNode) startSendingHeartbeats() {
	for {
		select {
		case <-this.heartbeatsTicker.C:
			if this.state == LEADER {
				this.Broadcast(nodes.CreateMessage(
					nodes.WithNodeId(this.GetNodeId()),
					nodes.WithFlags(types.FLAG_BROADCAST),
					nodes.WithType(types.MESSAGE_HEARTBEAT)))
			}
		}
	}
}

func (this *RaftNode) AppendEntries(entryValue uint32) {
	if this.IsFollower() {
		return
	}

	index, term := this.log.GetLastCommitted()

	this.log.AddUncommittedEntry(entryValue, this.currentTerm, index+1)

	this.Broadcast(nodes.CreateMessage(
		nodes.WithNodeId(this.GetNodeId()),
		nodes.WithType(types.MESSAGE_RAFT_LOG_APPEND_ENTRIES),
		nodes.WithFlags(types.FLAG_BROADCAST),
		nodes.WithContentsUInt32(
			uint32(this.currentTerm),
			index,
			uint32(term),
			1,
			entryValue)))
}

func (this *RaftNode) handleOutdatedEntries(message *nodes.Message) {
	fmt.Printf("[%d] Received OUTDATED_ENTRIES from node %d Node last index commited %d My last index commited %d\n",
		this.GetNodeId(), message.GetContentToUint32(), this.log.GetLastCommittedIndex())

	//Cannot implement the message system that I use doest allow sending arrays of objects. It would be very complex
}

func (this *RaftNode) handleAppendedEntry(message *nodes.Message) {
	this.log.NNodesAppendedLog++

	fmt.Printf("[%d] Recevied APPENDED_ENTRY from node %d Nº nodes appended: %d Is quorum satisfied? %t\n",
		this.GetNodeId(), message.NodeIdSender, this.log.NNodesAppendedLog+1, this.log.NNodesAppendedLog+1 >= this.GetConnectionManager().GetNumberConnections()/2+1)

	if !this.log.QuorumNodesAppendedSatisfied && this.log.NNodesAppendedLog+1 >= this.GetConnectionManager().GetNumberConnections()/2+1 {
		this.log.QuorumNodesAppendedSatisfied = true
		this.Broadcast(nodes.CreateMessage(nodes.WithNodeId(this.GetNodeId()),
			nodes.WithFlags(types.FLAG_BROADCAST),
			nodes.WithType(types.MESSAGE_RAFT_LOG_DO_COMMIT)))

		committedEntries := this.log.Commit()

		if this.onConsensusCallback != nil {
			this.onConsensusCallback(committedEntries)
		}

		fmt.Printf("[%d] Commiting %d entries The actual log %v\n",
			this.GetNodeId(), len(committedEntries), this.log.GetCommittedLogEntries())
	}
}

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
