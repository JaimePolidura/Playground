#include "allocator.hpp"

extern thread_local VM::Thread * self_thread;

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

    if(gcThreadInfo->nAllocatedBuffersWithoutGC + 1 < gcThreadInfo->nextGcAllocatedBuffer) {
        auto newAllocationBuffer = new Memory::MarkCompact::AllocationBuffer();
        gcThreadInfo->allocationBuffers[TO_ABSOLUTE_ALLOCATION_BUFFER_ADDR((uint64_t) newAllocationBuffer)] = newAllocationBuffer;
        newAllocationBuffer->setPrev(currentAllocationBuffer);
        gcThreadInfo->allocationBuffer = newAllocationBuffer;
        gcThreadInfo->nAllocatedBuffersWithoutGC++;

        return (void *) currentAllocationBuffer->allocateSize(size);
    }

    startGC();

    return nullptr;
}

void Memory::MarkCompact::MarkCompactAllocator::startGC() {
    Memory::MarkCompact::Marker marker{};
    marker.mark();

    Memory::MarkCompact::Compacter compacter{};
    compacter.compact();
}