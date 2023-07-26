package zab

import (
	"distributed-systems/src/nodes"
	"sync/atomic"
	"time"
)

func (this *ZabBroadcaster) onBroadcastMessageLeader(message *nodes.Message) {
	if message.IsType(MESSAGE_ACK) {
		this.removeMessagePendingAck(this.pendingFollowerAck, message)
	} else {
		this.broadcastLeader(message)
	}
}

func (this *ZabBroadcaster) broadcastLeader(message *nodes.Message) {
	if this.isAlreadyReceivedSeqNumByFollowerNodeId(message) {
		this.sendAckToNode(message.NodeIdSender, message)
		return
	}

	seqNumForMessage := atomic.AddUint32(&this.seqNum, 1)
	messageToBroadcast := nodes.CreateMessageBroadcast(this.selfNodeId, this.selfNodeId, string(message.Content))
	messageToBroadcast.SeqNum = seqNumForMessage
	messageToBroadcast.NodeIdOrigin = this.selfNodeId

	go func() {
		this.waitBroadcastLeaderTurn(seqNumForMessage)
		this.broadcastMessageToFollowers(messageToBroadcast)
		atomic.AddUint32(&this.seqNumToSendTurn, 1)
	}()

	this.addReceivedSeqNumByFollowerNodeId(message)

	this.sendAckToNode(message.NodeIdSender, message)
}

func (this *ZabBroadcaster) waitBroadcastLeaderTurn(seqNumBroadcastToWait uint32) {
	for seqNumBroadcastToWait != atomic.LoadUint32(&this.seqNumToSendTurn) {
		time.Sleep(0) //Deschedule go routine
	}
}

func (this *ZabBroadcaster) broadcastMessageToFollowers(message *nodes.Message) {
	this.forEachFollower(func(followerNodeConnection *nodes.NodeConnection) {
		followerNodeConnection.Write(message)

		if !message.IsType(MESSAGE_ACK_RETRANSMISSION) {
			this.addMessagePendingAck(this.pendingFollowerAck, followerNodeConnection.GetNodeId(), message)
		}
	})
}

func (this *ZabBroadcaster) forEachFollower(consumer func(connection *nodes.NodeConnection)) {
	for _, followerNodeConnection := range this.nodeConnectionsStore.ToArrayNodeConnections() {
		if followerNodeConnection.GetNodeId() != this.leaderNodeId {
			consumer(followerNodeConnection)
		}
	}
}

func (this *ZabBroadcaster) addReceivedSeqNumByFollowerNodeId(message *nodes.Message) {
	if _, contains := this.receivedSeqNumByFollowerNodeId[message.NodeIdSender]; !contains {
		this.receivedSeqNumByFollowerNodeId[message.NodeIdSender] = map[uint32]uint32{}
	}

	this.receivedSeqNumByFollowerNodeId[message.NodeIdSender][message.SeqNum] = message.SeqNum
}

func (this *ZabBroadcaster) isAlreadyReceivedSeqNumByFollowerNodeId(message *nodes.Message) bool {
	if _, contains := this.receivedSeqNumByFollowerNodeId[message.NodeIdSender]; !contains {
		this.receivedSeqNumByFollowerNodeId[message.NodeIdSender] = map[uint32]uint32{}
		return false
	}

	_, contained := this.receivedSeqNumByFollowerNodeId[message.NodeIdSender][message.SeqNum]

	return contained
}
