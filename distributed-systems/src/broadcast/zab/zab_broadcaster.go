package zab

import (
	"distributed-systems/src/broadcast"
	"distributed-systems/src/broadcast/fifo"
	"sync"
	"sync/atomic"
	"time"
)

type ZabBroadcaster struct {
	selfNodeId           uint32
	leaderNodeId         uint32
	nodeConnectionsStore *broadcast.NodeConnectionsStore

	//Leader
	seqNum           uint32
	seqNumToSendTurn uint32

	//Follower
	broadcastDataByNodeId map[uint32]*fifo.FifoNodeBroadcastData
}

func CreateZabBroadcaster(selfNodeId uint32, leaderNodeId uint32) *ZabBroadcaster {
	return &ZabBroadcaster{
		selfNodeId:            selfNodeId,
		leaderNodeId:          leaderNodeId,
		broadcastDataByNodeId: map[uint32]*fifo.FifoNodeBroadcastData{},
		seqNumToSendTurn:      1,
	}
}

func (this *ZabBroadcaster) Broadcast(message *broadcast.BroadcastMessage) {
	if this.isFollower() {
		this.sendBroadcastMessageToLeader(message)
	} else {
		seqNumForMessage := atomic.AddUint32(&this.seqNum, 1)
		message.SeqNum = seqNumForMessage
		message.NodeIdOrigin = this.selfNodeId

		go func() {
			this.waitBroadcastLeaderTurn(seqNumForMessage)

			this.broadcastMessageToFollowers(message)

			atomic.AddUint32(&this.seqNumToSendTurn, 1)
		}()
	}
}

func (this *ZabBroadcaster) sendBroadcastMessageToLeader(message *broadcast.BroadcastMessage) {
	message.SeqNum = atomic.AddUint32(&this.seqNum, 1)
	this.nodeConnectionsStore.Get(this.leaderNodeId).Write(message)
}

func (this *ZabBroadcaster) waitBroadcastLeaderTurn(seqNumBroadcastToWait uint32) {
	for seqNumBroadcastToWait != atomic.LoadUint32(&this.seqNumToSendTurn) {
		time.Sleep(0) //Deschedule go routine
	}
}

func (this *ZabBroadcaster) broadcastMessageToFollowers(message *broadcast.BroadcastMessage) {
	this.forEachFollower(func(followerNodeConnection *broadcast.NodeConnection) {
		followerNodeConnection.Write(message)
	})
}

func (this *ZabBroadcaster) OnBroadcastMessage(messages []*broadcast.BroadcastMessage, newMessageCallback func(newMessage *broadcast.BroadcastMessage)) {
	message := messages[0]

	if this.isFollower() {
		msgSeqNumbReceived := message.SeqNum
		lastSeqNumDelivered := this.getLastSeqNumDelivered(message.NodeIdOrigin)
		broadcastData := this.broadcastDataByNodeId[message.NodeIdOrigin]

		if msgSeqNumbReceived > lastSeqNumDelivered && message.NodeIdOrigin != this.selfNodeId {
			broadcastData.AddToBuffer(message)

			for _, messageInBuffer := range broadcastData.GetDeliverableMessages(msgSeqNumbReceived) {
				newMessageCallback(messageInBuffer)
			}
		}
	} else {
		this.Broadcast(message)
	}
}

func (this *ZabBroadcaster) SetNodeConnectionsStore(store *broadcast.NodeConnectionsStore) broadcast.Broadcaster {
	this.nodeConnectionsStore = store
	return this
}

func (this *ZabBroadcaster) getLastSeqNumDelivered(nodeId uint32) uint32 {
	if prevSeqNum, found := this.broadcastDataByNodeId[nodeId]; found {
		return prevSeqNum.GetLastSeqNumDelivered()
	} else {
		this.broadcastDataByNodeId[nodeId] = fifo.CreateFifoNodeBroadcastData()
		return 0
	}
}

func (this *ZabBroadcaster) isFollower() bool {
	return this.leaderNodeId != this.selfNodeId
}

func (this *ZabBroadcaster) forEachFollower(consumer func(connection *broadcast.NodeConnection)) {
	for _, followerNodeConnection := range this.nodeConnectionsStore.ToArrayNodeConnections() {
		if followerNodeConnection.GetNodeId() != this.leaderNodeId {
			consumer(followerNodeConnection)
		}
	}
}

func (this *ZabBroadcaster) createWaitGroupFollowers() *sync.WaitGroup {
	wait := &sync.WaitGroup{}
	size := int(this.nodeConnectionsStore.Size())
	wait.Add(size)
	return wait
}