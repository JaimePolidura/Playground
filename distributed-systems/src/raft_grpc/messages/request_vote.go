package messages

type RequestVoteRequest struct {
	Term         uint64
	CandidateId  uint32
	LastLogIndex uint32
	LastLogTerm  uint32
}
