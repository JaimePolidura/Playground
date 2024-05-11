#include "gc_mark.hpp"

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
    std::queue<Types::Object *> pending;
    pending.push(object);

    while(!pending.empty()){
        Types::Object * currentObject = pending.front();
        pending.pop();

        if(currentObject != nullptr && !isMarked(currentObject)){
            markObject(currentObject);

            switch (currentObject->type) {
                case Types::ObjectType::ARRAY: {
                    traverseArray(reinterpret_cast<Types::ArrayObject *>(currentObject), pending);
                }
                case Types::ObjectType::STRUCT: {
                    traverseStruct(reinterpret_cast<Types::StructObject *>(currentObject), pending);
                }
                default:
                    break;
            }
        }
    }
}

void Memory::MarkCompact::Marker::traverseStruct(Types::StructObject * structObject, std::queue<Types::Object *>& pending) {
    for(auto currentField = structObject->fields;
        currentField < (structObject->fields + structObject->n_fields);
        currentField++) {

        pending.push(currentField);
    }
}

void Memory::MarkCompact::Marker::traverseArray(Types::ArrayObject * arrayObject, std::queue<Types::Object *>& pending) {
    for(auto currentElement = arrayObject->content;
        currentElement < (arrayObject->content + arrayObject->size);
        currentElement++) {

        pending.push(currentElement);
    }
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