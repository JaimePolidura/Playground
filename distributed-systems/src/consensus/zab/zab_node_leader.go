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
				message := nodes.CreateMessageBroadcast(this.node.GetNodeId(), this.node.GetNodeId(), "").WithType(zab.MESSAGE_HEARTBEAT).WithFlag(nodes.FLAG_URGENT)
				this.node.Broadcast(message)
			}
		}
	}
}
