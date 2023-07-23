package broadcast

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type Message struct {
	NodeId  uint32
	SeqNum  uint32
	TTL     uint32
	Content []byte //Content size 1 byte
}

func Serialize(message *Message) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, message.NodeId)
	binary.Write(&buf, binary.BigEndian, message.SeqNum)
	binary.Write(&buf, binary.BigEndian, message.TTL)
	binary.Write(&buf, binary.BigEndian, uint8(len(message.Content)))
	binary.Write(&buf, binary.BigEndian, message.Content)

	return buf.Bytes()
}

func Deserialize(bytes []byte) (*Message, error) {
	if len(bytes) < 5 {
		return nil, errors.New("invalid raw message format")
	}

	NodeId := binary.BigEndian.Uint32(bytes)
	SeqNum := binary.BigEndian.Uint32(bytes[4:])
	TTL := binary.BigEndian.Uint32(bytes[8:])
	ContentSize := bytes[12]
	Content := bytes[13:ContentSize]

	return &Message{
		NodeId:  NodeId,
		SeqNum:  SeqNum,
		TTL:     TTL,
		Content: Content,
	}, nil
}
