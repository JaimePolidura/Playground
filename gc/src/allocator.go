package src

import (
	"gc/src/objects"
	"gc/src/vm"
)

type Allocator interface {
	AllocArray(self *vm.Thread, primitiveSize objects.PrimitiveType, nElements int32) *objects.ArrayObject
	AllocString(self *vm.Thread, values []byte) *objects.StringObject
	AllocStruct(self *vm.Thread, nFields int) *objects.StructObject

	ForceGc()
}
