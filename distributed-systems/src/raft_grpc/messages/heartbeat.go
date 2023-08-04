package messages

type HeartbeatRequest struct {
	Term         uint64
	SenderNodeId uint32
}
