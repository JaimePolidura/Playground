package log

type RaftLog struct {
	committedEntries []*RaftLogEntry
	commitIndex      uint32

	uncommittedEntries []*RaftLogEntry

	NNodesAppendedLog            uint32
	QuorumNodesAppendedSatisfied bool
}

func (this *RaftLog) Commit() []*RaftLogEntry {
	if len(this.uncommittedEntries) > 0 {
		toCommit := this.uncommittedEntries
		firstUnCommittedEntry := this.uncommittedEntries[len(this.uncommittedEntries)-1]

		this.trimCommittedEntriesLog(firstUnCommittedEntry.index)

		this.committedEntries = append(this.committedEntries, this.uncommittedEntries...)
		this.commitIndex += uint32(len(this.uncommittedEntries))

		this.uncommittedEntries = nil
		this.NNodesAppendedLog = 0
		this.QuorumNodesAppendedSatisfied = false
		this.uncommittedEntries = []*RaftLogEntry{}

		return toCommit
	} else {
		return []*RaftLogEntry{}
	}
}

func (this *RaftLog) AddUncommittedEntry(value uint32, term uint64, index uint32) {
	this.uncommittedEntries = append([]*RaftLogEntry{{value: value, term: term, index: index}}, this.uncommittedEntries...)
}

func (this *RaftLog) HasUnCommittedValue() bool {
	return this.uncommittedEntries != nil
}

func (this *RaftLog) HasIndex(index uint32) bool {
	return len(this.committedEntries) > 0 && uint32(len(this.committedEntries)-1) <= index
}

func (this *RaftLog) GetNextIndexToAppend() uint32 {
	return uint32(len(this.committedEntries))
}

func (this *RaftLog) GetTermByIndex(index uint32) uint64 {
	if len(this.committedEntries) == 0 {
		return 0
	} else {
		return this.committedEntries[len(this.committedEntries)-1].term
	}
}

func (this *RaftLog) GetLastCommitted() (_index uint32, _term uint64) {
	if len(this.committedEntries) == 0 {
		return 0, 0
	}

	index := uint32(len(this.committedEntries) - 1)
	term := this.committedEntries[len(this.committedEntries)-1].term

	return index, term
}

func (this *RaftLog) GetCommittedLogEntries() []uint32 {
	entriesValues := make([]uint32, len(this.committedEntries))

	for index, entry := range this.committedEntries {
		entriesValues[index] = entry.value
	}

	return entriesValues
}

func (this *RaftLog) GetLastCommittedIndex() uint32 {
	if len(this.committedEntries) == 0 {
		return 0
	} else {
		return uint32(len(this.committedEntries) - 1)
	}
}

func (this *RaftLog) trimCommittedEntriesLog(lastIndexInclusive uint32) {
	for i := int(lastIndexInclusive); i < len(this.committedEntries); i++ {
		this.committedEntries[i] = nil
	}
}

func CreateRaftLog() *RaftLog {
	return &RaftLog{
		committedEntries:   make([]*RaftLogEntry, 0),
		uncommittedEntries: make([]*RaftLogEntry, 0),
	}
}
