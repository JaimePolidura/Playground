package broadcast

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type BroadcastMessage struct {
	NodeIdOrigin uint32
	NodeIdSender uint32
	SeqNum       uint32
	TTL          int32
	Flags        uint8
	Content      []byte //Content size 1 byte
}

func (this *BroadcastMessage) GetSizeInBytes() uint32 {
	return 4 + 4 + 4 + 4 + 1 + 1 + uint32(len(this.Content))
}

func (this *BroadcastMessage) GetMessageId() uint64 {
	return uint64(this.NodeIdOrigin)<<32 | uint64(this.SeqNum)
}

func (this *BroadcastMessage) SetFlag(flag uint8) *BroadcastMessage {
	this.Flags |= flag
	return this
}

func (this *BroadcastMessage) HasFlag(flag uint8) bool {
	return this.Flags&flag != 0
}

func (this *BroadcastMessage) SetContentUin32(newContent uint32) {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, newContent)
	this.Content = buf.Bytes()
}

func (this *BroadcastMessage) GetContentToUint32() uint32 {
	return binary.BigEndian.Uint32(this.Content)
}

func CreateMessageWithFlags(nodeIdOrigin uint32, nodeIdSender uint32, content string, flags uint8) *BroadcastMessage {
	return &BroadcastMessage{
		NodeIdOrigin: nodeIdOrigin,
		NodeIdSender: nodeIdSender,
		Content:      []byte(content),
		Flags:        flags,
	}
}

func CreateMessage(nodeIdOrigin uint32, nodeIdSender uint32, content string) *BroadcastMessage {
	return &BroadcastMessage{
		NodeIdOrigin: nodeIdOrigin,
		NodeIdSender: nodeIdSender,
		Content:      []byte(content),
	}
}

func SerializeAll(messages []*BroadcastMessage) []byte {
	sizeBytes := sizeToBytes(GetSizeAllInBytes(messages))

	contentBytes := make([]byte, 0)
	for _, message := range messages {
		contentBytes = append(contentBytes, serializeNotIncludingSize(message)...)
	}

	return append(sizeBytes, contentBytes...)
}

func Serialize(message *BroadcastMessage) []byte {
	sizeBytes := sizeToBytes(message.GetSizeInBytes())
	contentBytes := serializeNotIncludingSize(message)

	return append(sizeBytes, contentBytes...)
}

func serializeNotIncludingSize(message *BroadcastMessage) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, message.NodeIdOrigin)
	binary.Write(&buf, binary.BigEndian, message.NodeIdSender)
	binary.Write(&buf, binary.BigEndian, message.SeqNum)
	binary.Write(&buf, binary.BigEndian, message.TTL)
	binary.Write(&buf, binary.BigEndian, message.Flags)
	binary.Write(&buf, binary.BigEndian, uint8(len(message.Content)))
	binary.Write(&buf, binary.BigEndian, message.Content)

	return buf.Bytes()
}

func sizeToBytes(size uint32) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, size)
	return buf.Bytes()
}

func Deserialize(bytes []byte, start uint32) (_message *BroadcastMessage, _endInclusive uint32, _error error) {
	if len(bytes) < 5 {
		return nil, 0, errors.New("invalid raw message format")
	}

	NodeIdOrigin := binary.BigEndian.Uint32(bytes[start:])
	NodeIdSender := binary.BigEndian.Uint32(bytes[start+4:])
	SeqNum := binary.BigEndian.Uint32(bytes[start+8:])
	TTL := int32(binary.BigEndian.Uint32(bytes[start+12:]))
	Flags := bytes[start+12+4]
	ContentSize := bytes[start+17]
	Content := bytes[start+18 : uint32(start)+18+uint32(ContentSize)]

	message := &BroadcastMessage{
		NodeIdOrigin: NodeIdOrigin,
		NodeIdSender: NodeIdSender,
		SeqNum:       SeqNum,
		TTL:          TTL,
		Flags:        Flags,
		Content:      Content,
	}

	return message, message.GetSizeInBytes() + start, nil
}

func GetSizeAllInBytes(messages []*BroadcastMessage) uint32 {
	totalSize := uint32(0)
	for _, message := range messages {
		totalSize += message.GetSizeInBytes()
	}

	return totalSize
}
