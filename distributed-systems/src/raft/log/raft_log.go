package log

type RaftLog struct {
	Entries     []RaftLogEntry
	commitIndex int32
	lastAdded   int32
}

type RaftLogEntry struct {
	Value uint32
	Term  uint64
	Index int32
}

func (this *RaftLog) Size() int32 {
	return int32(len(this.Entries))
}

func (this *RaftLog) GetFromIndexExclusive(index int32) []RaftLogEntry {
	if len(this.Entries) > 0 {
		return this.Entries[index+1:]
	} else {
		return []RaftLogEntry{}
	}
}

func (this *RaftLog) AddEntry(value uint32, term uint64) {
	this.Entries = append(this.Entries, RaftLogEntry{Term: term, Value: value, Index: int32(len(this.Entries))})
}

func (this *RaftLog) AddEntries(entries []RaftLogEntry) int32 {
	indexToAdd := int32(len(this.Entries))

	this.Entries = this.Entries[:this.commitIndex]
	this.Entries = append(this.Entries, entries...)

	return indexToAdd
}

func (this *RaftLog) HasIndex(index int32) bool {
	return len(this.Entries) > 0 && int32(len(this.Entries)-1) <= index
}

func (this *RaftLog) GetTermByIndex(index int32) uint64 {
	if len(this.Entries) == 0 {
		return 0
	} else {
		return this.Entries[len(this.Entries)-1].Term
	}
}

func (this *RaftLog) GetLastTerm() uint64 {
	if len(this.Entries) == 0 {
		return 0
	} else {
		return this.Entries[len(this.Entries)-1].Term
	}
}

func (this *RaftLog) GetLastIndex() int32 {
	if len(this.Entries) == 0 {
		return -1
	} else {
		return int32(len(this.Entries)) - 1
	}
}

func (this *RaftLog) GetNextIndex() int32 {
	return this.GetLastIndex() + 1
}

func CreateRaftLog() *RaftLog {
	return &RaftLog{
		Entries: make([]RaftLogEntry, 0),
	}
}
