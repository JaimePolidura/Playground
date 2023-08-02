package types

// Skips order, ack. Used for heartbeats
const FLAG_BYPASS_ORDERING = 1

// The follower will communicate to the other followers without pass through the leader
// Used for when leader has failed
const FLAG_BYPASS_LEADER = 2

const FLAG_BROADCAST = 3

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

const MESSAGE_MULTIPAXOS_REDIRECT_LEADER = 16
const MESSAGE_MULTIPAXOS_REPLAY_SUBMIT = 17
const MESSAGE_MULTIPAXOS_REPLAY_SUBMISSION = 18
const MESSAGE_MULTIPAXOS_ACCEPTED = 19
const MESSAGE_MULTIPAXOS_ACCEPT = 20

const MESSAGE_RAFT_REQUEST_ELECTION = 21
const MESSAGE_RAFT_OUTDATED_TERM = 22
const MESSAGE_RAFT_REQUEST_ELECTION_REJECTED_ALREADY_VOTED = 23
const MESSAGE_RAFT_REQUEST_ELECTION_VOTED = 24
const MESSAGE_RAFT_REQUEST_ELECTION_NODE_ELECTED = 25

const MESSAGE_RAFT_LOG_APPEND_ENTRIES = 26
const MESSAGE_RAFT_LOG_APPENDED_ENTRY = 27
const MESSAGE_RAFT_LOG_DO_COMMIT = 28
const MESSAGE_RAFT_LOG_OUTDATED_ENTRIES = 29
