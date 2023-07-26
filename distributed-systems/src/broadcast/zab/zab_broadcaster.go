package zab

import (
	"distributed-systems/src/broadcast"
	"distributed-systems/src/broadcast/fifo"
	"strconv"
)

const ACK_FLAG = 1
const ACK_RETRANSMISSION_FLAG = 2

type ZabBroadcaster struct {
	selfNodeId           uint32
	leaderNodeId         uint32
	nodeConnectionsStore *broadcast.NodeConnectionsStore
	
	//Leader
	seqNum                         uint32
	seqNumToSendTurn               uint32
	pendingFollowerAck             map[uint32]*fifo.FifoBufferMessages
	receivedSeqNumByFollowerNodeId map[uint32]map[uint32]uint32

	//Follower
	fifoBroadcastDataByNodeId map[uint32]*fifo.FifoNodeBroadcastData
	pendingLeaderAck          map[uint32]*fifo.FifoBufferMessages
}

func CreateZabBroadcaster(selfNodeId uint32, leaderNodeId uint32) *ZabBroadcaster {
	return &ZabBroadcaster{
		selfNodeId:                     selfNodeId,
		leaderNodeId:                   leaderNodeId,
		fifoBroadcastDataByNodeId:      map[uint32]*fifo.FifoNodeBroadcastData{},
		pendingFollowerAck:             map[uint32]*fifo.FifoBufferMessages{},
		pendingLeaderAck:               map[uint32]*fifo.FifoBufferMessages{},
		receivedSeqNumByFollowerNodeId: map[uint32]map[uint32]uint32{},
		seqNumToSendTurn:               1,
	}
}

func (this *ZabBroadcaster) Broadcast(message *broadcast.BroadcastMessage) {
	if this.isFollower() {
		this.sendBroadcastMessageToLeader(message)
	} else {
		this.broadcastLeader(message)
	}
}

func (this *ZabBroadcaster) OnBroadcastMessage(messages []*broadcast.BroadcastMessage, newMessageCallback func(newMessage *broadcast.BroadcastMessage)) {
	message := messages[0]

	if this.isFollower() {
		this.onBroadcastMessageFollower(message, newMessageCallback)
	} else {
		this.onBroadcastMessageLeader(message)
	}
}

func (this *ZabBroadcaster) SetNodeConnectionsStore(store *broadcast.NodeConnectionsStore) broadcast.Broadcaster {
	this.nodeConnectionsStore = store
	return this
}

func (this *ZabBroadcaster) addMessagePendingAck(pendingAck map[uint32]*fifo.FifoBufferMessages, nodeId uint32, message *broadcast.BroadcastMessage) {
	if _, contained := pendingAck[nodeId]; !contained {
		pendingAck[nodeId] = fifo.CreateFifoBufferMessages()
	}

	pendingAck[nodeId].Add(message)
}

func (this *ZabBroadcaster) isFollower() bool {
	return this.leaderNodeId != this.selfNodeId
}

func (this *ZabBroadcaster) sendAckToNode(nodeIdToSendAck uint32, messageToAck *broadcast.BroadcastMessage) {
	ackMessage := broadcast.CreateMessageWithFlags(messageToAck.NodeIdOrigin, this.selfNodeId, strconv.Itoa(int(messageToAck.SeqNum)),
		ACK_FLAG)
	ackMessage.SetContentUin32(messageToAck.SeqNum)

	this.nodeConnectionsStore.Get(nodeIdToSendAck).Write(ackMessage)
}

func (this *ZabBroadcaster) removeMessagePendingAck(pendingAcksMap map[uint32]*fifo.FifoBufferMessages, ackMessage *broadcast.BroadcastMessage) {
	seqNumAcked := ackMessage.GetContentToUint32()

	fifoMessagesPendingAck := pendingAcksMap[ackMessage.NodeIdSender]
	messagesToResend := fifoMessagesPendingAck.GetMessagesLessThanSeqNum(seqNumAcked)
	fifoMessagesPendingAck.RemoveBySeqNum(seqNumAcked)

	for _, messageToResend := range messagesToResend {
		messageToResend = messageToResend.SetFlag(ACK_RETRANSMISSION_FLAG)
		this.nodeConnectionsStore.Get(ackMessage.NodeIdSender).Write(messageToResend)
	}
}
