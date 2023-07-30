package types

// Skips order, ack. Used for heartbeats
const FLAG_BYPASS_ORDERING = 1

// The follower will communicate to the other followers without pass through the leader
// Used for when leader has failed
const FLAG_BYPASS_LEADER = 2

const MESSAGE_BROADCAST = 0
const MESSAGE_ACK = 1
const MESSAGE_HEARTBEAT = 3
const MESSAGE_ZAB_ELECTION_FAILURE_DETECTED = 4
const MESSAGE_ZAB_ELECTION_PROPOSAL = 5
const MESSAGE_ZAB_ELECTION_ACK_PROPOSAL = 6
const MESSAGE_ZAB_ELECTION_COMMIT = 7
const MESSAGE_DO_BROADCAST = 8
const MESSAGE_NODE_STOPPED = 9

const MESSAGE_PAXOS_PREPARE = 10
const MESSAGE_PAXOS_PROMISE = 12
const MESSAGE_PAXOS_PROMISE_ACCEPT = 13
const MESSAGE_PAXOS_ACCEPT = 14
const MESSAGE_PAXOS_ACCEPTED = 15
