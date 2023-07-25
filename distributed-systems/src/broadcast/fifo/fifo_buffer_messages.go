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

func (buffer *FifoBufferMessages) RetrieveAllDeliverable() []*broadcast.BroadcastMessage {
	toReturn := make([]*broadcast.BroadcastMessage, 0)

	if len(buffer.messages) == 0 {
		return toReturn
	}

	if len(buffer.messages) == 1 {
		toReturn = append(toReturn, buffer.messages[0])
		buffer.messages = make([]*broadcast.BroadcastMessage, 0)
		return toReturn
	}

	for index := 1; index < len(buffer.messages); index++ {
		actual := buffer.messages[index]
		prev := buffer.messages[index-1]

		if actual.SeqNum-prev.SeqNum == 1 {
			toReturn = append(toReturn, actual)
			buffer.messages = buffer.messages[1:]
		} else {
			break
		}
	}

	return toReturn
}

func (buffer *FifoBufferMessages) Add(message *broadcast.BroadcastMessage) {
	index := buffer.getIndexToAdd(message.SeqNum)
	buffer.messages = buffer.insert(index, message)
}

func (buffer *FifoBufferMessages) getIndexToAdd(seqNum uint32) int {
	indexToAdd := 0

	for i := 0; i < len(buffer.messages); i++ {
		hasNext := i+1 < len(buffer.messages)
		indexToAdd = i

		if buffer.messages[i].SeqNum < seqNum &&
			(!hasNext || buffer.messages[i+1].SeqNum > seqNum) {
			break
		}
	}

	return indexToAdd
}

func (buffer *FifoBufferMessages) insert(index int, value *broadcast.BroadcastMessage) []*broadcast.BroadcastMessage {
	if len(buffer.messages) == index {
		return append(buffer.messages, value)
	}

	buffer.messages = append(buffer.messages[:index+1], buffer.messages[index:]...)
	buffer.messages[index] = value

	return buffer.messages
}
