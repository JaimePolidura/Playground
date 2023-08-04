package messages

type RequestVoteRequest struct {
	Term         uint64
	CandidateId  uint32
	LastLogIndex uint32
	LastLogTerm  uint32
}

type RequestVoteResponse struct {
	Term        uint64
	VoteGranted bool
}
