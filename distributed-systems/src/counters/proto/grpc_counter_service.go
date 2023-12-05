package proto

import (
	"context"
	"distributed-systems/src/counters"
	"distributed-systems/src/counters/proto/counters_grpc"
	"net"
	"strconv"

	"google.golang.org/grpc"
)

type CounterGRPCServer struct {
	counters_grpc.UnimplementedCounterNodeServer

	nativeServer *grpc.Server
	node         *counters.Node
}

func CreateCounterGRPCServer(node *counters.Node) *CounterGRPCServer {
	countersGrpcServer := &CounterGRPCServer{node: node}

	lis, _ := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(int(node.Port)))

	grpcServer := grpc.NewServer()

	counters_grpc.RegisterCounterNodeServer(grpcServer, countersGrpcServer)

	countersGrpcServer.nativeServer = grpcServer

	go grpcServer.Serve(lis)

	return countersGrpcServer
}

func (c *CounterGRPCServer) Update(ctx context.Context, request *counters_grpc.UpdateCounterRequest) (*counters_grpc.UpdateCounterResponse, error) {
	res := c.node.OnUpdateFromNode(ctx, counters.UpdateCounterRequest{
		IsIncrement:               *request.IsIncrement,
		NextSelfSeqValue:          *request.NextSelfSeqValue,
		LastSeqValueSeenIncrement: *request.LastSeqValueSeenIncrement,
		LastSeqValueSeenDecrement: *request.LastSeqValueSeenDecrement,
		SelfNodeId:                *request.NodeId,
	})

	return &counters_grpc.UpdateCounterResponse{
		NeedsSyncIncrement:              &res.NeedsSyncIncrement,
		NextSelfSeqValueToSyncIncrement: &res.NextSelfSeqValueToSyncIncrement,
		NeedsSyncDecrement:              &res.NeedsSyncDecrement,
		NextSelfSeqValueToSyncDecrement: &res.NextSelfSeqValueToSyncDecrement,
	}, nil
}

func (c *CounterGRPCServer) Stop() {
	c.nativeServer.Stop()
}
