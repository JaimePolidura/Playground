package zab

import (
	"distributed-systems/src/broadcast"
	"distributed-systems/src/broadcast/zab"
	"distributed-systems/src/nodes"
)

import (
	"time"
)

type ZabNode struct {
	node *broadcast.Node

	state NodeState

	heartbeatTimeMs       uint64
	heartbeatTimeoutMs    uint64
	heartbeatSenderTicker *time.Ticker
	heartbeatTimerTimeout *time.Timer

	nodesIdRing                       []uint32
	selfNodeIdRingIndex               uint32
	nNodesThatHaveAckElectionProposal uint32

	epoch uint32

	leaderNodeId uint32
}

func CreateZabNode(selfNodeId uint32, port uint16, leaderNodeId uint32, heartbeatTimeMs uint64, heartbeatTimeoutMs uint64, broadcasterNode broadcast.Broadcaster) *ZabNode {
	zabNode := &ZabNode{
		node:                  broadcast.CreateNode(selfNodeId, port, broadcasterNode),
		heartbeatSenderTicker: time.NewTicker(time.Duration(heartbeatTimeMs * uint64(time.Millisecond))),
		heartbeatTimerTimeout: time.NewTimer(time.Duration(heartbeatTimeoutMs * uint64(time.Millisecond))),
		heartbeatTimeMs:       heartbeatTimeMs,
		leaderNodeId:          leaderNodeId,
		epoch:                 0,
		state:                 STARTING,
		nodesIdRing:           make([]uint32, 0),
	}

	if zabNode.IsLeader() {
		go zabNode.startSendingHeartbeats()
	}
	if zabNode.IsFollower() {
		go zabNode.startHeartbeatTimer()
	}

	zabNode.node.OnBroadcastMessage(zabNode.OnBroadcastMessage)
	zabNode.node.OnSingleMessage(zabNode.OnSingleMessage)

	return zabNode
}

func (this *ZabNode) OnBroadcastMessage(message *nodes.Message) {
	if message.IsType(zab.MESSAGE_HEARTBEAT) && this.IsFollower() {
		this.heartbeatTimerTimeout.Reset(this.GetDurationHeartbeatTimeout(this.heartbeatTimeMs))
		return
	}
}

func (this *ZabNode) OnSingleMessage(message *nodes.Message) {
	if message.IsType(zab.MESSAGE_ELECTION_FAILURE_DETECTED) {
		this.handleNodeFailureMessage()
	} else if message.IsType(zab.MESSAGE_ELECTION_PROPOSAL) {
		this.ackProposal(message)
	} else if message.IsType(zab.MESSAGE_ELECTION_ACK_PROPOSAL) {
		this.collectAckProposal()
	} else if message.IsType(zab.MESSAGE_ELECTION_COMMIT) {
		this.saveNewLeader(message.NodeIdSender)
	}
}

func (this *ZabNode) GetDurationHeartbeatTimeout(heartbeatSenderTicker uint64) time.Duration {
	return time.Duration(heartbeatSenderTicker*2 + heartbeatSenderTicker)
}

func (this *ZabNode) GetNodeId() uint32 {
	return this.node.GetNodeId()
}

func (this *ZabNode) GetNode() *broadcast.Node {
	return this.node
}

func (this *ZabNode) IsLeader() bool {
	return this.leaderNodeId == this.node.GetNodeId()
}

func (this *ZabNode) IsFollower() bool {
	return this.leaderNodeId != this.node.GetNodeId()
}

func (this *ZabNode) SetStateToBroadcast() {
	this.state = BROADCAST
}
