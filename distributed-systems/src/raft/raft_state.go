package raft

type RaftState int

const (
	FOLLOWER = iota
	CANDIDATE
	LEADER
)
