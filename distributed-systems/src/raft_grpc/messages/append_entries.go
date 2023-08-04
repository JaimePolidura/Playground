package messages

type AppendEntriesRequest struct {
	Term         uint64
	LeaderId     uint32
	PrevLogIndex uint32
	PrevLogTerm  uint32
	Entries      []Entry
	LeaderCommit uint32
}

type Entry struct {
	Term  uint64
	Index uint32
}

type AppendEntriesResponse struct {
	Term    uint64
	Success bool
}
