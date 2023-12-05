package proto

import (
	"context"
	"distributed-systems/src/counters"
	"distributed-systems/src/counters/proto/counters_grpc"
	"strconv"

	"google.golang.org/grpc"
)

type CounterGRPCClient struct {
	nativeClient counters_grpc.CounterNodeClient
}

func CreateCounterGRPCClient(otherPort uint16) *CounterGRPCClient {
	conn, err := grpc.Dial("127.0.0.1:"+strconv.Itoa(int(otherPort)), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	grpcClient := counters_grpc.NewCounterNodeClient(conn)

	return &CounterGRPCClient{
		nativeClient: grpcClient,
	}
}

func (r *CounterGRPCClient) Update(ctx context.Context, request *counters.UpdateCounterRequest) *counters.UpdateCounterResponse {
	response, err := r.nativeClient.Update(ctx, &counters_grpc.UpdateCounterRequest{
		IsIncrement:               &request.IsIncrement,
		NextSelfSeqValue:          &request.NextSelfSeqValue,
		LastSeqValueSeenIncrement: &request.LastSeqValueSeenIncrement,
		LastSeqValueSeenDecrement: &request.LastSeqValueSeenDecrement,
		NodeId:                    &request.SelfNodeId,
	})

	if err != nil {
		panic(err)
	}

	return &counters.UpdateCounterResponse{
		NeedsSyncIncrement:              *response.NeedsSyncIncrement,
		NextSelfSeqValueToSyncIncrement: *response.NextSelfSeqValueToSyncIncrement,
		NeedsSyncDecrement:              *response.NeedsSyncDecrement,
		NextSelfSeqValueToSyncDecrement: *response.NextSelfSeqValueToSyncDecrement,
	}
}
