#include "allocation_buffer.hpp"

void * Memory::MarkCompact::AllocationBuffer::allocateSize(size_t size) {
    if(this->nextFree + size >= MARK_COMPACT_ALLOCATION_BUFFER_SIZE) {
        return nullptr;
    }
    if(size >= MARK_COMPACT_ALLOCATION_BUFFER_SIZE){
        return nullptr;
    }

    void * ptr = &this->buffer[this->nextFree];
    this->nextFree += size;

    return ptr;
}

void Memory::MarkCompact::AllocationBuffer::resetMarkBit() {
    memset(this->markBitMap, 0, sizeof(this->markBitMap));
}

void Memory::MarkCompact::AllocationBuffer::setPrev(Memory::MarkCompact::AllocationBuffer * other) {
    this->prev = other;
}

void Memory::MarkCompact::AllocationBuffer::mark(Types::Object * object) {
    const auto [markBitMapIndex, offsetInByteBitMap] = this->getBitMapIndex(object);
    this->markBitMap[markBitMapIndex] |= static_cast<std::byte>(1 << offsetInByteBitMap);
}

void Memory::MarkCompact::AllocationBuffer::unmark(Types::Object * object) {
    const auto [markBitMapIndex, offsetInByteBitMap] = this->getBitMapIndex(object);
    this->markBitMap[markBitMapIndex] ^= static_cast<std::byte>(1 << offsetInByteBitMap);
}

bool Memory::MarkCompact::AllocationBuffer::isMarked(Types::Object * object) {
    const auto [markBitMapIndex, offsetInByteBitMap] = this->getBitMapIndex(object);
    std::byte valueInByte = this->markBitMap[markBitMapIndex];
    std::byte valueInBit = (valueInByte >> offsetInByteBitMap) & static_cast<std::byte>(0x01);

    return static_cast<bool>(valueInBit);
}

std::pair<int, int> Memory::MarkCompact::AllocationBuffer::getBitMapIndex(Types::Object * object) {
    uint64_t allocBufferIndex = (((uint64_t) object) << (sizeof(uint64_t) - BUFFER_ADDRESS_BIT_OFFSET))
            >> (sizeof(uint64_t) - BUFFER_ADDRESS_BIT_OFFSET);
    int markBitMapIndex = static_cast<int>(allocBufferIndex / (sizeof(Types::Object) * 8));
    int offsetInByteBitMap = markBitMapIndex - roundLessTo8(markBitMapIndex);

    return std::make_pair(markBitMapIndex, offsetInByteBitMap);
}