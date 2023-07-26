package broadcast

import "fmt"

type BroadcasterNode struct {
	selfNodeId uint32
	port       uint16

	nodeConnectionsStore *NodeConnectionsStore
	messageListener      *MessageListener
	broadcaster          Broadcaster
}

func CreateBroadcasterNode(nodeId uint32, port uint16, broadcaster Broadcaster) *BroadcasterNode {
	nodeConnectionsStore := CreateNodeConnectionStore()
	broadcaster.SetNodeConnectionsStore(nodeConnectionsStore)

	return &BroadcasterNode{
		selfNodeId:           nodeId,
		port:                 port,
		broadcaster:          broadcaster,
		messageListener:      CreateMessageListener(nodeId, port),
		nodeConnectionsStore: nodeConnectionsStore,
	}
}

func (this *BroadcasterNode) AddOtherNode(otherNodeId uint32, port uint32) {
	if otherNodeId != this.selfNodeId {
		this.nodeConnectionsStore.Add(otherNodeId, port, this.selfNodeId)
	}
}

func (this *BroadcasterNode) Broadcast(content string) {
	this.broadcaster.Broadcast(CreateMessage(this.selfNodeId, this.selfNodeId, content))
}

func (this *BroadcasterNode) OpenConnectionsToNodes(nodes []*BroadcasterNode) {
	for _, node := range nodes {
		if node.selfNodeId != this.selfNodeId {
			this.nodeConnectionsStore.Open(node.selfNodeId)
		}
	}
}

func (this *BroadcasterNode) StartListening() {
	this.messageListener.ListenAsync(func(message []*BroadcastMessage) {
		this.broadcaster.OnBroadcastMessage(message, func(newMessage *BroadcastMessage) {
			fmt.Printf("[%d] RECIEVED UNIQUE MESSAGE \"%s\"\n", this.selfNodeId, newMessage.Content)
		})
	})
}
