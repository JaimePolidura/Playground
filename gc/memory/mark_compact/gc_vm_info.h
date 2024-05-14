#pragma once

#include "memory/mark_compact/allocation_buffer.hpp"
#include "shared.hpp"

#define GET_VM_GC_INFO(vm) (reinterpret_cast<Memory::MarkCompact::VMInfo *>(vm.gc))
#define GET_ALLOCATION_BUFFER(current_vm, anyAddress) (GET_VM_GC_INFO(current_vm)->allocationBuffers[TO_ABSOLUTE_ALLOCATION_BUFFER_ADDR((uint64_t) anyAddress)])

namespace Memory::MarkCompact {
    struct VMInfo {
        std::map<absoluteAllocBufAddress_t, Memory::MarkCompact::AllocationBuffer *> allocationBuffers{};

        std::vector<Memory::MarkCompact::AllocationBuffer *> freeAllocationBuffers;

        Memory::MarkCompact::AllocationBuffer * submitAllocationBuffer();
        void resetAllMarkBits();
    };
}