package ack

import (
	"distributed-systems/src/broadcast/fifo"
	"distributed-systems/src/nodes"
)

type MessagesPendingAck struct {
	messagesPendingAckByNodeId map[uint32]*fifo.FifoBufferMessages
}

func CreateMessagesPendingAck() *MessagesPendingAck {
	return &MessagesPendingAck{
		messagesPendingAckByNodeId: make(map[uint32]*fifo.FifoBufferMessages),
	}
}

func (this *MessagesPendingAck) GetAllLessThanSeqNum(nodeId uint32, seqNum uint32) []*nodes.Message {
	return this.getSeqNumsByNodeId(nodeId).GetMessagesLessThanSeqNum(seqNum)
}

func (this *MessagesPendingAck) Add(nodeId uint32, message *nodes.Message) {
	this.getSeqNumsByNodeId(nodeId).Add(message)
}

func (this *MessagesPendingAck) Delete(nodeId uint32, seqNum uint32) {
	this.getSeqNumsByNodeId(nodeId).RemoveBySeqNum(seqNum)
}

func (this *MessagesPendingAck) getSeqNumsByNodeId(nodeId uint32) *fifo.FifoBufferMessages {
	if value, contained := this.messagesPendingAckByNodeId[nodeId]; contained {
		return value
	} else {
		this.messagesPendingAckByNodeId[nodeId] = fifo.CreateFifoBufferMessages()
		return this.messagesPendingAckByNodeId[nodeId]
	}
}
