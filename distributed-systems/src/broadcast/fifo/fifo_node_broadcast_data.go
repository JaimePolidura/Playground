package fifo

import "distributed-systems/src/broadcast"

type FifoNodeBroadcastData struct {
	lastSeqNumDelivered uint32

	buffer *FifoBufferMessages
}

func CreateFifoNodeBroadcastData() *FifoNodeBroadcastData {
	return &FifoNodeBroadcastData{
		buffer: CreateFifoBufferMessages(),
	}
}

func (data *FifoNodeBroadcastData) GetDeliverableMessages(seqNumbReceived uint32) []*broadcast.BroadcastMessage {
	if data.lastSeqNumDelivered+1 == seqNumbReceived {
		arr := data.buffer.RetrieveAllDeliverable()

		if len(arr) > 0 {
			data.lastSeqNumDelivered = arr[len(arr)-1].SeqNum
		}

		return arr

	} else {
		return make([]*broadcast.BroadcastMessage, 0)
	}
}

func (data *FifoNodeBroadcastData) AddToBuffer(message *broadcast.BroadcastMessage) {
	data.buffer.Add(message)
}

func (data *FifoNodeBroadcastData) GetLastSeqNumDelivered() uint32 {
	return data.lastSeqNumDelivered
}
