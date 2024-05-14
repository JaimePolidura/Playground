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
    auto currentAllocBufferFree = gcThreadInfo->allocationBuffer;
    auto currentAllocBufferScan = currentAllocBufferFree->getLast();

    Types::Object * free = reinterpret_cast<Types::Object *>(currentAllocBufferFree->buffer);
    Types::Object * scan = reinterpret_cast<Types::Object *>(
            currentAllocBufferScan->buffer + MARK_COMPACT_ALLOCATION_BUFFER_SIZE - sizeof(Types::Object)
    );

    while(free < scan) {
        scan = this->nextScan(scan, &currentAllocBufferScan);
        free = this->nextFree(free, &currentAllocBufferFree);

        if(scan > free &&
            currentAllocBufferScan > currentAllocBufferFree &&
            currentAllocBufferScan != nullptr &&
            currentAllocBufferFree != nullptr) {

            std::size_t freeSize = this->sizeFree(free, currentAllocBufferFree);
            std::size_t liveSize = Types::sizeofObject(scan);

            if (freeSize >= liveSize){
                this->moveObjects(free, scan, currentAllocBufferFree, currentAllocBufferScan);
            } else {
                scan -= sizeof(Types::Object);
            }
        }
    }
}

void Memory::MarkCompact::Compacter::updateReferences() {
    updateGlobalReferences();
    updateStackReferences();
    GET_VM_GC_INFO(current_vm)->resetAllMarkBits();

    deleteGlobalsForwardingAddress();
    deleteThreadsForwardingAddress();
    GET_VM_GC_INFO(current_vm)->resetAllMarkBits();
}

void Memory::MarkCompact::Compacter::updateStackReferences() {
    for (VM::Thread& thread: current_vm.threads) {
        for(int i = 0; i < thread.esp; i++) {
            Types::Object * currentObject = thread.stack[i];
            if(currentObject->gc != nullptr){
                thread.stack[i] = reinterpret_cast<Types::Object *>(currentObject->gc);
            }

            updateObjectReferencesDeep(thread.stack[i]);
        }
    }
}

void Memory::MarkCompact::Compacter::updateGlobalReferences() {
    std::vector<std::pair<std::string, VM::Global>> toUpdate{};

    for (const auto& [packageName, package]: current_vm.packages) {
        for (const auto& it: package->globals) {
            if(it.second.value->gc != nullptr) {
                toUpdate.push_back({it.first, it.second});
            }
        }

        for (const auto& [globalName, globalValue]: toUpdate) {
            auto movedObject = reinterpret_cast<Types::Object *>(globalValue.value->gc);
            package->globals[globalName] = {movedObject};
            this->updateObjectReferencesDeep(movedObject);
        }
    }
}

void Memory::MarkCompact::Compacter::updateObjectReferencesDeep(Types::Object * object) {
    std::queue<Types::Object *> pending{};
    pending.push(object);

    while(!pending.empty()) {
        Types::Object * currentObject = pending.front();
        pending.pop();

        auto allocationBuffer = GET_ALLOCATION_BUFFER(current_vm, currentObject);

        if(allocationBuffer->isMarked(currentObject)) {
            continue;
        }

        allocationBuffer->mark(currentObject);

        switch (currentObject->type) {
            case Types::ObjectType::ARRAY: {
                Types::ArrayObject * arrayObject = AS_ARRAY(currentObject);

                for(int i = 0; i < arrayObject->nElements; i++) {
                    Types::Object * arrayObjectItem = arrayObject->elements[i];
                    if(arrayObjectItem == nullptr) {
                        continue;
                    }
                    if(arrayObjectItem->gc != nullptr) {
                        arrayObject->elements[i] = reinterpret_cast<Types::Object *>(arrayObjectItem->gc);
                    }
                }

                break;
            }
            case Types::ObjectType::STRUCT: {
                Types::StructObject * structObject = AS_STRUCT(currentObject);
                for(int i = 0; i < structObject->nFields; i++) {
                    Types::Object * structObjectField = structObject->fields[i];
                    if(structObjectField == nullptr) {
                        continue;
                    }
                    if(structObjectField->gc != nullptr) {
                        structObject->fields[i] = reinterpret_cast<Types::Object *>(structObjectField->gc);
                    }
                }

                break;
            }
            default:
                break;
        }
    }
}

void Memory::MarkCompact::Compacter::deleteGlobalsForwardingAddress() {
    for (const auto &[packageName, package]: current_vm.packages) {
        for (const auto &[globalName, global]: package->globals) {
            deleteForwardingAddressDeep(global.value);
        }
    }
}

void Memory::MarkCompact::Compacter::deleteThreadsForwardingAddress() {
    for (const VM::Thread& currenThread: current_vm.threads) {
        for(int i = 0; i < currenThread.esp; i++){
            deleteForwardingAddressDeep(currenThread.stack[i]);
        }
    }
}

void deleteForwardingAddressDeep(Types::Object * object) {
    Types::traverseObjectDeep(object, [](Types::Object * currentObject) -> bool {
        if(GET_ALLOCATION_BUFFER(current_vm, currentObject)->isMarked(currentObject)){
            return false; //Stop
        }
        currentObject->gc = nullptr;
        return true;
    });
}

Types::Object * Memory::MarkCompact::Compacter::nextScan(Types::Object * prevPtr, AllocationBuffer ** currentAllocBufferIndirect) {
    AllocationBuffer * currentAllocBuffer = *currentAllocBufferIndirect;
    while (currentAllocBuffer != nullptr) {
        while(currentAllocBuffer->belongs(prevPtr) && !currentAllocBuffer->isMarked(prevPtr)) {
            prevPtr--;
        }
        if(currentAllocBuffer->belongs(prevPtr)){
            return prevPtr;
        }

        currentAllocBuffer = currentAllocBuffer->prev;

        if(currentAllocBuffer != nullptr){
            prevPtr = reinterpret_cast<Types::Object *>(
                    currentAllocBuffer->buffer + MARK_COMPACT_ALLOCATION_BUFFER_SIZE - sizeof(Types::Object)
            );
        }
    }

    *currentAllocBufferIndirect = currentAllocBuffer;

    return prevPtr;
}

Types::Object * Memory::MarkCompact::Compacter::nextFree(Types::Object * prevPtr, AllocationBuffer ** currentAllocBufferIndirect) {
    AllocationBuffer * currentAllocBuffer = *currentAllocBufferIndirect;
    auto vmGcInfo = GET_VM_GC_INFO(current_vm);

    while(currentAllocBuffer != nullptr){
        while(currentAllocBuffer->belongs(prevPtr) && currentAllocBuffer->isMarked(prevPtr)) {
            prevPtr += Types::sizeofObject(prevPtr);
        }

        if(currentAllocBuffer->belongs(prevPtr)){
            return prevPtr;
        }

        vmGcInfo->freeAllocationBuffers.push_back(currentAllocBuffer);
        currentAllocBuffer = currentAllocBuffer->next;

        if(currentAllocBuffer != nullptr){
            prevPtr = reinterpret_cast<Types::Object *>(currentAllocBuffer->buffer);
        }
    }

    *currentAllocBufferIndirect = currentAllocBuffer;

    return prevPtr;
}

std::size_t Memory::MarkCompact::Compacter::sizeFree(Types::Object * freePtr, AllocationBuffer * currentAllocBuffer) {
    std::size_t size = 0;

    while(currentAllocBuffer->belongs(freePtr) && !currentAllocBuffer->isMarked(freePtr)){
        size += sizeof(Types::Object);
        freePtr++;
    }

    return size;
}

void Memory::MarkCompact::Compacter::moveObjects(
        Types::Object * dst,
        Types::Object * src,
        AllocationBuffer * currentAllocBufferFree,
        AllocationBuffer * currentAllocBufferScan
) {
    void * gcDst = dst->gc;

    Types::copy(dst, src);
    dst->gc = gcDst;
    src->gc = dst;

    currentAllocBufferFree->mark(dst);
    currentAllocBufferScan->unmark(src);
}