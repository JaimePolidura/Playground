package log

type RaftLogEntry struct {
	value uint32
	term  uint32
	index uint32
}
