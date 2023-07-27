package zab

import (
	"distributed-systems/src/broadcast/fifo"
	"distributed-systems/src/nodes"
	"sync/atomic"
)

func (this *ZabBroadcaster) handleAckRetransmissionMessageFollower(message *nodes.Message) {
	msgSeqNumbReceived := message.SeqNum
	lastSeqNumDelivered := this.getLastSeqNumDelivered(message.NodeIdOrigin)
	broadcastData := this.fifoBroadcastDataByNodeId[message.NodeIdOrigin]

	if msgSeqNumbReceived > lastSeqNumDelivered && message.NodeIdOrigin != this.selfNodeId {
		broadcastData.AddToBuffer(message)

		for _, messageInBuffer := range broadcastData.RetrieveDeliverableMessages(msgSeqNumbReceived) {
			this.onBroadcastMessageCallback(messageInBuffer)
		}
	}

	this.sendAckToNode(this.leaderNodeId, message)
}

func (this *ZabBroadcaster) onBroadcastMessageFollower(message *nodes.Message) {
	msgSeqNumbReceived := message.SeqNum
	lastSeqNumDelivered := this.getLastSeqNumDelivered(message.NodeIdOrigin)
	broadcastData := this.fifoBroadcastDataByNodeId[message.NodeIdOrigin]

	if msgSeqNumbReceived > lastSeqNumDelivered && message.NodeIdOrigin != this.selfNodeId {
		broadcastData.AddToBuffer(message)

		for _, messageInBuffer := range broadcastData.RetrieveDeliverableMessages(msgSeqNumbReceived) {
			this.onBroadcastMessageCallback(messageInBuffer)
		}
	}

	this.sendAckToNode(this.leaderNodeId, message)
}

func (this *ZabBroadcaster) sendBroadcastMessageToLeader(message *nodes.Message) {
	message.SeqNum = atomic.AddUint32(&this.seqNum, 1)
	this.messagesPendingLeaderAck.Add(this.leaderNodeId, message)
	this.nodeConnectionsStore.Get(this.leaderNodeId).Write(message)
}

func (this *ZabBroadcaster) getLastSeqNumDelivered(nodeId uint32) uint32 {
	if prevSeqNum, found := this.fifoBroadcastDataByNodeId[nodeId]; found {
		return prevSeqNum.GetLastSeqNumDelivered()
	} else {
		this.fifoBroadcastDataByNodeId[nodeId] = fifo.CreateFifoNodeBroadcastData()
		return 0
	}
}
