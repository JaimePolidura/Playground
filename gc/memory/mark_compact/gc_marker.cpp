#include "gc_marker.hpp"

extern struct VM::VM current_vm;

void Memory::MarkCompact::Marker::mark() {
    current_vm.stopThreadsGC();
    markThreadsStack();
    markPackages();
}

void Memory::MarkCompact::Marker::markThreadsStack() {
    for(VM::Thread& thread : current_vm.threads){
        markStack(thread);
    }
}

void Memory::MarkCompact::Marker::markPackages() {
    for (const auto& [packageName, package] : current_vm.packages) {
        markPackage(package);
    }
}

void Memory::MarkCompact::Marker::markStack(VM::Thread& thread) {
    for(int i = 0; i < thread.esp; i++){
        traverseObject(thread.stack[i]);
    }
}

void Memory::MarkCompact::Marker::markPackage(std::shared_ptr<VM::Package> package) {
    for (const auto& [globalName, global]: package->globals) {
        traverseObject(global);
    }
}

void Memory::MarkCompact::Marker::traverseObject(Types::Object * object) {
    Types::traverseObjectDeep(object, [this](Types::Object * currentObject) -> bool {
        if(currentObject != nullptr && !this->isMarked(currentObject)){
            this->markObject(currentObject);
            return true;
        } else {
            return false;
        }
    });
}

void Memory::MarkCompact::Marker::markObject(Types::Object * object) {
    auto allocationBuffer = this->getAllocationBuffer(TO_ABSOLUTE_ALLOCATION_BUFFER_ADDR((uint64_t) object));
    allocationBuffer->mark(object);
}

bool Memory::MarkCompact::Marker::isMarked(Types::Object * object) {
    auto allocationBuffer = this->getAllocationBuffer(TO_ABSOLUTE_ALLOCATION_BUFFER_ADDR((uint64_t) object));
    return allocationBuffer->isMarked(object);
}

Memory::MarkCompact::AllocationBuffer * Memory::MarkCompact::Marker::getAllocationBuffer(absoluteAllocBufAddress_t addressToLookup) {
    auto allocationBufferAtAddress = this->allocationsBufferByAddress.find(addressToLookup);
    if(allocationBufferAtAddress != this->allocationsBufferByAddress.end()){
        return allocationBufferAtAddress->second;
    }

    for (const VM::Thread& currentThread: current_vm.threads) {
        auto gcThreadInfo = reinterpret_cast<Memory::MarkCompact::ThreadInfo *>(currentThread.gc);

        if(gcThreadInfo->allocationBuffers.contains(addressToLookup)){
            auto allocationBufferResult = gcThreadInfo->allocationBuffers[addressToLookup];
            this->allocationsBufferByAddress.insert({addressToLookup, allocationBufferResult});
            return allocationBufferResult;
        }
    }

    return nullptr;
}