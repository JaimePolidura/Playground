package broadcast

type NodeConnectionsStore struct {
	restNodesConnections map[uint32]NodeConnection
}

func CreateNodeConnectionStore() *NodeConnectionsStore {
	return &NodeConnectionsStore{
		restNodesConnections: make(map[uint32]NodeConnection),
	}
}

func (store *NodeConnectionsStore) Add(connection NodeConnection) {
	store.restNodesConnections[connection.nodeId] = connection
}

func (store *NodeConnectionsStore) ToArrayNodeConnections() []NodeConnection {
	var nodeConnections []NodeConnection

	for key := range store.restNodesConnections {
		nodeConnections = append(nodeConnections, store.restNodesConnections[key])
	}

	return nodeConnections
}
