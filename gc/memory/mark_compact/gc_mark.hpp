#pragma once

#include "shared.hpp"
#include "vm/vm.hpp"
#include "types/array.hpp"
#include "types/struct.hpp"
#include "types/string.hpp"
#include "memory/mark_compact/gc_thread_info.hpp"

namespace Memory::MarkCompact {
    class Marker {
    public:
        void mark();

    private:
        std::map<absoluteAllocBufAddress_t , AllocationBuffer *> allocationsBufferByAddress{};

        void markPackage(std::shared_ptr<VM::Package> package);
        void markStack(VM::Thread& thread);
        void markThreadsStack();
        void markPackages();
        void traverseObject(Types::Object * object);
        void traverseStruct(Types::StructObject * structObject, std::queue<Types::Object *>& pending);
        void traverseArray(Types::ArrayObject * arrayObject, std::queue<Types::Object *>& pending);
        void markObject(Types::Object * object);
        bool isMarked(Types::Object * object);

        AllocationBuffer * getAllocationBuffer(absoluteAllocBufAddress_t addressToLookup);
    };
}
