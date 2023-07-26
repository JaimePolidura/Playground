package broadcast

type NodeConnectionsStore struct {
	restNodesConnections map[uint32]*NodeConnection
}

func CreateNodeConnectionStore() *NodeConnectionsStore {
	return &NodeConnectionsStore{
		restNodesConnections: make(map[uint32]*NodeConnection),
	}
}

func (store *NodeConnectionsStore) Size() uint32 {
	return uint32(len(store.restNodesConnections))
}

func (store *NodeConnectionsStore) Add(otherNodeId uint32, otherNodePort uint32, selfNodeId uint32) {
	store.restNodesConnections[otherNodeId] = CreateNodeConnection(otherNodeId, otherNodePort, selfNodeId)
}

func (store *NodeConnectionsStore) Contains(nodeId uint32) bool {
	_, contained := store.restNodesConnections[nodeId]
	return contained
}

func (store *NodeConnectionsStore) Open(nodeId uint32) {
	nodeConnection := store.restNodesConnections[nodeId]
	nodeConnection.Open()
}

func (store *NodeConnectionsStore) Get(nodeId uint32) *NodeConnection {
	return store.restNodesConnections[nodeId]
}

func (store *NodeConnectionsStore) ToArrayNodeConnections() []*NodeConnection {
	nodeConnections := make([]*NodeConnection, 0)

	for key := range store.restNodesConnections {
		nodeConnections = append(nodeConnections, store.restNodesConnections[key])
	}

	return nodeConnections
}
