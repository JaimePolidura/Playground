package markcompact

import (
	"gc/src"
	"gc/src/objects"
	"gc/src/vm"
	"unsafe"
)

func mark() {
	currentVm := vm.Current

	for _, thread := range currentVm.Threads {
		markRegisters(thread)
		markStack(thread)
	}

	for _, currentPackage := range currentVm.Packages {
		markPackageGlobals(currentPackage)
	}
}

func markPackageGlobals(currentPackage *vm.Package) {
	for _, global := range currentPackage.Globals {
		traverseObject(global)
	}
}

func markStack(thread *vm.Thread) {
	for i := 0; i < int(thread.StackTop); i++ {
		traverseObject(thread.Stack[i])
	}
}

func markRegisters(thread *vm.Thread) {
	for i := 0; i < len(thread.Registers); i++ {
		if thread.Registers[i] != nil {
			traverseObject(thread.Registers[i])
		}
	}
}

func traverseObject(object *objects.Object) {
	if !isMarked(object) {
		return
	}

	markObject(object)
	queue := src.Queue{}
	queue.Enqueue(object)

	for !queue.IsEmpty() {
		currentObject := queue.Dequeue().(*objects.Object)

		if !isMarked(currentObject) {
			markObject(currentObject)

			switch currentObject.Type {
			case uint8(objects.String):
				markObject(currentObject)
				break
			case uint8(objects.Array):
				arrayObject := (*objects.ArrayObject)(unsafe.Pointer(currentObject))
				if arrayObject.ContentType == objects.Primitive {
					break
				}

				for _, arrayElement := range arrayObject.Content {
					queue.Enqueue(arrayElement)
				}
				break
			case uint8(objects.Struct):
				structObject := (*objects.StructObject)(unsafe.Pointer(currentObject))
				for _, arrayElement := range structObject.Fields {
					queue.Enqueue(arrayElement)
				}

				break
			}
		}
	}
}

func markObject(object *objects.Object) {
}

func isMarked(object *objects.Object) bool {
	return true
}
