package nodes

type NodeConnectionsStore struct {
	restNodesConnections map[uint32]*NodeConnection
}

func CreateNodeConnectionStore() *NodeConnectionsStore {
	return &NodeConnectionsStore{
		restNodesConnections: make(map[uint32]*NodeConnection),
	}
}

func (this *NodeConnectionsStore) Size() uint32 {
	return uint32(len(this.restNodesConnections))
}

func (this *NodeConnectionsStore) AddConnection(connection *NodeConnection) {
	this.restNodesConnections[connection.nodeId] = connection
}

func (this *NodeConnectionsStore) Add(otherNodeId uint32, otherNodePort uint32, selfNodeId uint32) {
	this.restNodesConnections[otherNodeId] = CreateNodeConnection(otherNodeId, otherNodePort, selfNodeId)
}

func (this *NodeConnectionsStore) Contains(nodeId uint32) bool {
	_, contained := this.restNodesConnections[nodeId]
	return contained
}

func (this *NodeConnectionsStore) Open(nodeId uint32) {
	nodeConnection := this.restNodesConnections[nodeId]
	nodeConnection.Open()
}

func (this *NodeConnectionsStore) CloseAll() {
	for _, connection := range this.restNodesConnections {
		connection.Close()
	}
}

func (this *NodeConnectionsStore) Get(nodeId uint32) *NodeConnection {
	return this.restNodesConnections[nodeId]
}

func (this *NodeConnectionsStore) ToArrayNodeConnections() []*NodeConnection {
	nodeConnections := make([]*NodeConnection, 0)

	for key := range this.restNodesConnections {
		nodeConnections = append(nodeConnections, this.restNodesConnections[key])
	}

	return nodeConnections
}
