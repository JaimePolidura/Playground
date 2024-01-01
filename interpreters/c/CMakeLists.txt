cmake_minimum_required(VERSION 3.27)
project(c C)

set(CMAKE_C_STANDARD 23)

add_executable(c main.c
        shared.h
        bytecode.h
        chunk/chunk.h
        chunk/chunk.c
        memory/memory.h
        memory/memory.c
        chunk/chunk_disassembler.h
        chunk/chunk_disassembler.c
        types/types.h
        types/types.c)
