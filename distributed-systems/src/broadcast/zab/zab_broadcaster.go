package zab

import (
	"distributed-systems/src/broadcast"
	"distributed-systems/src/broadcast/fifo"
	"distributed-systems/src/broadcast/zab/ack"
	"distributed-systems/src/nodes"
	"distributed-systems/src/nodes/types"
	"fmt"
)

type ZabBroadcaster struct {
	selfNodeId             uint32
	leaderNodeId           uint32
	nodesConnectionManager *nodes.ConnectionManager

	//Leader
	seqNum           uint32
	seqNumToSendTurn uint32

	messagesDeliveredToFollowers *ack.MessagesAlreadyDelivered
	messagesPendingFollowersAck  *ack.MessagesPendingAck

	//Follower
	fifoBroadcastDataByNodeId  map[uint32]*fifo.FifoNodeBroadcastData
	onBroadcastMessageCallback func(newMessage *nodes.Message)
	messagesPendingLeaderAck   *ack.MessagesPendingAck
	largestSeqNumReceived      uint32

	onBroadcastMessage func(newMessage *nodes.Message)
}

func CreateZabBroadcaster(selfNodeId uint32, leaderNodeId uint32, retransmissionTimeout uint64) *ZabBroadcaster {
	broadcaster := &ZabBroadcaster{
		fifoBroadcastDataByNodeId:    map[uint32]*fifo.FifoNodeBroadcastData{},
		messagesDeliveredToFollowers: ack.CreateMessagesAlreadyDelivered(),
		messagesPendingFollowersAck:  ack.CreateMessagesPendingAck(retransmissionTimeout),
		messagesPendingLeaderAck:     ack.CreateMessagesPendingAck(retransmissionTimeout),
		selfNodeId:                   selfNodeId,
		leaderNodeId:                 leaderNodeId,
		seqNumToSendTurn:             1,
	}

	broadcaster.messagesPendingFollowersAck.SetOnRetransmissionCallback(broadcaster.doRetransmission)
	broadcaster.messagesPendingLeaderAck.SetOnRetransmissionCallback(broadcaster.doRetransmission)

	return broadcaster
}

func (this *ZabBroadcaster) SetOnBroadcastMessage(callback func(newMessage *nodes.Message)) *ZabBroadcaster {
	this.onBroadcastMessageCallback = callback
	return this
}

func (this *ZabBroadcaster) GetLargestSeqNumbReachievedLeader() uint32 {
	return this.largestSeqNumReceived
}

func (this *ZabBroadcaster) OnElectionStarted() {
	this.messagesPendingLeaderAck.StopRetransmissionTimer()
}

func (this *ZabBroadcaster) OnNewLeader(newLeaderNodeId uint32, newSeqNum uint32) {
	this.leaderNodeId = newLeaderNodeId
	this.seqNum = newSeqNum

	if this.isFollower() {
		this.messagesPendingLeaderAck.RestartRetransmissionTimer()
	}
	if this.isLeader() {
		this.seqNumToSendTurn = newSeqNum + 1
	}
}

func (this *ZabBroadcaster) doRetransmission(nodeIdToRetransmit uint32, message *nodes.Message) {
	fmt.Printf("[%d] ACK Timeout passed. Starting retransmission to node %d with SeqNum %d Message type %d\n",
		this.selfNodeId, nodeIdToRetransmit, message.SeqNum, message.Type)

	this.nodesConnectionManager.Send(nodeIdToRetransmit, message)
}

func (this *ZabBroadcaster) OnStop() {
	this.messagesPendingFollowersAck.StopRetransmissionTimer()
	this.messagesPendingLeaderAck.StopRetransmissionTimer()
}

func (this *ZabBroadcaster) Broadcast(message *nodes.Message) {
	if message.HasFlag(types.FLAG_BYPASS_LEADER) {
		this.sendMessageToFollowers(message)
		return
	}

	if this.isFollower() {
		this.sendBroadcastMessageToLeader(message)
	} else {
		this.sendMessageToFollowers(message)
	}
}

func (this *ZabBroadcaster) OnBroadcastMessage(message *nodes.Message) {
	if this.isFollower() {
		this.onBroadcastMessageFollower(message)
	}
}

func (this *ZabBroadcaster) SetNodeConnectionsManager(nodesConnectionManager *nodes.ConnectionManager) broadcast.Broadcaster {
	this.nodesConnectionManager = nodesConnectionManager
	return this
}

func (this *ZabBroadcaster) SetOnBroadcastMessageCallback(callback func(newMessage *nodes.Message)) broadcast.Broadcaster {
	this.onBroadcastMessageCallback = callback
	return this
}

func (this *ZabBroadcaster) HandleDoBroadcast(message *nodes.Message) {
	this.sendMessageToFollowers(message.WithType(types.MESSAGE_BROADCAST))
}

func (this *ZabBroadcaster) HandleAckMessage(message *nodes.Message) {
	seqNumAcked := message.GetContentToUint32()

	fmt.Printf("[%d] Received ACK message from node %d with SeqNum %d\n", this.selfNodeId, message.NodeIdSender, seqNumAcked)

	if this.isFollower() {
		this.messagesPendingLeaderAck.Delete(message.NodeIdSender, seqNumAcked)
	} else {
		this.messagesPendingFollowersAck.Delete(message.NodeIdSender, seqNumAcked)
		this.messagesDeliveredToFollowers.Add(message.NodeIdSender, seqNumAcked)
	}
}

func (this *ZabBroadcaster) addMessagePendingAck(pendingAck map[uint32]*fifo.FifoBufferMessages, nodeId uint32, message *nodes.Message) {
	if _, contained := pendingAck[nodeId]; !contained {
		pendingAck[nodeId] = fifo.CreateFifoBufferMessages()
	}

	pendingAck[nodeId].Add(message)
}

func (this *ZabBroadcaster) isFollower() bool {
	return this.leaderNodeId != this.selfNodeId
}

func (this *ZabBroadcaster) isLeader() bool {
	return this.leaderNodeId == this.selfNodeId
}

func (this *ZabBroadcaster) sendAckToNode(nodeIdToSendAck uint32, messageToAck *nodes.Message) {
	if this.selfNodeId != nodeIdToSendAck && messageToAck.HasNotFlag(types.FLAG_BYPASS_ORDERING) {
		ackMessage := nodes.CreateMessage(
			nodes.WithContentUInt32(messageToAck.SeqNum),
			nodes.WithOrigin(messageToAck.NodeIdOrigin),
			nodes.WithSenderNodeId(this.selfNodeId),
			nodes.WithType(types.MESSAGE_ACK))

		fmt.Printf("[%d] Sending ACK to node %d with SeqNum %d\n", this.selfNodeId, nodeIdToSendAck, messageToAck.SeqNum)

		this.nodesConnectionManager.Send(nodeIdToSendAck, ackMessage)
	}
}
