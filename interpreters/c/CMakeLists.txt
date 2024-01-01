cmake_minimum_required(VERSION 3.27)
project(c C)

set(CMAKE_C_STANDARD 23)

add_executable(c main.c
        shared.h
        bytecode.h
        chunk.h
        chunk.c
        memory.h
        memory.c)
