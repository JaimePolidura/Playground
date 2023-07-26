package zab

import (
	"distributed-systems/src/nodes"
)

func (this *ZabNode) startSendingHeartbeats() {
	defer this.heartbeatSenderTicker.Stop()

	for {
		select {
		case <-this.heartbeatSenderTicker.C:
			if this.node.GetNodeId() == this.leaderNodeId {
				message := nodes.CreateMessageWithType(this.node.GetNodeId(), this.node.GetNodeId(), "", MESSAGE_HEARTBEAT).AddFlag(nodes.BROADCAST_FLAG)
				this.node.Broadcast(message)
			}
		}
	}
}
