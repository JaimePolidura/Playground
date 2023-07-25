package broadcast

import "fmt"

type BroadcasterNode struct {
	nodeId uint32
	port   uint16

	nodeConnectionsStore *NodeConnectionsStore
	messageListener      *MessageListener
	broadcaster          Broadcaster
}

func CreateBroadcasterNode(nodeId uint32, port uint16, broadcaster Broadcaster) *BroadcasterNode {
	nodeConnectionsStore := CreateNodeConnectionStore()
	broadcaster.SetNodeConnectionsStore(nodeConnectionsStore)

	return &BroadcasterNode{
		nodeId:               nodeId,
		port:                 port,
		broadcaster:          broadcaster,
		messageListener:      CreateMessageListener(nodeId, port),
		nodeConnectionsStore: nodeConnectionsStore,
	}
}

func (broadcasterNode *BroadcasterNode) AddOtherNode(nodeId uint32, port uint32) {
	if nodeId != broadcasterNode.nodeId {
		broadcasterNode.nodeConnectionsStore.Add(nodeId, port)
	}
}

func (broadcasterNode *BroadcasterNode) Broadcast(content string) {
	broadcasterNode.broadcaster.Broadcast(CreateMessage(broadcasterNode.nodeId, broadcasterNode.nodeId, content))
}

func (broadcasterNode *BroadcasterNode) OpenConnectionsToNodes(nodes []*BroadcasterNode) {
	for _, node := range nodes {
		if node.nodeId != broadcasterNode.nodeId {
			broadcasterNode.nodeConnectionsStore.Open(node.nodeId)
		}
	}
}

func (broadcasterNode *BroadcasterNode) StartListening() {
	broadcasterNode.messageListener.ListenAsync(func(message []*BroadcastMessage) {
		broadcasterNode.broadcaster.OnBroadcastMessage(message, func(newMessage *BroadcastMessage) {
			fmt.Printf("[%d] RECIEVED UNIQUE MESSAGE \"%s\"\n", broadcasterNode.nodeId, newMessage.Content)
		})
	})
}
