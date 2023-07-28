package nodes

type ConnectionManager struct {
	selfNodeId uint32
	selfPort   uint16

	messageListener  *MessageListener
	connectionsStore *NodeConnectionsStore

	NewMessage chan *Message
}

func CreateConnectionManager(nodeId uint32, port uint16) *ConnectionManager {
	return &ConnectionManager{
		selfNodeId:       nodeId,
		selfPort:         port,
		messageListener:  CreateMessageListener(nodeId, port),
		connectionsStore: CreateNodeConnectionStore(),
		NewMessage:       make(chan *Message, 100),
	}
}

func (this *ConnectionManager) Stop() {
	this.messageListener.Stop()
	this.connectionsStore.CloseAll()
}

func (this *ConnectionManager) OpenAllConnections() {
	for _, connection := range this.connectionsStore.restNodesConnections {
		connection.Open()
	}
}

func (this *ConnectionManager) GetNumberConnections() uint32 {
	return this.connectionsStore.Size()
}

func (this *ConnectionManager) Add(nodeId uint32, port uint32) {
	if !this.connectionsStore.Contains(nodeId) {
		this.connectionsStore.Add(nodeId, port, this.selfNodeId)
	}
}

func (this *ConnectionManager) Open(nodeId uint32) {
	if this.connectionsStore.Contains(nodeId) {
		this.connectionsStore.Get(nodeId).Open()
	}
}

func (this *ConnectionManager) SendAllExcept(nodeIdExcept uint32, message *Message) {
	for _, connection := range this.connectionsStore.ToArrayNodeConnections() {
		if connection.nodeId != nodeIdExcept {
			connection.Write(message)
		}
	}
}

func (this *ConnectionManager) Send(nodeId uint32, message *Message) {
	if this.connectionsStore.Contains(nodeId) {
		this.connectionsStore.Get(nodeId).Write(message)
	}
}

func (this *ConnectionManager) ForEachConnectionExcept(exceptNodeId uint32, consumer func(connection *NodeConnection)) {
	for _, connection := range this.connectionsStore.ToArrayNodeConnections() {
		if exceptNodeId != connection.nodeId {
			consumer(connection)
		}
	}
}

func (this *ConnectionManager) GetNodesId() []uint32 {
	nodesId := make([]uint32, this.connectionsStore.Size())

	for index, node := range this.connectionsStore.restNodesConnections {
		nodesId[index] = node.nodeId
	}

	return nodesId
}

func (this *ConnectionManager) StartListeningAsync() {
	this.messageListener.ListenAsync()

	go func() {
		for newMessage := range this.messageListener.newMessage {
			this.NewMessage <- newMessage.message
		}
	}()
}
