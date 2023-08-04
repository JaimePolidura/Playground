package raft_grpc

import (
	"context"
	"distributed-systems/src/raft_grpc"
	"distributed-systems/src/raft_grpc/grpc/proto"
	"distributed-systems/src/raft_grpc/messages"
	"google.golang.org/grpc"
	"net"
	"strconv"
)

type RaftGRPCServer struct {
	proto.UnimplementedRaftNodeServer

	raftNode *raft_grpc.RaftNode
}

func CreateRaftGRPCServerAndRun(raftNode *raft_grpc.RaftNode) *RaftGRPCServer {
	raftGrpcServer := &RaftGRPCServer{raftNode: raftNode}

	lis, _ := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(int(raftNode.Port)))
	grpcServer := grpc.NewServer()
	proto.RegisterRaftNodeServer(grpcServer, raftGrpcServer)

	go grpcServer.Serve(lis)

	return raftGrpcServer
}

func (this *RaftGRPCServer) RequestVote(context context.Context, request *proto.RequestVoteRequest) (*proto.RequestVoteResponse, error) {
	response := this.raftNode.RequestVote(context, &messages.RequestVoteRequest{
		Term:         *request.Term,
		CandidateId:  *request.CandidateId,
		LastLogIndex: *request.LastLogIndex,
		LastLogTerm:  *request.LastLogTerm,
	})

	return &proto.RequestVoteResponse{
		Term:        &response.Term,
		VoteGranted: &response.VoteGranted,
	}, nil
}

func (this *RaftGRPCServer) ReceiveLeaderHeartbeat(context context.Context, request *proto.HeartbeatRequest) (*proto.Void, error) {
	this.raftNode.ReceiveLeaderHealthCheck(context, &messages.HeartbeatRequest{Term: *request.Term})

	return &proto.Void{}, nil
}

func (this *RaftGRPCServer) AppendEntries(context context.Context, request *proto.AppendEntriesRequest) (*proto.AppendEntriesResponse, error) {
	entries := make([]messages.Entry, len(request.Entries))
	for index, entry := range request.Entries {
		entries[index] = messages.Entry{Index: *entry.Index, Term: *entry.Term}
	}

	response := this.raftNode.AppendEntries(context, &messages.AppendEntriesRequest{
		Term:         *request.Term,
		LeaderId:     *request.LeaderId,
		PrevLogIndex: *request.PrevLogTerm,
		PrevLogTerm:  *request.PrevLogTerm,
		LeaderCommit: *request.LeaderCommit,
		Entries:      entries,
	})

	return &proto.AppendEntriesResponse{
		Term:    &response.Term,
		Success: &response.Success,
	}, nil
}
