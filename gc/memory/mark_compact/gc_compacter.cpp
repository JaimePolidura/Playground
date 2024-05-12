#include "gc_compacter.hpp"

extern struct VM::VM current_vm;

void Memory::MarkCompact::Compacter::compact() {
    compactThreads();
    updateReferences();
}

void Memory::MarkCompact::Compacter::compactThreads() {
    for (const VM::Thread& thread: current_vm.threads) {
        compactThread(thread);
    }
}

void Memory::MarkCompact::Compacter::compactThread(const VM::Thread& thread) {
    auto gcThreadInfo = reinterpret_cast<Memory::MarkCompact::ThreadInfo *>(thread.gc);
    auto currentAllocBuffer = gcThreadInfo->allocationBuffer;

    Types::Object * free = reinterpret_cast<Types::Object *>(currentAllocBuffer->buffer);
    Types::Object * live = reinterpret_cast<Types::Object *>(
            currentAllocBuffer->buffer + MARK_COMPACT_ALLOCATION_BUFFER_SIZE - sizeof(Types::Object)
    );

    while(free < live) {
        live = this->nextLive(live, currentAllocBuffer);
        free = this->nextFree(free, currentAllocBuffer);

        if(live > free) {
            std::size_t freeSize = this->sizeFree(free, currentAllocBuffer);
            std::size_t liveSize = Types::sizeofObject(live->type);

            if(freeSize >= liveSize){
                this->moveObjects(free, live, currentAllocBuffer);
            } else {
                live -= sizeof(Types::Object);
            }
        }
    }
}

void Memory::MarkCompact::Compacter::updateReferences() {
    updateGlobalReferences();
    
    for (const VM::Thread& thread: current_vm.threads) {
        for(int i = 0; i < thread.esp; i++){
        }
    }
}

void Memory::MarkCompact::Compacter::updateGlobalReferences() {
    std::vector<std::pair<std::string, Types::Object *>> toUpdate{};

    for (const auto& [packageName, package]: current_vm.packages) {
        for (const auto& it: package->globals) {
            if(it.second->gc != nullptr) {
                toUpdate.push_back(it);
            }
        }

        for (const auto& [globalName, globalValue]: toUpdate) {
            package->globals[globalName] = globalValue;
        }
    }
}


Types::Object * Memory::MarkCompact::Compacter::nextLive(Types::Object * prevPtr, AllocationBuffer * currentAllocBuffer) {
    while(!currentAllocBuffer->isMarked(prevPtr)) {
        prevPtr--;
    }

    return prevPtr;
}

Types::Object * Memory::MarkCompact::Compacter::nextFree(Types::Object * prevPtr, AllocationBuffer * currentAllocBuffer) {
    while(currentAllocBuffer->isMarked(prevPtr)) {
        prevPtr += Types::sizeofObject(prevPtr->type);
    }

    return prevPtr;
}

std::size_t Memory::MarkCompact::Compacter::sizeFree(Types::Object * freePtr, AllocationBuffer * currentAllocBuffer) {
    std::size_t size = 0;

    while(!currentAllocBuffer->isMarked(freePtr)){
        size += sizeof(Types::Object);
        freePtr++;
    }

    return size;
}

void Memory::MarkCompact::Compacter::moveObjects(Types::Object * dst, Types::Object * src, AllocationBuffer * currentAllocBuffer) {
    *dst = *src;
    currentAllocBuffer->mark(dst);
    currentAllocBuffer->unmark(src);
    src->gc = dst;
}