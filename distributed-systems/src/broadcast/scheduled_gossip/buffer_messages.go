package scheduled_gossip

import (
	"distributed-systems/src/broadcast"
	"sync"
)

type BufferMessages struct {
	messages []*broadcast.BroadcastMessage
	lock     sync.Mutex
}

func CreateBufferMessages() *BufferMessages {
	return &BufferMessages{
		messages: make([]*broadcast.BroadcastMessage, 0),
	}
}

func (buffer *BufferMessages) RetrieveAll() []*broadcast.BroadcastMessage {
	buffer.lock.Lock()
	messagesToReturn := buffer.messages
	buffer.messages = make([]*broadcast.BroadcastMessage, 0)
	buffer.lock.Unlock()

	return messagesToReturn
}

func (buffer *BufferMessages) Add(message *broadcast.BroadcastMessage) {
	buffer.lock.Lock()
	buffer.messages = append(buffer.messages, message)
	buffer.lock.Unlock()
}
