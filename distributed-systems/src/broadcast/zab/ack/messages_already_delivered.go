package ack

type MessagesAlreadyDelivered struct {
	seqNumsAlreadyDeliveredByNodeId map[uint32]map[uint32]uint32
}

func CreateMessagesAlreadyDelivered() *MessagesAlreadyDelivered {
	return &MessagesAlreadyDelivered{
		seqNumsAlreadyDeliveredByNodeId: make(map[uint32]map[uint32]uint32),
	}
}

func (this *MessagesAlreadyDelivered) Add(nodeId uint32, seqNum uint32) {
	seqNumsByNodeId := this.getSeqNumsByNodeId(nodeId)
	seqNumsByNodeId[seqNum] = seqNum
}

func (this *MessagesAlreadyDelivered) IsAlreadyDelivered(nodeId uint32, seqNum uint32) bool {
	_, contained := this.getSeqNumsByNodeId(nodeId)[seqNum]
	return contained
}

func (this *MessagesAlreadyDelivered) getSeqNumsByNodeId(nodeId uint32) map[uint32]uint32 {
	if value, contained := this.seqNumsAlreadyDeliveredByNodeId[nodeId]; contained {
		return value
	} else {
		this.seqNumsAlreadyDeliveredByNodeId[nodeId] = make(map[uint32]uint32)
		return this.seqNumsAlreadyDeliveredByNodeId[nodeId]
	}
}
