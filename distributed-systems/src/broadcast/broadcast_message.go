package broadcast

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type Message struct {
	NodeId  uint32
	SeqNum  uint32
	TTL     int32
	Content []byte //Content size 1 byte
}

func (message *Message) GetSizeInBytes() uint32 {
	return 4 + 4 + 4 + 1 + uint32(len(message.Content))
}

func (message *Message) GetMessageId() uint64 {
	return uint64(message.NodeId)<<32 | uint64(message.SeqNum)
}

func CreateMessage(nodeId uint32, content string) *Message {
	return &Message{
		NodeId:  nodeId,
		Content: []byte(content),
		TTL:     0,
		SeqNum:  0,
	}
}

func Serialize(message *Message) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, message.GetSizeInBytes())
	binary.Write(&buf, binary.BigEndian, message.NodeId)
	binary.Write(&buf, binary.BigEndian, message.SeqNum)
	binary.Write(&buf, binary.BigEndian, message.TTL)
	binary.Write(&buf, binary.BigEndian, uint8(len(message.Content)))
	binary.Write(&buf, binary.BigEndian, message.Content)

	a := buf.Bytes()
	return a
}

func Deserialize(bytes []byte, start uint32) (_message *Message, _endInclusive uint32, _error error) {
	if len(bytes) < 5 {
		return nil, 0, errors.New("invalid raw message format")
	}

	NodeId := binary.BigEndian.Uint32(bytes[start:])
	SeqNum := binary.BigEndian.Uint32(bytes[start+4:])
	TTL := int32(binary.BigEndian.Uint32(bytes[start+8:]))
	ContentSize := bytes[start+12]
	Content := bytes[start+13 : uint32(start)+13+uint32(ContentSize)]

	message := &Message{
		NodeId:  NodeId,
		SeqNum:  SeqNum,
		TTL:     TTL,
		Content: Content,
	}

	return message, message.GetSizeInBytes() + start, nil
}
