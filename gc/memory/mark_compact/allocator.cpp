#include "allocator.hpp"

extern thread_local VM::Thread * self_thread;
extern struct VM::VM current_vm;

Types::StringObject * Memory::MarkCompact::MarkCompactAllocator::allocString(char * data) {
    auto stringObject = reinterpret_cast<Types::StringObject *>(this->allocateSize(Types::StringObject::size(data)));
    stringObject->object.type = Types::ObjectType::STRING;
    stringObject->nChars = strlen(data);
    std::memcpy(stringObject->contents, data, strlen(data));
    return stringObject;
}

Types::StructObject * Memory::MarkCompact::MarkCompactAllocator::allocStruct(int nFields) {
    auto structObject = reinterpret_cast<Types::StructObject *>(this->allocateSize(Types::StructObject::size(nFields)));
    structObject->object.type = Types::ObjectType::STRUCT;
    structObject->nFields = nFields;
    return structObject;
}

Types::ArrayObject * Memory::MarkCompact::MarkCompactAllocator::allocArray(int nElements) {
    auto arrayObject = reinterpret_cast<Types::ArrayObject *>(this->allocateSize(Types::ArrayObject::size(nElements)));
    arrayObject->nElements = nElements;
    arrayObject->object.type = Types::ARRAY;
    return arrayObject;
}

void * Memory::MarkCompact::MarkCompactAllocator::allocateSize(size_t size) {
    auto gcThreadInfo = reinterpret_cast<Memory::MarkCompact::ThreadInfo *>(self_thread->gc);
    auto currentAllocationBuffer = gcThreadInfo->allocationBuffer;

    void * ptr = currentAllocationBuffer->allocateSize(size);
    if(ptr != nullptr){
        return ptr;
    }

    auto newAllocBuf = nextAllocBuffer(size);

    return newAllocBuf->allocateSize(size);
}

Memory::MarkCompact::AllocationBuffer * Memory::MarkCompact::MarkCompactAllocator::nextAllocBuffer(size_t size) {
    auto gcThreadInfo = reinterpret_cast<Memory::MarkCompact::ThreadInfo *>(self_thread->gc);
    auto currentAllocationBuffer = gcThreadInfo->allocationBuffer;

    while(currentAllocationBuffer->next != nullptr && !currentAllocationBuffer->hasRoom(size)){
        currentAllocationBuffer = currentAllocationBuffer->next;
    }

    if(currentAllocationBuffer->hasRoom(size)) {
        gcThreadInfo->allocationBuffer = currentAllocationBuffer;
        return currentAllocationBuffer;
    }

    if(gcThreadInfo->nAllocatedBuffersWithoutGC + 1 < gcThreadInfo->nextGcAllocatedBuffer) {
        auto gcVmInfo = GET_VM_GC_INFO(current_vm);
        auto newAllocationBuffer = gcVmInfo->submitAllocationBuffer();;
        gcVmInfo->allocationBuffers[TO_ABSOLUTE_ALLOCATION_BUFFER_ADDR((uint64_t) newAllocationBuffer->buffer[0])] = newAllocationBuffer;
        newAllocationBuffer->prev = currentAllocationBuffer;
        currentAllocationBuffer->next = newAllocationBuffer;
        gcThreadInfo->allocationBuffer = newAllocationBuffer;
        gcThreadInfo->nAllocatedBuffersWithoutGC++;
        return newAllocationBuffer;
    }

    startGC();
}

void Memory::MarkCompact::MarkCompactAllocator::startGC() {
    Memory::MarkCompact::Marker marker{};
    marker.mark();

    Memory::MarkCompact::Compacter compacter{};
    compacter.compact();

    updateThreadAllocationBuffers();
}

void Memory::MarkCompact::MarkCompactAllocator::updateThreadAllocationBuffers() {
    for (const auto &thread: current_vm.threads) {
        auto gcThreadInfo = reinterpret_cast<Memory::MarkCompact::ThreadInfo *>(thread.gc);
        if(gcThreadInfo->allocationBuffer != nullptr) {
            auto newPointingAllocationBuffer = gcThreadInfo->allocationBuffer->getLast();
            gcThreadInfo->nAllocatedBuffersWithoutGC =  newPointingAllocationBuffer->getNAllocationBuffersFromLast();
            gcThreadInfo->allocationBuffer = newPointingAllocationBuffer;
        }
    }

    auto selfThreadGcInfo = reinterpret_cast<Memory::MarkCompact::ThreadInfo *>(self_thread->gc);
    selfThreadGcInfo->nextGcAllocatedBuffer = selfThreadGcInfo->nAllocatedBuffersWithoutGC * 2;
}