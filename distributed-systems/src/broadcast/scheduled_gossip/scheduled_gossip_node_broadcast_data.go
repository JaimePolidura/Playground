package scheduled_gossip

import "distributed-systems/src/broadcast"

type ScheduledGossipBroadcastData struct {
	pendingConfirmMessagesId map[uint64]*broadcast.Message
}

func CreateScheduledGossipBroadcastData() *ScheduledGossipBroadcastData {
	return &ScheduledGossipBroadcastData{
		pendingConfirmMessagesId: make(map[uint64]*broadcast.Message),
	}
}

func (data *ScheduledGossipBroadcastData) GetAll() []*broadcast.Message {
	messages := make([]*broadcast.Message, 0)

	for _, message := range messages {
		messages = append(messages, message)
	}

	return messages
}

func (data *ScheduledGossipBroadcastData) IsEmpty() bool {
	return len(data.pendingConfirmMessagesId) == 0
}

func (data *ScheduledGossipBroadcastData) AddAllIfNotContainedOrRemove(messages []*broadcast.Message) {
	for _, message := range messages {
		data.AddIfNotContainedOrRemove(message)
	}
}

func (data *ScheduledGossipBroadcastData) AddIfNotContainedOrRemove(message *broadcast.Message) {
	if _, contained := data.pendingConfirmMessagesId[message.GetMessageId()]; contained {
		delete(data.pendingConfirmMessagesId, message.GetMessageId())
	} else {
		data.pendingConfirmMessagesId[message.GetMessageId()] = message
	}
}
