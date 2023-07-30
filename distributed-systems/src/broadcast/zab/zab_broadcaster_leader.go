package zab

import (
	"distributed-systems/src/nodes"
	"distributed-systems/src/nodes/types"
	"fmt"
	"sync/atomic"
	"time"
)

func (this *ZabBroadcaster) sendMessageToFollowers(message *nodes.Message) {
	seqNumForMessage := this.getSeqNumForMessage(message)
	messageToBroadcast := message.Clone()
	messageToBroadcast.SeqNum = seqNumForMessage
	messageToBroadcast.NodeIdSender = this.selfNodeId
	messageToBroadcast.NodeIdOrigin = message.NodeIdSender

	if message.HasFlag(types.FLAG_BYPASS_ORDERING) {
		this.sendNonOrderedMessageToFollowers(message.NodeIdSender, messageToBroadcast)
	} else {
		this.orderAndSendMessageToFollowers(seqNumForMessage, message.NodeIdSender, messageToBroadcast)
	}
}

func (this *ZabBroadcaster) sendNonOrderedMessageToFollowers(originalNodeSender uint32, message *nodes.Message) {
	this.nodesConnectionManager.ForEachConnectionExcept(originalNodeSender, func(followerNodeConnection *nodes.NodeConnection) {
		if this.messagesDeliveredToFollowers.IsAlreadyDelivered(followerNodeConnection.GetNodeId(), message.SeqNum) {
			return
		}

		if followerNodeConnection.GetNodeId() != this.leaderNodeId {
			followerNodeConnection.Write(message)
		}
	})
}

func (this *ZabBroadcaster) orderAndSendMessageToFollowers(seqNumForMessage uint32, nodeIdSender uint32, message *nodes.Message) {
	fmt.Printf("[%d] Broadcasting message from node %d old SeqNum %d with new SeqNum %d of message type %d\n", this.selfNodeId,
		message.NodeIdSender, message.SeqNum, message.SeqNum, message.Type)

	this.waitBroadcastLeaderTurn(seqNumForMessage)
	this.sendAckToNode(nodeIdSender, message)
	atomic.AddUint32(&this.seqNumToSendTurn, 1)

	this.nodesConnectionManager.ForEachConnectionExcept(nodeIdSender, func(followerNodeConnection *nodes.NodeConnection) {
		if this.messagesDeliveredToFollowers.IsAlreadyDelivered(followerNodeConnection.GetNodeId(), message.SeqNum) {
			return
		}

		this.messagesPendingFollowersAck.Add(followerNodeConnection.GetNodeId(), message)

		if followerNodeConnection.GetNodeId() != this.leaderNodeId {
			followerNodeConnection.Write(message)
		} else if followerNodeConnection.GetNodeId() == this.leaderNodeId {
			this.onBroadcastMessageCallback(message)
		}
	})
}

func (this *ZabBroadcaster) getSeqNumForMessage(message *nodes.Message) uint32 {
	if message.HasNotFlag(types.FLAG_BYPASS_ORDERING) {
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
