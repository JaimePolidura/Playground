package code

import (
	"errors"
)

type LinkedList[T comparable] struct {
	first * linkedListNode[T]
	size uint32
}

func(linkedlist * LinkedList[T]) Size () uint32 {
	return linkedlist.size
}

func(linkedlist * LinkedList[T]) Contains (value T) bool {
	return linkedlist.getNodeByData(value) != nil
}

func (linkedlist * LinkedList[T]) Stream() * Stream[T]{
	return &Stream[T]{iterator: linkedlist.Iterate()}
}

func (linkedlist * LinkedList[T]) Clear() {
	linkedlist.first = nil
	linkedlist.size = 0
}

func(linkedlist * LinkedList[T]) Get (targetIndex uint32) (T, error) {
	var defaultIfNotFound T
	var actualIndex uint32 = 0
	node := linkedlist.first

	for node != nil {
		if actualIndex == targetIndex {
			return node.data, errors.New("not found")
		}
		if actualIndex > targetIndex {
			return defaultIfNotFound, errors.New("not found")
		}

		node = node.next
		actualIndex++
	}

	return defaultIfNotFound, nil
}

func(linkedlist * LinkedList[T]) Remove (value T) bool {
	nodeToRemove := linkedlist.getNodeByData(value)

	if nodeToRemove == nil {
		return false
	}

	backNodeToRemove := nodeToRemove.back
	nextNodeToRemove := nodeToRemove.next

	if backNodeToRemove != nil && nextNodeToRemove != nil {
		backNodeToRemove.next = nextNodeToRemove
		nextNodeToRemove.back = backNodeToRemove
	}else if backNodeToRemove != nil {
		backNodeToRemove.next = nil
	}else if nextNodeToRemove != nil {
		nextNodeToRemove.back = nil
	}

	if nodeToRemove == linkedlist.first {
		linkedlist.first = nextNodeToRemove
	}

	linkedlist.size = linkedlist.size - 1

	return true
}

func(linkedlist * LinkedList[T]) Add (value T) uint32 {
	newNode := new(linkedListNode[T])
	newNode.data = value

	if linkedlist.first == nil {
		linkedlist.first = newNode

	} else {
		lastNode := linkedlist.getLastNode()
		newNode.back = lastNode
		lastNode.next = newNode
	}

	linkedlist.size = linkedlist.size + 1

	return linkedlist.size
}

func(linkedlist * LinkedList[T]) getNodeByData(data T) * linkedListNode[T] {
	node := linkedlist.first

	for node != nil {
		if node.data == data {
			return node
		}

		node = node.next
	}

	return nil
}

func(linkedlist * LinkedList[T]) getLastNode() * linkedListNode[T] {
	node := linkedlist.first

	for node != nil && node.next != nil {
		node = node.next
	}

	return node
}

func(linkedlist * LinkedList[T]) Iterate() Iterator[T] {
	return &LinkedListIterator[T]{linkedlist.first}
}

type LinkedListIterator[T comparable] struct {
	actualNode * linkedListNode[T]
}

func(iterator * LinkedListIterator[T]) HastNext() bool {
	return iterator.actualNode != nil
}

func(iterator * LinkedListIterator[T]) Next() T {
	data := iterator.actualNode.data
	iterator.actualNode = iterator.actualNode.next

	return data
}

type linkedListNode[T comparable] struct {
	next * linkedListNode[T]
	back * linkedListNode[T]
	data T
}
