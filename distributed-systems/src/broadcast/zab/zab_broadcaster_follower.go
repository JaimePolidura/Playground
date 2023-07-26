package zab

import (
	"distributed-systems/src/broadcast"
	"distributed-systems/src/broadcast/fifo"
	"sync/atomic"
)

func (this *ZabBroadcaster) onBroadcastMessageFollower(message *broadcast.BroadcastMessage, newMessageCallback func(newMessage *broadcast.BroadcastMessage)) {
	if !message.HasFlag(ACK_FLAG) {
		msgSeqNumbReceived := message.SeqNum
		lastSeqNumDelivered := this.getLastSeqNumDelivered(message.NodeIdOrigin)
		broadcastData := this.fifoBroadcastDataByNodeId[message.NodeIdOrigin]

		if msgSeqNumbReceived > lastSeqNumDelivered && message.NodeIdOrigin != this.selfNodeId {
			broadcastData.AddToBuffer(message)

			for _, messageInBuffer := range broadcastData.RetrieveDeliverableMessages(msgSeqNumbReceived) {
				newMessageCallback(messageInBuffer)
			}
		}

		this.sendAckToNode(this.leaderNodeId, message)
		if !message.HasFlag(ACK_RETRANSMISSION_FLAG) {
			this.addMessagePendingAck(this.pendingLeaderAck, this.leaderNodeId, message)
		}

	} else {
		this.removeMessagePendingAck(this.pendingLeaderAck, message)
	}
}

func (this *ZabBroadcaster) sendBroadcastMessageToLeader(message *broadcast.BroadcastMessage) {
	message.SeqNum = atomic.AddUint32(&this.seqNum, 1)
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
