package scheduled_gossip

import (
	"distributed-systems/src/nodes"
)

type ScheduledGossipBroadcastData struct {
	pendingConfirmMessagesId map[uint64]*nodes.Message
}

func CreateScheduledGossipBroadcastData() *ScheduledGossipBroadcastData {
	return &ScheduledGossipBroadcastData{
		pendingConfirmMessagesId: make(map[uint64]*nodes.Message),
	}
}

func (data *ScheduledGossipBroadcastData) GetAll() []*nodes.Message {
	messages := make([]*nodes.Message, 0)

	for _, message := range messages {
		messages = append(messages, message)
	}

	return messages
}

func (data *ScheduledGossipBroadcastData) IsEmpty() bool {
	return len(data.pendingConfirmMessagesId) == 0
}

func (data *ScheduledGossipBroadcastData) AddAllIfNotContainedOrRemove(messages []*nodes.Message) {
	for _, message := range messages {
		data.AddIfNotContainedOrRemove(message)
	}
}

func (data *ScheduledGossipBroadcastData) AddIfNotContainedOrRemove(message *nodes.Message) {
	if _, contained := data.pendingConfirmMessagesId[message.GetMessageId()]; contained {
		delete(data.pendingConfirmMessagesId, message.GetMessageId())
	} else {
		data.pendingConfirmMessagesId[message.GetMessageId()] = message
	}
}
