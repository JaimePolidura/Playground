#pragma once

#include "shared.hpp"
#include "types/types.hpp"

// 2 MB
#define BUFFER_ADDRESS_BIT_OFFSET 24
#define MARK_COMPACT_ALLOCATION_BUFFER_SIZE (1 << BUFFER_ADDRESS_BIT_OFFSET)

#define TO_ABSOLUTE_ALLOCATION_BUFFER_ADDR(relative) (absoluteAllocBufAddress_t) (relative >> BUFFER_ADDRESS_BIT_OFFSET)

//Full address, which points directly to the allocated objet into the buffer. For example 0x90A9125B12AAFDF
typedef uint64_t relativeAllocBufAddress_t;
//Address of allocation buffer allocated. For example relative address: 0x90A9125B12AAFDF, absolute 0x0000090A9125B1
typedef uint64_t absoluteAllocBufAddress_t;

namespace Memory::MarkCompact {
    class AllocationBuffer {
    public:
        void * allocateSize(size_t size);

        void mark(Types::Object * object);

        void unmark(Types::Object * object);

        bool isMarked(Types::Object * object);

        void resetMarkBit();

        bool belongs(Types::Object * object);

        bool hasRoom(std::size_t size);

        AllocationBuffer * getLast();

        int getNAllocationBuffersFromLast();

    private:
        //Returns {index on markBitMap, offset in the byte}
        std::pair<int, int> getBitMapIndex(Types::Object * object);

        std::byte buffer[MARK_COMPACT_ALLOCATION_BUFFER_SIZE];
        AllocationBuffer * prev{nullptr};
        AllocationBuffer * next{nullptr};
        uint64_t nextFree{0};
        std::byte markBitMap[MARK_COMPACT_ALLOCATION_BUFFER_SIZE / (sizeof(Types::Object) * 8)];

        friend class MarkCompactAllocator;
        friend class Compacter;
    };
}