package log

type RaftLogEntry struct {
	value uint32
	term  uint64
	index uint32
}
