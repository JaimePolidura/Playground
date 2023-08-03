package zab

import (
	"distributed-systems/src/nodes"
	"distributed-systems/src/nodes/types"
	"distributed-systems/src/utils"
	"fmt"
	"time"
)

func (this *ZabNode) handleHeartbeatMessage(message *nodes.Message) {
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

			fmt.Printf("[%d] Detected cantidate %d failure. Broadcasting MESSAGE_ZAB_ELECTION_FAILURE_DETECTED\n",
				this.GetNodeId(), this.nodesIdRing[this.getBackIndexInRingByIndex(this.getRingIndexByNodeId(this.GetNodeId()))])

			this.Broadcast(nodes.CreateMessage(
				nodes.WithNodeId(this.GetNodeId()),
				nodes.WithType(types.MESSAGE_ZAB_ELECTION_FAILURE_DETECTED),
				nodes.WithFlags(types.FLAG_BYPASS_LEADER, types.FLAG_BYPASS_ORDERING, types.FLAG_BROADCAST),
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
				fmt.Printf("[%d] Detected leader %d failure. Broadcasting MESSAGE_ZAB_ELECTION_FAILURE_DETECTED\n",
					this.GetNodeId(), this.leaderNodeId)

				this.changeStateFromBroadcastToElection(this.leaderNodeId)

				this.Broadcast(nodes.CreateMessage(
					nodes.WithNodeId(this.GetNodeId()),
					nodes.WithType(types.MESSAGE_ZAB_ELECTION_FAILURE_DETECTED),
					nodes.WithFlags(types.FLAG_BYPASS_LEADER, types.FLAG_BYPASS_ORDERING, types.FLAG_BROADCAST),
					nodes.WithContentUInt32(this.leaderNodeId)))

				if this.leaderNodeId == this.prevNodeRing {
					this.proposeMySelfAsTheNewLeader()
				}
			}
		}
	}
}

func (this *ZabNode) handleNodeFailureMessage(message *nodes.Message) {
	failedNodeId := message.GetContentToUint32()

	if this.isElectionAlreadyOnGoingByFailedNode(failedNodeId) {
		return
	}

	this.changeStateFromBroadcastToElection(failedNodeId)

	failureNodeIsPrev := this.prevNodeRing == failedNodeId

	fmt.Printf("[%d] Receieved MESSAGE_ZAB_ELECTION_FAILURE_DETECTED from node %d of failed node %d. Is prev node? %t\n",
		this.GetNodeId(), failedNodeId, message.NodeIdSender, failureNodeIsPrev)

	if failureNodeIsPrev {
		this.proposeMySelfAsTheNewLeader()
	}
}

func (this *ZabNode) proposeMySelfAsTheNewLeader() {
	fmt.Printf("[%d] Sending proposal to be leader to followers\n", this.GetNodeId())

	go this.startSendingHeartbeats()

	this.Broadcast(nodes.CreateMessage(
		nodes.WithNodeId(this.GetNodeId()),
		nodes.WithType(types.MESSAGE_ZAB_ELECTION_PROPOSAL),
		nodes.WithFlags(types.FLAG_BYPASS_LEADER, types.FLAG_BYPASS_ORDERING, types.FLAG_BROADCAST)))
}

func (this *ZabNode) handleElectionProposalMessage(message *nodes.Message) {
	proposerNodeId := message.NodeIdSender
	largestSeqNum := this.GetBroadcaster().(*ZabBroadcaster).GetLargestSeqNumbReachievedLeader()

	fmt.Printf("[%d] Accepting proposal from node %d to be the new leader\n", this.GetNodeId(), proposerNodeId)

	this.GetConnectionManager().Send(proposerNodeId, nodes.CreateMessage(
		nodes.WithContentUInt32(largestSeqNum),
		nodes.WithNodeId(this.GetNodeId()),
		nodes.WithType(types.MESSAGE_ZAB_ELECTION_ACK_PROPOSAL),
		nodes.WithFlags(types.FLAG_BYPASS_LEADER, types.FLAG_BYPASS_ORDERING)))
}

func (this *ZabNode) handleElectionAckProposalMessage(message *nodes.Message) {
	nodeIdVoter := message.NodeIdSender

	if this.hasAlreadyVoted(this.GetNodeId(), nodeIdVoter) {
		return
	}

	this.nNodesThatHaveAckElectionProposal++
	largestSeqSumSeenByFollower := message.GetContentToUint32()
	nNodesQuorum := this.GetConnectionManager().GetNumberConnections()/2 + 1
	isQuorumSatisfied := this.nNodesThatHaveAckElectionProposal+1 >= nNodesQuorum //The candidate votes it self
	this.largestSeqNumSeenFromFollowers = utils.MaxUint32(this.largestSeqNumSeenFromFollowers, largestSeqSumSeenByFollower)
	this.registerNodeVote(this.GetNodeId(), nodeIdVoter)

	fmt.Printf("[%d] Received propolsal ACK from node %d Nº Nodes voted for me: %d Min nodes needed: %d Quorum satisfied: %t\n",
		this.GetNodeId(), message.NodeIdSender, this.nNodesThatHaveAckElectionProposal, nNodesQuorum, isQuorumSatisfied)

	if isQuorumSatisfied {
		fmt.Printf("[%d] Quorum leader proposal satisfied. New leader elected: %d With new SeqNum %d Sending commit to the rest of the nodes\n",
			this.GetNodeId(), this.GetNodeId(), largestSeqSumSeenByFollower)

		this.Broadcast(nodes.CreateMessage(
			nodes.WithNodeId(this.GetNodeId()),
			nodes.WithType(types.MESSAGE_ZAB_ELECTION_COMMIT),
			nodes.WithContentUInt32(largestSeqSumSeenByFollower),
			nodes.WithFlags(types.FLAG_BYPASS_LEADER, types.FLAG_BYPASS_ORDERING, types.FLAG_BROADCAST)))

		this.changeStateFromElectionToBroadcast(this.GetNodeId(), this.largestSeqNumSeenFromFollowers)
	}
}

func (this *ZabNode) handleElectionCommitMessage(message *nodes.Message) {
	fmt.Printf("[%d] Received commit from node %d New leader established!\n", this.GetNodeId(), message.NodeIdSender)

	this.changeStateFromElectionToBroadcast(message.NodeIdSender, message.GetContentToUint32())
}

func (this *ZabNode) changeStateFromElectionToBroadcast(newLeaderNodeId uint32, newSeqNum uint32) {
	this.state = BROADCAST
	this.leaderNodeId = newLeaderNodeId

	zabBroadcaster := this.GetBroadcaster().(*ZabBroadcaster)
	zabBroadcaster.OnNewLeader(newLeaderNodeId, newSeqNum)

	this.nNodesThatHaveAckElectionProposal = 0

	this.heartbeatTimerTimeout.Reset(time.Duration(this.heartbeatTimeoutMs * uint64(time.Millisecond)))

	if this.heartbeatCandidateTimerTimeout != nil {
		this.heartbeatCandidateTimerTimeout.Stop()
	}

	utils.ClearMap(this.nodesVotesRegistry[newLeaderNodeId])

	this.EnableBroadcast()
}

func (this *ZabNode) changeStateFromBroadcastToElection(failedNode uint32) {
	this.state = ELECTION
	this.DisableBroadcast()
	this.GetBroadcaster().(*ZabBroadcaster).OnElectionStarted()
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

func (this *ZabNode) setupHeartbeatCandidateTimerTimeout(failedNode uint32) { //NodeId <- Self
	if this.heartbeatCandidateTimerTimeout != nil {
		this.heartbeatCandidateTimerTimeout.Stop()
	}

	indexFailedNode := this.getRingIndexByNodeId(failedNode)
	indexSelfNode := this.getRingIndexByNodeId(this.GetNodeId())
	indexCandidate := this.getNextIndexInRingByIndex(indexFailedNode)
	distanceFromFailedCandidate := this.getRingClockwiseDistanceByIndex(indexCandidate, indexSelfNode)

	if distanceFromFailedCandidate == 0 {
		return
	}
	if distanceFromFailedCandidate <= this.nHeartbeatCandidateTimersTimeout {
		timeout := time.Duration(uint64(distanceFromFailedCandidate) * (this.heartbeatTimeoutMs + 1) * uint64(time.Millisecond))
		this.heartbeatCandidateTimerTimeout = time.NewTimer(timeout)
		go this.startCandidateHeartbeatTimerTimeout()
	}
}

func (this *ZabNode) restartHeartbeatCandidateTimerTimeout() {
	indexFailedNode := this.getRingIndexByNodeId(this.lastFailedNodeId)
	indexSelfNode := this.getRingIndexByNodeId(this.GetNodeId())
	indexCandidate := this.getNextIndexInRingByIndex(indexFailedNode)

	distanceFromFailedCandidate := this.getRingClockwiseDistanceByIndex(indexCandidate, indexSelfNode)
	timeout := time.Duration(uint64(distanceFromFailedCandidate) * (this.heartbeatTimeoutMs + 1) * uint64(time.Millisecond))
	this.heartbeatCandidateTimerTimeout.Reset(timeout)
}
