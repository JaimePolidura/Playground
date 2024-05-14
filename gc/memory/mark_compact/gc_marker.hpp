#pragma once

#include "shared.hpp"
#include "vm/vm.hpp"
#include "types/array.hpp"
#include "types/struct.hpp"
#include "types/string.hpp"
#include "types/utils.hpp"
#include "memory/mark_compact/gc_thread_info.hpp"
#include "memory/mark_compact/gc_vm_info.h"

namespace Memory::MarkCompact {
    class Marker {
    public:
        void mark();

    private:
        void markPackage(std::shared_ptr<VM::Package> package);
        void markStack(VM::Thread& thread);
        void markThreadsStack();
        void markPackages();
        void traverseObject(Types::Object * rootObject);
        void markObject(Types::Object * object);
        bool isMarked(Types::Object * object);
    };
}
