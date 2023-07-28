package fifo

import (
	"distributed-systems/src/nodes"
)

type FifoNodeBroadcastData struct {
	lastSeqNumDelivered uint32

	buffer *FifoBufferMessages
}

func CreateFifoNodeBroadcastData() *FifoNodeBroadcastData {
	return &FifoNodeBroadcastData{
		buffer: CreateFifoBufferMessages(),
	}
}

func (this *FifoNodeBroadcastData) RetrieveDeliverableMessages(seqNumbReceived uint32) []*nodes.Message {
	if this.lastSeqNumDelivered+1 == seqNumbReceived {
		arr := this.buffer.RetrieveAllDeliverable()

		if len(arr) > 0 {
			this.lastSeqNumDelivered = arr[len(arr)-1].SeqNum
		}

		return arr

	} else {
		return make([]*nodes.Message, 0)
	}
}

func (this *FifoNodeBroadcastData) AddToBuffer(message *nodes.Message) {
	this.buffer.Add(message)
}

func (this *FifoNodeBroadcastData) GetLastSeqNumDelivered() uint32 {
	return this.lastSeqNumDelivered
}
