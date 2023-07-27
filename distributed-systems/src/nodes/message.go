package nodes

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// Skips order, ack. Used for heartbeats
const FLAG_URGENT = 1

// The follower will communicate to the other followers without pass through the leader
// Used for when leader has failed
const FLAG_BYPASS_LEADER = 2
const MESSAGE_BROADCAST = 0
const MESSAGE_DO_BROADCAST = 8

type Message struct {
	Opts
}

func (this *Message) Clone() *Message {
	return &Message{Opts{
		NodeIdOrigin: this.NodeIdOrigin,
		NodeIdSender: this.NodeIdSender,
		SeqNum:       this.SeqNum,
		TTL:          this.TTL,
		Type:         this.Type,
		Flags:        this.Flags,
		Content:      this.Content,
	}}
}

func (this *Message) GetSizeInBytes() uint32 {
	return 4 + 4 + 4 + 4 + 1 + 1 + 1 + uint32(len(this.Content))
}

func (this *Message) GetMessageId() uint64 {
	return uint64(this.NodeIdOrigin)<<32 | uint64(this.SeqNum)
}

func (this *Message) AddFlag(flag uint8) *Message {
	this.Flags |= flag
	return this
}

func (this *Message) WithType(typeToSet uint8) *Message {
	this.Type = typeToSet
	return this
}

func (this *Message) HasFlag(flag uint8) bool {
	return this.Flags&flag != 0
}

func (this *Message) HasNotFlag(flag uint8) bool {
	return this.Flags&flag == 0
}

func (this *Message) IsType(Type uint8) bool {
	return this.Type&Type != 0
}

func (this *Message) SetContentUin32(newContent uint32) {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, newContent)
	this.Content = buf.Bytes()
}

func (this *Message) GetContentToUint32() uint32 {
	return binary.BigEndian.Uint32(this.Content)
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
	binary.Write(&buf, binary.BigEndian, message.Type)
	binary.Write(&buf, binary.BigEndian, message.Flags)
	binary.Write(&buf, binary.BigEndian, uint8(len(message.Content)))
	binary.Write(&buf, binary.BigEndian, message.Content)

	return buf.Bytes()
}

func sizeToBytes(size uint32) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, size)
	a := buf.Bytes()
	return a
}

func Deserialize(bytes []byte, start uint32) (_message *Message, _endInclusive uint32, _error error) {
	if len(bytes) < 5 {
		return nil, 0, errors.New("invalid raw message format")
	}

	NodeIdOrigin := binary.BigEndian.Uint32(bytes[start:])
	NodeIdSender := binary.BigEndian.Uint32(bytes[start+4:])
	SeqNum := binary.BigEndian.Uint32(bytes[start+8:])
	TTL := int32(binary.BigEndian.Uint32(bytes[start+12:]))
	Type := bytes[start+12+4]
	Flags := bytes[start+17]
	ContentSize := bytes[start+18]
	Content := bytes[start+19 : uint32(start)+19+uint32(ContentSize)]

	message := &Message{Opts{
		NodeIdOrigin: NodeIdOrigin,
		NodeIdSender: NodeIdSender,
		Content:      Content,
		SeqNum:       SeqNum,
		Flags:        Flags,
		Type:         Type,
		TTL:          TTL,
	}}

	return message, message.GetSizeInBytes() + start, nil
}

func WithSenderNodeId(nodeId uint32) OptFunc {
	return func(opts *Opts) {
		opts.NodeIdSender = nodeId
	}
}

func WithSeqNum(seqNum uint32) OptFunc {
	return func(opts *Opts) {
		opts.SeqNum = seqNum
	}
}

func WithOrigin(nodeId uint32) OptFunc {
	return func(opts *Opts) {
		opts.NodeIdSender = nodeId
	}
}

func WithContentString(content string) OptFunc {
	return func(opts *Opts) {
		opts.Content = []byte(content)
	}
}

func WithContentBytes(content []byte) OptFunc {
	return func(opts *Opts) {
		opts.Content = content
	}
}

func WithFlags(flags ...uint8) OptFunc {
	return func(opts *Opts) {
		for _, flag := range flags {
			opts.Flags |= flag
		}
	}
}

func WithType(Type uint8) OptFunc {
	return func(opts *Opts) {
		opts.Type = Type
	}
}

func WithTTL(TTL int32) OptFunc {
	return func(opts *Opts) {
		opts.TTL = TTL
	}
}

func WithNodeId(nodeId uint32) OptFunc {
	return func(opts *Opts) {
		opts.NodeIdOrigin = nodeId
		opts.NodeIdSender = nodeId
	}
}

func GetSizeAllInBytes(messages []*Message) uint32 {
	totalSize := uint32(0)
	for _, message := range messages {
		totalSize += message.GetSizeInBytes()
	}

	return totalSize
}

type OptFunc func(*Opts)

func CreateMessage(opts ...OptFunc) *Message {
	o := defaultOpts()
	for _, fn := range opts {
		fn(&o)
	}

	return &Message{o}
}

func defaultOpts() Opts {
	return Opts{
		Content: []byte{},
	}
}

type Opts struct {
	NodeIdOrigin uint32
	NodeIdSender uint32
	SeqNum       uint32
	TTL          int32
	Type         uint8
	Flags        uint8
	Content      []byte //Content size 1 byte
}
