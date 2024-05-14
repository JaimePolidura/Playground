#include "gc_marker.hpp"

extern struct VM::VM current_vm;

void Memory::MarkCompact::Marker::mark() {
    current_vm.stopThreadsGC();
    markThreadsStack();
    markPackages();
    GET_VM_GC_INFO(current_vm)->resetAllMarkBits();
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
        traverseObject(global.value);
    }
}

void Memory::MarkCompact::Marker::traverseObject(Types::Object * rootObject) {
    Types::traverseObjectDeep(rootObject, [this](Types::Object * currentObject) -> bool {
        if(currentObject != nullptr && !this->isMarked(currentObject)){
            this->markObject(currentObject);
            return true; ////Keep going
        } else {
            return false; //Stop
        }
    });
}

void Memory::MarkCompact::Marker::markObject(Types::Object * object) {
    auto allocationBuffer = GET_ALLOCATION_BUFFER(current_vm, object);
    allocationBuffer->mark(object);
}

bool Memory::MarkCompact::Marker::isMarked(Types::Object * object) {
    auto allocationBuffer = GET_ALLOCATION_BUFFER(current_vm, object);
    return allocationBuffer->isMarked(object);
}