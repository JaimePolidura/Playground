package scheduled_gossip

import "distributed-systems/src/broadcast"

type ScheduledGossipBroadcastData struct {
	pendingConfirmMessagesId map[uint64]*broadcast.BroadcastMessage
}

func CreateScheduledGossipBroadcastData() *ScheduledGossipBroadcastData {
	return &ScheduledGossipBroadcastData{
		pendingConfirmMessagesId: make(map[uint64]*broadcast.BroadcastMessage),
	}
}

func (data *ScheduledGossipBroadcastData) GetAll() []*broadcast.BroadcastMessage {
	messages := make([]*broadcast.BroadcastMessage, 0)

	for _, message := range messages {
		messages = append(messages, message)
	}

	return messages
}

func (data *ScheduledGossipBroadcastData) IsEmpty() bool {
	return len(data.pendingConfirmMessagesId) == 0
}

func (data *ScheduledGossipBroadcastData) AddAllIfNotContainedOrRemove(messages []*broadcast.BroadcastMessage) {
	for _, message := range messages {
		data.AddIfNotContainedOrRemove(message)
	}
}

func (data *ScheduledGossipBroadcastData) AddIfNotContainedOrRemove(message *broadcast.BroadcastMessage) {
	if _, contained := data.pendingConfirmMessagesId[message.GetMessageId()]; contained {
		delete(data.pendingConfirmMessagesId, message.GetMessageId())
	} else {
		data.pendingConfirmMessagesId[message.GetMessageId()] = message
	}
}
