package zab

import (
	"distributed-systems/src/broadcast/zab"
	"distributed-systems/src/nodes"
)

func (this *ZabNode) startSendingHeartbeats() {
	defer this.heartbeatSenderTicker.Stop()

	for {
		select {
		case <-this.heartbeatSenderTicker.C:
			message := nodes.CreateMessage(
				nodes.WithNodeId(this.GetNodeId()),
				nodes.WithType(zab.MESSAGE_HEARTBEAT),
				nodes.WithFlags(nodes.FLAG_BYPASS_ORDERING, nodes.FLAG_BYPASS_LEADER))

			this.node.Broadcast(message)
		}
	}
}
