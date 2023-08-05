package raft

import (
	"context"
	"distributed-systems/src/raft/messages"
	"distributed-systems/src/utils"
	"time"
)

func (this *RaftNode) Append(value uint32) {
	if this.IsFollower() {
		panic("wtf bro?!")
	}

	this.log.AddEntry(value, this.currentTerm)
}

func (this *RaftNode) startSendingHeartbeats() {
	for {
		<-this.heartbeatsTicker.C

		if this.state == LEADER {
			this.SendHeartbeats()
		}
	}
}

func (this *RaftNode) SendHeartbeats() {
	for _, peer := range this.peers {
		go func(_peer *Peer) { //CONCURRENT MAP W
			nextIndex := utils.GetInt32FromSyncMap(this.nextIndex, _peer.NodeId)
			prevLogIndex := nextIndex - 1
			prevLogTerm := uint64(0)
			if prevLogIndex >= 0 {
				prevLogTerm = this.log.GetTermByIndex(prevLogIndex)
			}

			entriesToSend := this.log.GetFromIndexExclusive(prevLogIndex)

			response := _peer.RaftNodeService.AppendEntries(context.Background(), &messages.AppendEntriesRequest{
				Term:         this.currentTerm,
				LeaderId:     this.leaderNodeId,
				PrevLogIndex: prevLogIndex,
				PrevLogTerm:  prevLogTerm,
				Entries:      entriesToSend,
				LeaderCommit: this.commitIndex,
			})

			if !response.Success && response.Term > this.currentTerm {
				this.updateOutdatedTerm(response.Term)
				return
			} else if !response.Success {
				this.nextIndex.Store(_peer.NodeId, nextIndex-1)
				return
			}

			//Successful
			this.nextIndex.Store(_peer.NodeId, nextIndex+int32(len(entriesToSend)))
			this.matchIndex.Store(_peer.NodeId, nextIndex+int32(len(entriesToSend))-1)

			initialCommitIndex := this.commitIndex

			for i := this.commitIndex + 1; i < this.log.Size(); i++ {
				if this.log.GetTermByIndex(i) == this.currentTerm {
					followerMatchCount := 1

					for _, peer := range this.peers {
						if utils.GetInt32FromSyncMap(this.matchIndex, peer.NodeId) >= i {
							followerMatchCount++
						}
					}

					if followerMatchCount >= len(this.peers)/2+1 {
						this.commitIndex++
					}
				}
			}

			if initialCommitIndex != this.commitIndex {
				//On consensus
			}
		}(peer)
	}
}

func (this *RaftNode) startLeader() {
	this.state = LEADER
	if this.heartbeatsTicker == nil {
		this.heartbeatsTicker = time.NewTicker(this.heartbeatTickerMs)
	}

	for _, peer := range this.peers {
		this.nextIndex.Store(peer.NodeId, this.log.GetNextIndex())
		this.matchIndex.Store(peer.NodeId, this.log.GetNextIndex()-1)
	}

	go this.startSendingHeartbeats()
}

func (this *RaftNode) stopLeader() {
	this.state = FOLLOWER
	this.heartbeatsTicker.Stop()
}
