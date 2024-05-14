#include "gc_vm_info.h"

void Memory::MarkCompact::VMInfo::resetAllMarkBits() {
    for (const auto& [ignored, allocBuffer]: this->allocationBuffers) {
        allocBuffer->resetMarkBit();
    }
}

//TODO Race conditions
Memory::MarkCompact::AllocationBuffer * Memory::MarkCompact::VMInfo::submitAllocationBuffer() {
    if(this->freeAllocationBuffers.empty()){
        return new Memory::MarkCompact::AllocationBuffer();
    }

    auto allocationBuffer = this->freeAllocationBuffers.back();
    this->freeAllocationBuffers.pop_back();

    return allocationBuffer;
}