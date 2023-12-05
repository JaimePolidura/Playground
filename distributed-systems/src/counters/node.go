package counters

import (
	"context"
	"sync/atomic"

	"golang.org/x/sys/cpu"
)

type Node struct {
	NodeId int
	Port   uint16

	Peers    []CounterNodeService
	counters []CounterByNode
}

type CounterByNode struct {
	_             cpu.CacheLinePad
	SeqIncrements uint64
	_             cpu.CacheLinePad
	SeqDecrements uint64
}

type CounterNodeService interface {
	Update(ctx context.Context, request *UpdateCounterRequest) *UpdateCounterResponse
}

func (n *Node) Get() uint64 {
	acc := uint64(0)

	for _, counterByNode := range n.counters {
		acc += counterByNode.SeqIncrements - counterByNode.SeqDecrements
	}

	return acc
}

func (n *Node) Increment() uint64 {
	counterThisNode := &n.counters[n.NodeId]
	nIncrementsPtr := &counterThisNode.SeqIncrements
	newSeqValue := atomic.AddUint64(nIncrementsPtr, 1)

	n.replicate(newSeqValue, true)

	return newSeqValue
}

func (n *Node) Decrement() uint64 {
	counterThisNode := &n.counters[n.NodeId]
	nIncrementsPtr := &counterThisNode.SeqDecrements
	newSeqValue := atomic.AddUint64(nIncrementsPtr, 1)

	n.replicate(newSeqValue, false)

	return newSeqValue
}

func (n *Node) replicate(newValue uint64, isIncrement bool) {
	for peerNodeId, peer := range n.Peers {
		if peerNodeId != n.NodeId {
			go func(peer CounterNodeService, peerNodeId int) {
				updateResponse := peer.Update(context.Background(), &UpdateCounterRequest{
					LastSeqValueSeenIncrement: atomic.LoadUint64(&n.counters[peerNodeId].SeqIncrements),
					LastSeqValueSeenDecrement: atomic.LoadUint64(&n.counters[peerNodeId].SeqDecrements),
					SelfNodeId:                uint32(n.NodeId),
					IsIncrement:               isIncrement,
					NextSelfSeqValue:          newValue,
				})

				if updateResponse.NeedsSyncIncrement {
					n.updatePeerCounterSeqToSync(peerNodeId, updateResponse.NextSelfSeqValueToSyncIncrement, true)
				}
				if updateResponse.NeedsSyncDecrement {
					n.updatePeerCounterSeqToSync(peerNodeId, updateResponse.NextSelfSeqValueToSyncDecrement, false)
				}
			}(peer, peerNodeId)
		}
	}
}

func (n *Node) OnUpdateFromNode(ctx context.Context, request UpdateCounterRequest) UpdateCounterResponse {
	n.updatePeerCounterSeqToSync(int(request.SelfNodeId), request.NextSelfSeqValue, request.IsIncrement)

	var response UpdateCounterResponse

	selfSeqCounterDecrementValue := atomic.LoadUint64(&n.counters[n.NodeId].SeqDecrements)
	if selfSeqCounterDecrementValue > request.LastSeqValueSeenDecrement {
		response.NextSelfSeqValueToSyncDecrement = selfSeqCounterDecrementValue
		response.NeedsSyncDecrement = true
	}

	selfSeqCounterIncrementValue := atomic.LoadUint64(&n.counters[n.NodeId].SeqIncrements)
	if selfSeqCounterIncrementValue > request.LastSeqValueSeenIncrement {
		response.NextSelfSeqValueToSyncIncrement = selfSeqCounterIncrementValue
		response.NeedsSyncIncrement = true
	}

	return response
}

func (n *Node) updatePeerCounterSeqToSync(peerId int, newSeqValue uint64, isIncrement bool) {
	actualCounterSeqPtr := n.getCounterSeqPtr(peerId, isIncrement)
	actualCounterSeqVal := atomic.LoadUint64(actualCounterSeqPtr)

	for newSeqValue > actualCounterSeqVal && atomic.CompareAndSwapUint64(actualCounterSeqPtr, actualCounterSeqVal, newSeqValue) {
		actualCounterSeqPtr = n.getCounterSeqPtr(peerId, isIncrement)
		actualCounterSeqVal = atomic.LoadUint64(actualCounterSeqPtr)
	}
}

func (n *Node) getCounterSeqPtr(peerId int, isIncrement bool) *uint64 {
	if isIncrement {
		return &n.counters[peerId].SeqIncrements
	} else {
		return &n.counters[peerId].SeqDecrements
	}
}

func CreateNode(nodeId int, port uint16) *Node {
	return &Node{
		NodeId: nodeId,
		Port:   port,
		Peers:  make([]CounterNodeService, 0),
	}
}

func (n *Node) SetPeers(peers []CounterNodeService) {
	for _, peer := range peers {
		n.Peers = append(n.Peers, peer)
	}

	n.counters = make([]CounterByNode, len(peers))
}
