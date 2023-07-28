package ack

import (
	"distributed-systems/src/broadcast/fifo"
	"distributed-systems/src/nodes"
	"time"
)

type MessagesPendingAck struct {
	messagesPendingAckByNodeId map[uint32]*MessagesPendingAckEntry
	timeoutRetransmissionMs    time.Duration

	retransmissionCallback func(nodeId uint32, message *nodes.Message)
}

type MessagesPendingAckEntry struct {
	nodeId                       uint32
	messages                     *fifo.FifoBufferMessages
	retransmissionTimersBySeqNum map[uint32]*time.Timer
}

func CreateMessagesPendingAckEntry(nodeId uint32) *MessagesPendingAckEntry {
	return &MessagesPendingAckEntry{
		messages:                     fifo.CreateFifoBufferMessages(),
		retransmissionTimersBySeqNum: make(map[uint32]*time.Timer),
	}
}

func CreateMessagesPendingAck(timeoutRetransmissionMs uint64) *MessagesPendingAck {
	return &MessagesPendingAck{
		messagesPendingAckByNodeId: make(map[uint32]*MessagesPendingAckEntry),
		timeoutRetransmissionMs:    time.Duration(timeoutRetransmissionMs * uint64(time.Millisecond)),
	}
}

func (this *MessagesPendingAck) SetOnRetransmissionCallback(retransmissionCallback func(nodeId uint32, message *nodes.Message)) {
	this.retransmissionCallback = retransmissionCallback
}

func (this *MessagesPendingAck) GetAllLessThanSeqNum(nodeId uint32, seqNum uint32) []*nodes.Message {
	return this.getSeqNumsByNodeId(nodeId).messages.GetMessagesLessThanSeqNum(seqNum)
}

func (this *MessagesPendingAck) Add(nodeId uint32, message *nodes.Message) {
	pendingMessagesAckEntry := this.getSeqNumsByNodeId(nodeId)

	pendingMessagesAckEntry.messages.Add(message)

	if _, contained := pendingMessagesAckEntry.retransmissionTimersBySeqNum[message.SeqNum]; !contained {
		newTimer := time.NewTimer(this.timeoutRetransmissionMs)
		pendingMessagesAckEntry.retransmissionTimersBySeqNum[message.SeqNum] = newTimer
		go func() { this.onRetransmission(newTimer, nodeId, message) }()
	}
}

func (this *MessagesPendingAck) onRetransmission(timer *time.Timer, nodeId uint32, message *nodes.Message) {
	<-timer.C
	this.retransmissionCallback(nodeId, message)
	timer.Reset(this.timeoutRetransmissionMs)
}

func (this *MessagesPendingAck) Delete(nodeId uint32, seqNum uint32) {
	pendingMessagesAckEntry := this.getSeqNumsByNodeId(nodeId)

	pendingMessagesAckEntry.messages.RemoveBySeqNum(seqNum)

	if timer, contained := pendingMessagesAckEntry.retransmissionTimersBySeqNum[seqNum]; contained {
		delete(pendingMessagesAckEntry.retransmissionTimersBySeqNum, seqNum)
		timer.Stop()
	}
}

func (this *MessagesPendingAck) getSeqNumsByNodeId(nodeId uint32) *MessagesPendingAckEntry {
	if value, contained := this.messagesPendingAckByNodeId[nodeId]; contained {
		return value
	} else {
		this.messagesPendingAckByNodeId[nodeId] = CreateMessagesPendingAckEntry(nodeId)
		return this.messagesPendingAckByNodeId[nodeId]
	}
}
