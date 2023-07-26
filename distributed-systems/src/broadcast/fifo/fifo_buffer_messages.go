package fifo

import (
	"distributed-systems/src/broadcast"
)

type FifoBufferMessages struct {
	messages []*broadcast.BroadcastMessage //Natural order
}

func CreateFifoBufferMessages() *FifoBufferMessages {
	return &FifoBufferMessages{
		messages: make([]*broadcast.BroadcastMessage, 0),
	}
}

func (this *FifoBufferMessages) GetMessagesLessThanSeqNum(seqNum uint32) []*broadcast.BroadcastMessage {
	toReturn := make([]*broadcast.BroadcastMessage, 0)

	for _, message := range this.messages {
		if message.SeqNum < seqNum {
			toReturn = append(toReturn, message)
		}
	}

	return toReturn
}

func (this *FifoBufferMessages) RemoveBySeqNum(seqNum uint32) {
	index := this.getIndexBySeqNum(seqNum)

	if index != -1 {
		this.messages = append(this.messages[:index], this.messages[index+1:]...)
	}
}

func (this *FifoBufferMessages) RetrieveAllDeliverable() []*broadcast.BroadcastMessage {
	toReturn := make([]*broadcast.BroadcastMessage, 0)

	if len(this.messages) == 0 {
		return toReturn
	}

	if len(this.messages) == 1 {
		toReturn = append(toReturn, this.messages[0])
		this.messages = make([]*broadcast.BroadcastMessage, 0)
		return toReturn
	}

	for index := 1; index < len(this.messages); index++ {
		actual := this.messages[index]
		prev := this.messages[index-1]

		if actual.SeqNum-prev.SeqNum == 1 {
			toReturn = append(toReturn, actual)
			this.messages = this.messages[1:]
		} else {
			break
		}
	}

	return toReturn
}

func (this *FifoBufferMessages) Add(message *broadcast.BroadcastMessage) {
	index := this.getIndexToAdd(message.SeqNum)
	this.messages = this.insert(index, message)
}

func (this *FifoBufferMessages) getIndexBySeqNum(seqNum uint32) int {
	for i := 0; i < len(this.messages); i++ {
		if this.messages[i].SeqNum == seqNum {
			return i
		}
	}

	return -1
}

func (this *FifoBufferMessages) getIndexToAdd(seqNum uint32) int {
	indexToAdd := 0

	for i := 0; i < len(this.messages); i++ {
		hasNext := i+1 < len(this.messages)
		indexToAdd = i

		if this.messages[i].SeqNum < seqNum &&
			(!hasNext || this.messages[i+1].SeqNum > seqNum) {
			break
		}
	}

	return indexToAdd
}

func (this *FifoBufferMessages) insert(index int, value *broadcast.BroadcastMessage) []*broadcast.BroadcastMessage {
	if len(this.messages) == index {
		return append(this.messages, value)
	}

	this.messages = append(this.messages[:index+1], this.messages[index:]...)
	this.messages[index] = value

	return this.messages
}
