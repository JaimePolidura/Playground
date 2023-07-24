package fifo

import "distributed-systems/src/broadcast"

type FifoGossipNodeBroadcastData struct {
	lastSeqNumDelivered uint32

	buffer *BufferMessages
}

func CreateFifoGossipNodeBroadcastData() *FifoGossipNodeBroadcastData {
	return &FifoGossipNodeBroadcastData{
		buffer: CreateBufferMessages(),
	}
}

func (data *FifoGossipNodeBroadcastData) GetDeliverableMessages(seqNumbReceived uint32) []*broadcast.Message {
	if data.lastSeqNumDelivered+1 == seqNumbReceived {
		arr := data.buffer.GetAllDeliverable()

		if len(arr) > 0 {
			data.lastSeqNumDelivered = arr[len(arr)-1].SeqNum
		}

		return arr

	} else {
		return make([]*broadcast.Message, 0)
	}
}

func (data *FifoGossipNodeBroadcastData) AddToBuffer(message *broadcast.Message) {
	data.buffer.Add(message)
}

func (data *FifoGossipNodeBroadcastData) GetLastSeqNumDelivered() uint32 {
	return data.lastSeqNumDelivered
}
