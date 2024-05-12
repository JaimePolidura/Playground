cmake_minimum_required(VERSION 3.28)
project(untitled7)

set(CMAKE_CXX_STANDARD 26)

include_directories(${CMAKE_SOURCE_DIR})

add_executable(untitled7 main.cpp
        types/array.hpp
        types/types.hpp
        shared.hpp
        types/string.hpp
        types/struct.hpp
        memory/allocator.hpp
        vm/thread.hpp
        vm/vm.hpp
        vm/package.hpp
        vm/vm.cpp
        vm/thread.cpp
        memory/mark_compact/allocation_buffer.hpp
        memory/mark_compact/allocation_buffer.cpp
        memory/mark_compact/allocator.hpp
        memory/mark_compact/allocator.cpp
        memory/mark_compact/gc_thread_info.hpp
        memory/mark_compact/gc_marker.hpp
        memory/mark_compact/gc_marker.cpp
        memory/mark_compact/gc_compacter.hpp
        memory/mark_compact/gc_compacter.cpp
        types/utils.hpp
        types/utils.cpp
)
