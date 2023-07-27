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
			if this.IsLeader() && this.state == BROADCAST {
				message := nodes.CreateMessageWithType(this.node.GetNodeId(), this.node.GetNodeId(), "", zab.MESSAGE_HEARTBEAT).AddFlag(nodes.BROADCAST_FLAG)
				this.node.Broadcast(message)
			}
		}
	}
}
