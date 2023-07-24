package scheduled_gossip

import (
	"distributed-systems/src/broadcast"
	"sync"
	"time"
)

type ScheduledGossipBroadcaster struct {
	buffer *BufferMessages

	millisPeriodGossipTask uint64

	gossipTaskLock sync.Mutex

	nodeConnectionsStore *broadcast.NodeConnectionsStore
}

func CreateScheduledGossipBroadcaster(millisPeriodGossipTask uint64) *ScheduledGossipBroadcaster {
	scheduledGossipBroadcaster := &ScheduledGossipBroadcaster{
		millisPeriodGossipTask: millisPeriodGossipTask,
		buffer:                 CreateBufferMessages(),
	}

	scheduledGossipBroadcaster.scheduleGossipTask()

	return scheduledGossipBroadcaster
}

func (broadcaster *ScheduledGossipBroadcaster) scheduleGossipTask() {
	ticker := time.NewTicker(time.Duration(broadcaster.millisPeriodGossipTask) * time.Millisecond)

	for {
		select {
		case <-ticker.C:
			broadcaster.startGossip()
		}
	}
}

func (broadcaster *ScheduledGossipBroadcaster) startGossip() {
	broadcaster.gossipTaskLock.Lock()

	//messages := broadcaster.buffer.GetAll()

	broadcaster.gossipTaskLock.Unlock()
}

func (broadcaster *ScheduledGossipBroadcaster) Broadcast(message *broadcast.Message) {
	broadcaster.buffer.Add(message)
}

func (broadcaster *ScheduledGossipBroadcaster) OnBroadcastMessage(message []*broadcast.Message, newMessageCallback func(newMessage *broadcast.Message)) {
}

func (broadcaster *ScheduledGossipBroadcaster) SetNodeConnectionsStore(store *broadcast.NodeConnectionsStore) broadcast.Broadcaster {
	broadcaster.nodeConnectionsStore = store
	return broadcaster
}
