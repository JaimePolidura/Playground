package broadcast

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type Message struct {
	NodeIdOrigin uint32
	NodeIdSender uint32
	SeqNum       uint32
	TTL          int32
	Content      []byte //Content size 1 byte
}

func (message *Message) GetSizeInBytes() uint32 {
	return 4 + 4 + 4 + 4 + 1 + uint32(len(message.Content))
}

func (message *Message) GetMessageId() uint64 {
	return uint64(message.NodeIdOrigin)<<32 | uint64(message.SeqNum)
}

func CreateMessage(nodeIdOrigin uint32, nodeIdSender uint32, content string) *Message {
	return &Message{
		NodeIdOrigin: nodeIdOrigin,
		NodeIdSender: nodeIdSender,
		Content:      []byte(content),
	}
}

func SerializeAll(messages []*Message) []byte {
	sizeBytes := sizeToBytes(GetSizeAllInBytes(messages))

	contentBytes := make([]byte, 0)
	for _, message := range messages {
		contentBytes = append(contentBytes, serializeNotIncludingSize(message)...)
	}

	return append(sizeBytes, contentBytes...)
}

func Serialize(message *Message) []byte {
	sizeBytes := sizeToBytes(message.GetSizeInBytes())
	contentBytes := serializeNotIncludingSize(message)

	return append(sizeBytes, contentBytes...)
}

func serializeNotIncludingSize(message *Message) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, message.NodeIdOrigin)
	binary.Write(&buf, binary.BigEndian, message.NodeIdSender)
	binary.Write(&buf, binary.BigEndian, message.SeqNum)
	binary.Write(&buf, binary.BigEndian, message.TTL)
	binary.Write(&buf, binary.BigEndian, uint8(len(message.Content)))
	binary.Write(&buf, binary.BigEndian, message.Content)

	return buf.Bytes()
}

func sizeToBytes(size uint32) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, size)
	return buf.Bytes()
}

func Deserialize(bytes []byte, start uint32) (_message *Message, _endInclusive uint32, _error error) {
	if len(bytes) < 5 {
		return nil, 0, errors.New("invalid raw message format")
	}

	NodeIdOrigin := binary.BigEndian.Uint32(bytes[start:])
	NodeIdSender := binary.BigEndian.Uint32(bytes[start+4:])
	SeqNum := binary.BigEndian.Uint32(bytes[start+8:])
	TTL := int32(binary.BigEndian.Uint32(bytes[start+12:]))
	ContentSize := bytes[start+16]
	Content := bytes[start+17 : uint32(start)+17+uint32(ContentSize)]

	message := &Message{
		NodeIdOrigin: NodeIdOrigin,
		NodeIdSender: NodeIdSender,
		SeqNum:       SeqNum,
		TTL:          TTL,
		Content:      Content,
	}

	return message, message.GetSizeInBytes() + start, nil
}

func GetSizeAllInBytes(messages []*Message) uint32 {
	totalSize := uint32(0)
	for _, message := range messages {
		totalSize += message.GetSizeInBytes()
	}

	return totalSize
}
