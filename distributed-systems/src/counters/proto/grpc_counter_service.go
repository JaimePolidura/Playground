package proto

import (
	"context"
	"distributed-systems/src/counters"
	"distributed-systems/src/counters/proto/counters_grpc"
	"distributed-systems/src/raft/grpc/proto"
	"net"
	"strconv"

	"google.golang.org/grpc"
)

type CounterGRPCServer struct {
	proto.UnimplementedRaftNodeServer

	nativeServer *grpc.Server
	node         *counters.Node
}

func CreateCounterGRPCServer(node *counters.Node) *CounterGRPCServer {
	countersGrpcServer := &CounterGRPCServer{node: node}

	lis, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(int(node.Port)))
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterRaftNodeServer(grpcServer, countersGrpcServer)

	countersGrpcServer.nativeServer = grpcServer

	go grpcServer.Serve(lis)

	return countersGrpcServer
}

func (c *CounterGRPCServer) Update(ctx context.Context, request *counters_grpc.UpdateCounterRequest) (*counters.UpdateCounterResponse, error) {
	res := c.node.OnUpdateFromNode(ctx, counters.UpdateCounterRequest{
		IsIncrement:               *request.IsIncrement,
		NextSelfSeqValue:          *request.NextSelfSeqValue,
		LastSeqValueSeenIncrement: *request.LastSeqValueSeenIncrement,
		LastSeqValueSeenDecrement: *request.LastSeqValueSeenDecrement,
		SelfNodeId:                *request.NodeId,
	})

	return &counters.UpdateCounterResponse{
		NeedsSyncIncrement:              res.NeedsSyncIncrement,
		NextSelfSeqValueToSyncIncrement: res.NextSelfSeqValueToSyncIncrement,
		NeedsSyncDecrement:              res.NeedsSyncDecrement,
		NextSelfSeqValueToSyncDecrement: res.NextSelfSeqValueToSyncDecrement,
	}, nil
}

func (c *CounterGRPCServer) Stop() {
	c.nativeServer.Stop()
}
