package zab

import (
	"distributed-systems/src/nodes"
	"distributed-systems/src/nodes/types"
)

func (this *ZabNode) startSendingHeartbeats() {
	defer this.heartbeatSenderTicker.Stop()

	for {
		select {
		case <-this.heartbeatSenderTicker.C:
			message := nodes.CreateMessage(
				nodes.WithNodeId(this.GetNodeId()),
				nodes.WithType(types.MESSAGE_HEARTBEAT),
				nodes.WithFlags(types.FLAG_BYPASS_ORDERING, types.FLAG_BYPASS_LEADER))

			this.node.Broadcast(message)
		}
	}
}
