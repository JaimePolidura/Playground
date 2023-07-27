package zab

import (
	"distributed-systems/src/nodes"
	"fmt"
	"sync/atomic"
	"time"
)

func (this *ZabBroadcaster) broadcastToFollowers(message *nodes.Message) {
	seqNumForMessage := this.getSeqNumForMessage(message)
	messageToBroadcast := message.Clone()
	messageToBroadcast.SeqNum = seqNumForMessage
	messageToBroadcast.NodeIdSender = this.selfNodeId
	messageToBroadcast.NodeIdOrigin = message.NodeIdSender

	if message.HasFlag(nodes.FLAG_URGENT) {
		this.broadcastMessageToFollowers(message.NodeIdSender, messageToBroadcast)
	} else {
		this.executeNonUrgentMessage(seqNumForMessage, message.NodeIdSender, messageToBroadcast)
	}
}

func (this *ZabBroadcaster) executeNonUrgentMessage(seqNumForMessage uint32, nodeIdSender uint32, message *nodes.Message) {
	fmt.Printf("[%d] Broadcasting message from node %d old SeqNum %d with new SeqNum %d of message type %d\n", this.selfNodeId,
		message.NodeIdSender, message.SeqNum, message.SeqNum, message.Type)

	go func() {
		this.waitBroadcastLeaderTurn(seqNumForMessage)

		this.sendAckToNode(nodeIdSender, message)
		this.broadcastMessageToFollowers(nodeIdSender, message)

		atomic.AddUint32(&this.seqNumToSendTurn, 1)
	}()
}

func (this *ZabBroadcaster) getSeqNumForMessage(message *nodes.Message) uint32 {
	if message.HasNotFlag(nodes.FLAG_URGENT) {
		return atomic.AddUint32(&this.seqNum, 1)
	} else {
		return message.SeqNum
	}
}

func (this *ZabBroadcaster) waitBroadcastLeaderTurn(seqNumBroadcastToWait uint32) {
	for seqNumBroadcastToWait != atomic.LoadUint32(&this.seqNumToSendTurn) {
		time.Sleep(0) //Deschedule go routine
	}
}

func (this *ZabBroadcaster) broadcastMessageToFollowers(originalNodeSender uint32, message *nodes.Message) {
	message.Type = nodes.MESSAGE_BROADCAST

	this.forEachFollowerExcept(originalNodeSender, func(followerNodeConnection *nodes.NodeConnection) {
		if followerNodeConnection.GetNodeId() != this.leaderNodeId {
			followerNodeConnection.Write(message)

		} else if followerNodeConnection.GetNodeId() == this.leaderNodeId && message.HasNotFlag(nodes.FLAG_URGENT) {
			this.onBroadcastMessageCallback(message)
		}

		if message.HasNotFlag(nodes.FLAG_URGENT) {
			this.messagesPendingFollowersAck.Add(followerNodeConnection.GetNodeId(), message)
		}
	})
}

func (this *ZabBroadcaster) forEachFollowerExcept(exceptNodeId uint32, consumer func(connection *nodes.NodeConnection)) {
	for _, followerNodeConnection := range this.nodeConnectionsStore.ToArrayNodeConnections() {
		if followerNodeConnection.GetNodeId() != exceptNodeId {
			consumer(followerNodeConnection)
		}
	}
}
