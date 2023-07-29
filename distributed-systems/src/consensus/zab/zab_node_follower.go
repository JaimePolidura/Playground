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

	if this.heartbeatCandidateTimerTimeout != nil {
		this.restartHeartbeatCandidateTimerTimeout()
	}
}

func (this *ZabNode) startCandidateHeartbeatTimerTimeout() {
	for {
		select {
		case <-this.heartbeatCandidateTimerTimeout.C:
			if this.state == STOPPED {
				return
			}

			fmt.Printf("[%d] Detected cantidate %d failure. Broadcasting MESSAGE_ELECTION_FAILURE_DETECTED\n",
				this.GetNodeId(), this.nodesIdRing[this.getBackIndexInRingByIndex(this.getRingIndexByNodeId(this.GetNodeId()))])

			this.node.Broadcast(nodes.CreateMessage(
				nodes.WithNodeId(this.GetNodeId()),
				nodes.WithType(zab.MESSAGE_ELECTION_FAILURE_DETECTED),
				nodes.WithFlags(nodes.FLAG_BYPASS_LEADER, nodes.FLAG_BYPASS_ORDERING),
				nodes.WithContentUInt32(this.prevNodeRing)))
			this.proposeMySelfAsTheNewLeader()
		}
	}
}

func (this *ZabNode) startLeaderHeartbeatTimerTimeout() {
	for {
		select {
		case <-this.heartbeatTimerTimeout.C:
			if this.state == STOPPED {
				return
			}

			if this.IsFollower() && this.state == BROADCAST && !this.isElectionAlreadyOnGoingByFailedNode(this.leaderNodeId) {
				fmt.Printf("[%d] Detected leader %d failure. Broadcasting MESSAGE_ELECTION_FAILURE_DETECTED\n",
					this.GetNodeId(), this.leaderNodeId)

				this.changeStateFromBroadcastToElection(this.leaderNodeId)

				this.node.Broadcast(nodes.CreateMessage(
					nodes.WithNodeId(this.GetNodeId()),
					nodes.WithType(zab.MESSAGE_ELECTION_FAILURE_DETECTED),
					nodes.WithFlags(nodes.FLAG_BYPASS_LEADER, nodes.FLAG_BYPASS_ORDERING),
					nodes.WithContentUInt32(this.leaderNodeId)))

				if this.leaderNodeId == this.prevNodeRing {
					this.proposeMySelfAsTheNewLeader()
				}
			}
		}
	}
}

func (this *ZabNode) handleNodeFailureMessage(messages []*nodes.Message) {
	message := messages[0]
	failedNodeId := message.GetContentToUint32()

	if this.isElectionAlreadyOnGoingByFailedNode(failedNodeId) {
		return
	}

	this.changeStateFromBroadcastToElection(failedNodeId)

	failureNodeIsPrev := this.prevNodeRing == failedNodeId

	fmt.Printf("[%d] Receieved MESSAGE_ELECTION_FAILURE_DETECTED from node %d of failed node %d. Is prev node? %t\n",
		this.GetNodeId(), failedNodeId, message.NodeIdSender, failureNodeIsPrev)

	if failureNodeIsPrev {
		this.proposeMySelfAsTheNewLeader()
	}
}

func (this *ZabNode) proposeMySelfAsTheNewLeader() {
	fmt.Printf("[%d] Sending proposal to be leader to followers\n", this.GetNodeId())

	go this.startSendingHeartbeats()

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

func (this *ZabNode) handleElectionAckProposalMessage(messages []*nodes.Message) {
	message := messages[0]
	nodeIdVoter := message.NodeIdSender

	if this.hasAlreadyVoted(this.GetNodeId(), nodeIdVoter) {
		return
	}

	this.nNodesThatHaveAckElectionProposal++
	largestSeqSumSeenByFollower := messages[0].GetContentToUint32()
	nNodesQuorum := this.GetConnectionManager().GetNumberConnections()/2 + 1
	isQuorumSatisfied := this.nNodesThatHaveAckElectionProposal+1 >= nNodesQuorum //The candidate votes it self
	this.largestSeqNumSeenFromFollowers = utils.MaxUint32(this.largestSeqNumSeenFromFollowers, largestSeqSumSeenByFollower)
	this.registerNodeVote(this.GetNodeId(), nodeIdVoter)

	fmt.Printf("[%d] Received propolsal ACK from node %d Nº Nodes voted for me: %d Min nodes needed: %d Quorum satiesfied: %t\n",
		this.GetNodeId(), messages[0].NodeIdSender, this.nNodesThatHaveAckElectionProposal, nNodesQuorum, isQuorumSatisfied)

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

	this.nNodesThatHaveAckElectionProposal = 0

	this.heartbeatTimerTimeout.Reset(time.Duration(this.heartbeatTimeoutMs * uint64(time.Millisecond)))

	if this.heartbeatCandidateTimerTimeout != nil {
		this.heartbeatCandidateTimerTimeout.Stop()
	}

	utils.ClearMap(this.nodesVotesRegistry[newLeaderNodeId])

	this.node.EnableBroadcast()
}

func (this *ZabNode) changeStateFromBroadcastToElection(failedNode uint32) {
	this.state = ELECTION
	this.node.DisableBroadcast()
	this.node.GetBroadcaster().(*zab.ZabBroadcaster).OnElectionStarted()
	this.heartbeatTimerTimeout.Stop()

	if failedNode != this.GetNodeId() {
		this.setupHeartbeatCandidateTimerTimeout(failedNode)
	}
}

func (this *ZabNode) isElectionAlreadyOnGoingByFailedNode(node uint32) bool {
	this.lastFailedNodeLock.Lock()
	if this.lastFailedNodeId == this.leaderNodeId {
		this.lastFailedNodeLock.Unlock()
		return true
	}
	this.lastFailedNodeId = this.leaderNodeId
	this.lastFailedNodeLock.Unlock()

	return false
}

func (this *ZabNode) hasAlreadyVoted(nodeIdCandidateToVote uint32, nodeIdToCheckIfHasVoted uint32) bool {
	if _, contained := this.nodesVotesRegistry[nodeIdCandidateToVote]; !contained {
		this.nodesVotesRegistry[nodeIdCandidateToVote] = make(map[uint32]uint32)
	}

	registeredVotes, _ := this.nodesVotesRegistry[nodeIdCandidateToVote]
	_, contained := registeredVotes[nodeIdToCheckIfHasVoted]

	return contained
}

func (this *ZabNode) registerNodeVote(nodeIdCandidateToVote uint32, nodeIdThatVoted uint32) {
	if _, contained := this.nodesVotesRegistry[nodeIdCandidateToVote]; !contained {
		this.nodesVotesRegistry[nodeIdCandidateToVote] = make(map[uint32]uint32)
	}

	registeredVotes, _ := this.nodesVotesRegistry[nodeIdCandidateToVote]
	registeredVotes[nodeIdThatVoted] = nodeIdThatVoted
}
