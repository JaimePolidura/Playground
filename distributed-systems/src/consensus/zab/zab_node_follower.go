package zab

import (
	"distributed-systems/src/broadcast/zab"
	"distributed-systems/src/nodes"
	"distributed-systems/src/utils"
	"fmt"
	"time"
)

func (this *ZabNode) handleHeartbeatMessage(message []*nodes.Message) {
	this.heartbeatTimerTimeout.Reset(time.Duration(this.heartbeatTimeoutMs * uint64(time.Millisecond)))
}

func (this *ZabNode) startHeartbeatTimer() {
	for {
		select {
		case <-this.heartbeatTimerTimeout.C:
			if this.IsFollower() && this.state == BROADCAST {
				if !this.electionLock.TryLock() {
					return
				}

				fmt.Printf("[%d] Detected leader %d failure. Broadcasting MESSAGE_ELECTION_FAILURE_DETECTED\n",
					this.GetNodeId(), this.leaderNodeId)

				this.changeStateFromBroadcastToElection()

				this.node.Broadcast(nodes.CreateMessage(
					nodes.WithNodeId(this.GetNodeId()),
					nodes.WithType(zab.MESSAGE_ELECTION_FAILURE_DETECTED),
					nodes.WithFlags(nodes.FLAG_BYPASS_LEADER, nodes.FLAG_BYPASS_ORDERING)))

				if this.leaderNodeId == this.prevNodeRing {
					this.proposeMySelfAsTheNewLeader()
				}
			}
		}
	}
}

func (this *ZabNode) handleNodeFailureMessage(message []*nodes.Message) {
	if !this.electionLock.TryLock() {
		return
	}

	this.changeStateFromBroadcastToElection()

	failureNodeIsPrev := this.prevNodeRing == this.leaderNodeId

	fmt.Printf("[%d] Receieved MESSAGE_ELECTION_FAILURE_DETECTED from node %d. Is prev node? %t\n",
		this.GetNodeId(), message[0].NodeIdSender, failureNodeIsPrev)

	if failureNodeIsPrev {
		this.proposeMySelfAsTheNewLeader()
	}
}

func (this *ZabNode) proposeMySelfAsTheNewLeader() {
	fmt.Printf("[%d] Sending proposal to be leader to followers\n", this.GetNodeId())

	this.node.Broadcast(nodes.CreateMessage(
		nodes.WithNodeId(this.GetNodeId()),
		nodes.WithType(zab.MESSAGE_ELECTION_PROPOSAL),
		nodes.WithFlags(nodes.FLAG_BYPASS_LEADER, nodes.FLAG_BYPASS_ORDERING)))
}

func (this *ZabNode) handleElectionProposalMessage(messages []*nodes.Message) {
	proposerNodeId := messages[0].NodeIdSender
	largestSeqNum := this.node.GetBroadcaster().(*zab.ZabBroadcaster).GetLargestSeqNumbReachievedLeader()

	fmt.Printf("[%d] Accepting proposal from node %d to be the new leader\n", this.GetNodeId(), proposerNodeId)

	this.GetNode().GetConnectionManager().Send(proposerNodeId, nodes.CreateMessage(
		nodes.WithContentUInt32(largestSeqNum),
		nodes.WithNodeId(this.GetNodeId()),
		nodes.WithType(zab.MESSAGE_ELECTION_ACK_PROPOSAL),
		nodes.WithFlags(nodes.FLAG_BYPASS_LEADER, nodes.FLAG_BYPASS_ORDERING)))
}

func (this *ZabNode) handleElectionAckProposalMessage(message []*nodes.Message) {
	this.nNodesThatHaveAckElectionProposal++
	largestSeqSumSeenByFollower := message[0].GetContentToUint32()
	nNodesQuorum := this.GetConnectionManager().GetNumberConnections()/2 + 1
	isQuorumSatisfied := this.nNodesThatHaveAckElectionProposal >= nNodesQuorum
	this.largestSeqNumSeenFromFollowers = utils.MaxUint32(this.largestSeqNumSeenFromFollowers, largestSeqSumSeenByFollower)

	fmt.Printf("[%d] Received propolsal ACK from node %d Nº Nodes voted for me: %d Min nodes needed: %d Quorum satiesfied: %t\n",
		this.GetNodeId(), message[0].NodeIdSender, this.nNodesThatHaveAckElectionProposal, nNodesQuorum, isQuorumSatisfied)

	if isQuorumSatisfied {
		fmt.Printf("[%d] Quorum leader proposal satiesfied. New leader elected: %d With new SeqNum %d Sending commit to the rest of the nodes\n",
			this.GetNodeId(), this.GetNodeId(), largestSeqSumSeenByFollower)

		this.node.Broadcast(nodes.CreateMessage(
			nodes.WithNodeId(this.GetNodeId()),
			nodes.WithType(zab.MESSAGE_ELECTION_COMMIT),
			nodes.WithContentUInt32(largestSeqSumSeenByFollower),
			nodes.WithFlags(nodes.FLAG_BYPASS_LEADER, nodes.FLAG_BYPASS_ORDERING)))

		this.changeStateFromElectionToBroadcast(this.GetNodeId(), this.largestSeqNumSeenFromFollowers)
	}
}

func (this *ZabNode) handleElectionCommitMessage(messages []*nodes.Message) {
	message := messages[0]

	fmt.Printf("[%d] Received commit from node %d New leader established!\n", this.GetNodeId(), message.NodeIdSender)

	this.changeStateFromElectionToBroadcast(message.NodeIdSender, message.GetContentToUint32())
}

func (this *ZabNode) changeStateFromElectionToBroadcast(newLeaderNodeId uint32, newSeqNum uint32) {
	this.state = BROADCAST
	this.leaderNodeId = newLeaderNodeId

	zabBroadcaster := this.node.GetBroadcaster().(*zab.ZabBroadcaster)
	zabBroadcaster.OnNewLeader(newLeaderNodeId, newSeqNum)

	this.heartbeatTimerTimeout = time.NewTimer(time.Duration(this.heartbeatTimeoutMs * uint64(time.Millisecond)))
	this.node.EnableBroadcast()
	this.electionLock.Unlock()
}

func (this *ZabNode) changeStateFromBroadcastToElection() {
	this.state = ELECTION
	this.node.DisableBroadcast()
	this.node.GetBroadcaster().(*zab.ZabBroadcaster).OnElectionStarted()
	this.heartbeatTimerTimeout.Stop()
}
