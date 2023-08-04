package raft_grpc

import (
	"context"
	"distributed-systems/src/raft_grpc/grpc/proto"
	"distributed-systems/src/raft_grpc/messages"
	"fmt"
	"google.golang.org/grpc"
	"strconv"
)

type RaftGRPCClient struct {
	nativeClient proto.RaftNodeClient
}

func CreateRaftGRPCClient(otherPort uint16) RaftGRPCClient {
	conn, _ := grpc.Dial("127.0.0.1:"+strconv.Itoa(int(otherPort)), grpc.WithInsecure())

	grpcClient := proto.NewRaftNodeClient(conn)

	return RaftGRPCClient{
		nativeClient: grpcClient,
	}
}

func (this *RaftGRPCClient) RequestVote(context context.Context, request *messages.RequestVoteRequest) *messages.RequestVoteResponse {
	response, err := this.nativeClient.RequestVote(context, &proto.RequestVoteRequest{
		Term:         &request.Term,
		CandidateId:  &request.CandidateId,
		LastLogIndex: &request.LastLogIndex,
		LastLogTerm:  &request.LastLogTerm,
	})

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return &messages.RequestVoteResponse{
		Term:        *response.Term,
		VoteGranted: *response.VoteGranted,
	}
}

func (this *RaftGRPCClient) ReceiveLeaderHealthCheck(context context.Context, request *messages.HeartbeatRequest) {
	this.nativeClient.ReceiveLeaderHeartbeat(context, &proto.HeartbeatRequest{Term: &request.Term})
}

func (this *RaftGRPCClient) AppendEntries(context context.Context, request *messages.AppendEntriesRequest) *messages.AppendEntriesResponse {
	entries := make([]*proto.Entry, len(request.Entries))
	for index, entry := range request.Entries {
		entries[index] = &proto.Entry{Index: &entry.Index, Term: &entry.Term}
	}

	response, err := this.nativeClient.AppendEntries(context, &proto.AppendEntriesRequest{
		Term:         &request.Term,
		LeaderId:     &request.LeaderId,
		PrevLogIndex: &request.PrevLogIndex,
		PrevLogTerm:  &request.PrevLogTerm,
		Entries:      entries,
		LeaderCommit: &request.LeaderCommit,
	})

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return &messages.AppendEntriesResponse{
		Term:    *response.Term,
		Success: *response.Success,
	}
}
