package messages

type RequestVoteRequest struct {
	Term         uint64
	CandidateId  uint32
	LastLogIndex int32
	LastLogTerm  uint64
}

type RequestVoteResponse struct {
	Term        uint64
	VoteGranted bool
}
