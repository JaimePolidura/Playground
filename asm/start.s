.global _start
.intel_syntax noprefix

.section .text
my_empty_function:
    ret

on_overflow:
    ret

_start:
    mov eax, 10
    mov ebx, 12
    add eax, ebx # eax = eax + ebx
    jc on_overflow # if overflow flag is set

    # sys_write (print hello world)
    mov rax, 1 # read syscall
    mov rdi, 1 # stdout
    lea rsi, [hello_world_buffer] # hello world address
    mov rdx, hello_world_size # hello world size
    syscall
    
    call my_empty_function

    # sys_exit
    mov rax, 60
    mov rdi, 1
    syscall

.section .data
hello_world_buffer:
    .asciz "Hello, world!\n"

hello_world_size:
    .long 14
