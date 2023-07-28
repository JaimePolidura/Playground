package zab

import (
	"distributed-systems/src/broadcast"
	"distributed-systems/src/broadcast/fifo"
	"distributed-systems/src/broadcast/zab/ack"
	"distributed-systems/src/nodes"
	"fmt"
)

const MESSAGE_ACK = 1
const MESSAGE_HEARTBEAT = 3
const MESSAGE_ELECTION_FAILURE_DETECTED = 4
const MESSAGE_ELECTION_PROPOSAL = 5
const MESSAGE_ELECTION_ACK_PROPOSAL = 6
const MESSAGE_ELECTION_COMMIT = 7
const MESSAGE_DO_BROADCAST = 8

type ZabBroadcaster struct {
	selfNodeId           uint32
	leaderNodeId         uint32
	nodeConnectionsStore *nodes.NodeConnectionsStore

	//Leader
	seqNum           uint32
	seqNumToSendTurn uint32

	messagesDeliveredToFollowers *ack.MessagesAlreadyDelivered
	messagesPendingFollowersAck  *ack.MessagesPendingAck

	//Follower
	fifoBroadcastDataByNodeId  map[uint32]*fifo.FifoNodeBroadcastData
	onBroadcastMessageCallback func(newMessage *nodes.Message)
	messagesPendingLeaderAck   *ack.MessagesPendingAck
}

func CreateZabBroadcaster(selfNodeId uint32, leaderNodeId uint32, retransmissionTimeout uint64, onBroadcastMessageCallback func(newMessage *nodes.Message)) *ZabBroadcaster {
	broadcaster := &ZabBroadcaster{
		fifoBroadcastDataByNodeId:    map[uint32]*fifo.FifoNodeBroadcastData{},
		messagesDeliveredToFollowers: ack.CreateMessagesAlreadyDelivered(),
		messagesPendingFollowersAck:  ack.CreateMessagesPendingAck(retransmissionTimeout),
		messagesPendingLeaderAck:     ack.CreateMessagesPendingAck(retransmissionTimeout),
		onBroadcastMessageCallback:   onBroadcastMessageCallback,
		selfNodeId:                   selfNodeId,
		leaderNodeId:                 leaderNodeId,
		seqNumToSendTurn:             1,
	}

	broadcaster.messagesPendingFollowersAck.SetOnRetransmissionCallback(broadcaster.doRetransmission)
	broadcaster.messagesPendingLeaderAck.SetOnRetransmissionCallback(broadcaster.doRetransmission)

	return broadcaster
}

func (this *ZabBroadcaster) doRetransmission(nodeIdToRetransmit uint32, message *nodes.Message) {
	fmt.Printf("[%d] ACK Timeout passed. Starting retransmission to node %d with SeqNum %d Message type %d\n",
		this.selfNodeId, nodeIdToRetransmit, message.SeqNum, message.Type)

	nodeConnection := this.nodeConnectionsStore.Get(nodeIdToRetransmit)
	nodeConnection.Write(message)
}

func (this *ZabBroadcaster) Broadcast(message *nodes.Message) {
	if message.HasFlag(nodes.FLAG_BYPASS_LEADER) {
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

func (this *ZabBroadcaster) SetOnBroadcastMessageCallback(callback func(newMessage *nodes.Message)) broadcast.Broadcaster {
	this.onBroadcastMessageCallback = callback
	return this
}

func (this *ZabBroadcaster) SetNodeConnectionsStore(store *nodes.NodeConnectionsStore) broadcast.Broadcaster {
	this.nodeConnectionsStore = store
	return this
}

func (this *ZabBroadcaster) HandleDoBroadcast(messages []*nodes.Message) {
	for _, message := range messages {
		this.sendMessageToFollowers(message.WithType(nodes.MESSAGE_BROADCAST))
	}
}

func (this *ZabBroadcaster) HandleAckMessage(messages []*nodes.Message) {
	message := messages[0]
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

func (this *ZabBroadcaster) sendAckToNode(nodeIdToSendAck uint32, messageToAck *nodes.Message) {
	if this.selfNodeId != nodeIdToSendAck && !messageToAck.HasFlag(nodes.FLAG_URGENT) {
		ackMessage := nodes.CreateMessage(
			nodes.WithOrigin(messageToAck.NodeIdOrigin),
			nodes.WithSenderNodeId(this.selfNodeId),
			nodes.WithType(MESSAGE_ACK))
		ackMessage.SetContentUin32(messageToAck.SeqNum)

		fmt.Printf("[%d] Sending ACK to node %d with SeqNum %d\n", this.selfNodeId, nodeIdToSendAck, messageToAck.SeqNum)

		this.nodeConnectionsStore.Get(nodeIdToSendAck).Write(ackMessage)
	}
}
