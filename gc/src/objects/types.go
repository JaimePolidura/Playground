package objects

type PrimitiveType uint8
type ObjectType uint8

const (
	i8  PrimitiveType = 8
	i16               = 16
	i32               = 32
	i64               = 64
)

const (
	String ObjectType = iota
	Array
	Struct
	Primitive
)

type Object struct {
	Type uint8
}
