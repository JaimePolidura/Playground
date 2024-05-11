#include "allocator.hpp"

extern thread_local VM::Thread * self_thread;

Types::StringObject * Memory::MarkCompact::MarkCompactAllocator::allocString(char * data) {
}

Types::StructObject * Memory::MarkCompact::MarkCompactAllocator::allocStruct(int nFields, Types::Object * fields) {
}

Types::ArrayObject * Memory::MarkCompact::MarkCompactAllocator::allocArray(Types::ObjectType contentType, int nElements) {
}

void * Memory::MarkCompact::MarkCompactAllocator::allocateSize(size_t size) {
    auto gcThreadInfo = reinterpret_cast<Memory::MarkCompact::ThreadInfo *>(self_thread->gc);
    auto currentAllocationBuffer = gcThreadInfo->allocationBuffer;

    void * ptr = (void *) currentAllocationBuffer->allocateSize(size);
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

    Memory::MarkCompact::compact();
}