#pragma once

#include "shared.hpp"
#include "memory/mark_compact/allocation_buffer.hpp"

namespace Memory::MarkCompact {
    struct ThreadInfo {
        AllocationBuffer * allocationBuffer{nullptr};
        int nAllocatedBuffersWithoutGC{0};
        int nextGcAllocatedBuffer{0};
    };
}