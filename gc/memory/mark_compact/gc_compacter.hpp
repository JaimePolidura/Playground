#pragma once

#include "shared.hpp"
#include "vm/vm.hpp"
#include "types/array.hpp"
#include "types/struct.hpp"
#include "types/string.hpp"
#include "memory/mark_compact/gc_thread_info.hpp"
#include "types/utils.hpp"

namespace Memory::MarkCompact {
    class Compacter {
    public:
        void compact();
    private:
        void compactThreads();
        void compactThread(const VM::Thread& thread);
        inline Types::Object * nextLive(Types::Object * prevPtr, AllocationBuffer * currentAllocBuffer);
        inline Types::Object * nextFree(Types::Object * prevPtr, AllocationBuffer * currentAllocBuffer);
        inline std::size_t sizeFree(Types::Object * freePtr, AllocationBuffer * currentAllocBuffer);
        void moveObjects(Types::Object * dst, Types::Object * src, AllocationBuffer * currentAllocBuffer);

        void updateReferences();
        void updateStackReferences();
        void updateGlobalReferences();
        void updateObjectReferencesDeep(Types::Object * object);
    };
}