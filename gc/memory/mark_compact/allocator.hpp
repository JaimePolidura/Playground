#pragma once

#include "shared.hpp"

#include "memory/allocator.hpp"
#include "vm/thread.hpp"
#include "memory/mark_compact/allocation_buffer.hpp"
#include "memory/mark_compact/gc_thread_info.hpp"
#include "memory/mark_compact/gc_compact.hpp"
#include "memory/mark_compact/gc_mark.hpp"

namespace Memory::MarkCompact {
    class MarkCompactAllocator : public Memory::Allocator{
    public:
        Types::StringObject * allocString(char * data) override;

        Types::StructObject * allocStruct(int nFields, Types::Object * fields) override;

        Types::ArrayObject * allocArray(Types::ObjectType contentType, int nElements) override;

    private:
        void * allocateSize(size_t size);

        void startGC();
    };
}