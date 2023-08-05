package messages

import "distributed-systems/src/raft/log"

type AppendEntriesRequest struct {
	Term         uint64
	LeaderId     uint32
	PrevLogIndex int32
	PrevLogTerm  uint64
	Entries      []log.RaftLogEntry
	LeaderCommit int32
}

type AppendEntriesResponse struct {
	Term    uint64
	Success bool
}
