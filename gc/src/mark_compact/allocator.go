package markcompact

import (
	"gc/src/objects"
	"gc/src/vm"
	"unsafe"
)

func AllocArray(self *vm.Thread, primitiveSize objects.PrimitiveType, nElements int32) *objects.ArrayObject {
	sizeToAllocate := int32(primitiveSize)*nElements + int32(unsafe.Sizeof(objects.ArrayObject{}))
	arrayObject := allocate[objects.ArrayObject](sizeToAllocate)
	for index, _ := range arrayObject.Content {
		arrayObject.Content[index] = 0
	}
	return arrayObject
}

func AllocString(self *vm.Thread, length int, values []byte) *objects.StringObject {
	sizeToAllocate := int32(len(values)) + int32(unsafe.Sizeof(objects.StringObject{}))
	stringObject := allocate[objects.StringObject](sizeToAllocate)
	stringObject.Content = values
	return stringObject
}

func AllocStruct(self *vm.Thread, nFields int) *objects.StructObject {
	sizeToAllocate := int32(8*nFields) + int32(unsafe.Sizeof(objects.StructObject{}))
	structObject := allocate[objects.StructObject](sizeToAllocate)
	for index, _ := range structObject.Fields {
		structObject.Fields[index] = nil
	}
	return structObject
}

func allocate[T any](sizeToAllocate int32) *T {
	markCompactGc := vm.Current.GC.(*MarkCompactGC)

	relativePtr := markCompactGc.TryAllocate(sizeToAllocate)

	if relativePtr == 0 {
		StartGC()
		relativePtr = markCompactGc.TryAllocate(sizeToAllocate)
	}

	return (*T)(unsafe.Pointer(unsafe.Offsetof(MarkCompactGC{}.Data) + uintptr(relativePtr)))
}
