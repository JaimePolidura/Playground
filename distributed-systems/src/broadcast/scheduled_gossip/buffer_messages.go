package scheduled_gossip

import (
	"distributed-systems/src/broadcast"
	"sync"
)

type BufferMessages struct {
	messages []*broadcast.Message
	lock     sync.Mutex
}

func CreateBufferMessages() *BufferMessages {
	return &BufferMessages{
		messages: make([]*broadcast.Message, 0),
	}
}

func (buffer *BufferMessages) RetrieveAll() []*broadcast.Message {
	buffer.lock.Lock()
	messagesToReturn := buffer.messages
	buffer.messages = make([]*broadcast.Message, 0)
	buffer.lock.Unlock()

	return messagesToReturn
}

func (buffer *BufferMessages) Add(message *broadcast.Message) {
	buffer.lock.Lock()
	buffer.messages = append(buffer.messages, message)
	buffer.lock.Unlock()
}
