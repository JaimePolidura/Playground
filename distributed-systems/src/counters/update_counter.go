package counters

type UpdateCounterRequest struct {
	SelfNodeId       uint32
	IsIncrement      bool
	NextSelfSeqValue uint64

	LastSeqValueSeenIncrement uint64
	LastSeqValueSeenDecrement uint64
}

type UpdateCounterResponse struct {
	NeedsSyncIncrement              bool
	NextSelfSeqValueToSyncIncrement uint64

	NeedsSyncDecrement              bool
	NextSelfSeqValueToSyncDecrement uint64
}
