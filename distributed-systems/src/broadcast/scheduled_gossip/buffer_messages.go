package scheduled_gossip

import (
	"distributed-systems/src/nodes"
	"sync"
)

type BufferMessages struct {
	messages []*nodes.Message
	lock     sync.Mutex
}

func CreateBufferMessages() *BufferMessages {
	return &BufferMessages{
		messages: make([]*nodes.Message, 0),
	}
}

func (buffer *BufferMessages) RetrieveAll() []*nodes.Message {
	buffer.lock.Lock()
	messagesToReturn := buffer.messages
	buffer.messages = make([]*nodes.Message, 0)
	buffer.lock.Unlock()

	return messagesToReturn
}

func (buffer *BufferMessages) Add(message *nodes.Message) {
	buffer.lock.Lock()
	buffer.messages = append(buffer.messages, message)
	buffer.lock.Unlock()
}
