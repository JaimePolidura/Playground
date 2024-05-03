package vm

import "gc/src/objects"

type ThreadState uint8

const (
	Running = iota
	Death
	Waiting
)

type Thread struct {
	Registers [16]*objects.Object
	Stack     [256]*objects.Object
	StackTop  int32

	State ThreadState
	GC    *interface{}
}
