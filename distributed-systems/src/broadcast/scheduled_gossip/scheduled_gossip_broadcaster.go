package scheduled_gossip

import (
	"distributed-systems/src/broadcast"
	"sync"
	"sync/atomic"
	"time"
)

type ScheduledGossipBroadcaster struct {
	selfNodeId             uint32
	nodesToPickToBroadcast uint32
	seqNum                 uint32
	millisPeriodGossipTask uint64
	initialTTL             int32
	topology               []uint32

	broadcastDataByNodeId          map[uint32]*ScheduledGossipBroadcastData
	seqNumsDeliveredByOriginNodeId map[uint32]uint32
	newMessageCallback             func(newMessage *broadcast.Message)
	nodeConnectionsStore           *broadcast.NodeConnectionsStore
	buffer                         *BufferMessages
	gossipTaskLock                 sync.Mutex
}

func CreateScheduledGossipBroadcaster(millisPeriodGossipTask uint64, topology []uint32) *ScheduledGossipBroadcaster {
	scheduledGossipBroadcaster := &ScheduledGossipBroadcaster{
		broadcastDataByNodeId:          make(map[uint32]*ScheduledGossipBroadcastData),
		seqNumsDeliveredByOriginNodeId: make(map[uint32]uint32),
		millisPeriodGossipTask:         millisPeriodGossipTask,
		buffer:                         CreateBufferMessages(),
		topology:                       topology,
	}

	scheduledGossipBroadcaster.scheduleGossipTask()

	return scheduledGossipBroadcaster
}

func (this *ScheduledGossipBroadcaster) scheduleGossipTask() {
	ticker := time.NewTicker(time.Duration(this.millisPeriodGossipTask) * time.Millisecond)

	for {
		select {
		case <-ticker.C:
			this.startGossip()
		}
	}
}

func (this *ScheduledGossipBroadcaster) startGossip() {
	this.gossipTaskLock.Lock()

	for _, nodeConnection := range this.GetNodesConnectionPendingToSync() {
		messagesInBuffer := this.buffer.RetrieveAll()
		broadcastData := this.getBroadcastDataByNodeId(nodeConnection.GetNodeId())
		messagesPendingConfirm := broadcastData.GetAll()
		messagesToSend := this.removeDuplicatedMessages(append(messagesInBuffer, messagesPendingConfirm...))
		messagesToSend = this.decreaseTTL(messagesToSend)

		broadcastData.AddAllIfNotContainedOrRemove(messagesInBuffer)
		nodeConnection.WriteAll(messagesToSend)
	}

	this.gossipTaskLock.Unlock()
}

func (this *ScheduledGossipBroadcaster) removeDuplicatedMessages(duplicatedMessages []*broadcast.Message) []*broadcast.Message {
	messagesIdsSeen := make(map[uint64]uint64)
	notDuplicatedMessages := make([]*broadcast.Message, 0)

	for _, message := range duplicatedMessages {
		if _, contained := messagesIdsSeen[message.GetMessageId()]; !contained {
			notDuplicatedMessages = append(notDuplicatedMessages, message)
			messagesIdsSeen[message.GetMessageId()] = message.GetMessageId()
		}
	}

	return notDuplicatedMessages
}

func (this *ScheduledGossipBroadcaster) GetNodesConnectionPendingToSync() []*broadcast.NodeConnection {
	topology := this.getNodeConnectionsTopology()
	connectionsToSync := topology

	for nodeId, _ := range this.broadcastDataByNodeId {
		for _, nodeConnectionInTopology := range topology {
			if nodeConnectionInTopology.GetNodeId() != nodeId {
				connectionsToSync = append(connectionsToSync, this.nodeConnectionsStore.Get(nodeId))
			}
		}
	}

	return connectionsToSync
}

func (this *ScheduledGossipBroadcaster) decreaseTTL(messages []*broadcast.Message) []*broadcast.Message {
	for _, message := range messages {
		message.TTL = message.TTL - 1
	}

	return messages
}

func (this *ScheduledGossipBroadcaster) Broadcast(message *broadcast.Message) {
	this.doBroadcast([]*broadcast.Message{message}, true)
}

func (this *ScheduledGossipBroadcaster) doBroadcast(messages []*broadcast.Message, firstTime bool) {
	for _, message := range messages {
		if firstTime {
			atomic.AddUint32(&this.seqNum, 1)
		}

		if firstTime {
			message.SeqNum = this.seqNum
			message.TTL = this.initialTTL
		} else {
			message.TTL = message.TTL - 1
		}

		lastDeliveredSeqNumByOrigin := this.getLastSeqNumsDeliveredByOriginNodeId(message.NodeIdOrigin)

		if message.TTL != 0 && lastDeliveredSeqNumByOrigin < message.SeqNum {
			this.buffer.Add(message)
		}
		if !firstTime && lastDeliveredSeqNumByOrigin < message.SeqNum {
			this.seqNumsDeliveredByOriginNodeId[message.NodeIdOrigin] = lastDeliveredSeqNumByOrigin + 1
			this.newMessageCallback(message)
			this.broadcastDataByNodeId[message.NodeIdSender].AddIfNotContainedOrRemove(message)
		}
	}
}

func (this *ScheduledGossipBroadcaster) OnBroadcastMessage(messages []*broadcast.Message, newMessageCallback func(newMessage *broadcast.Message)) {
	this.newMessageCallback = newMessageCallback
	this.doBroadcast(messages, false)
}

func (this *ScheduledGossipBroadcaster) SetNodeConnectionsStore(store *broadcast.NodeConnectionsStore) broadcast.Broadcaster {
	this.nodeConnectionsStore = store
	return this
}

func (this *ScheduledGossipBroadcaster) getNodeConnectionsTopology() []*broadcast.NodeConnection {
	connectionsTopology := make([]*broadcast.NodeConnection, len(this.topology))

	for index, nodeId := range this.topology {
		connectionsTopology[index] = this.nodeConnectionsStore.Get(nodeId)
	}

	return connectionsTopology
}

func (this *ScheduledGossipBroadcaster) getBroadcastDataByNodeId(nodeId uint32) *ScheduledGossipBroadcastData {
	if data, contained := this.broadcastDataByNodeId[nodeId]; contained {
		return data
	} else {
		this.broadcastDataByNodeId[nodeId] = CreateScheduledGossipBroadcastData()
		return this.broadcastDataByNodeId[nodeId]
	}
}

func (this *ScheduledGossipBroadcaster) getLastSeqNumsDeliveredByOriginNodeId(originNodeId uint32) uint32 {
	if last, contained := this.seqNumsDeliveredByOriginNodeId[originNodeId]; contained {
		return last
	} else {
		this.seqNumsDeliveredByOriginNodeId[originNodeId] = 0
		return 0
	}
}
