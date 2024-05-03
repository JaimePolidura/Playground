.global _start
.intel_syntax noprefix

.section .text

branch:
    shl eax, 2 # eax << 2, may set the carry flag
    shr eax, 2 # eax >> 2
    jmp finish

_start:
    mov eax, 10
    mov ebx, 10
    cmp eax, ebx
    
    je branch # jump if equal

finish:
    # sys_exit
    mov rax, 60
    mov rdi, 1
    syscall

.section .data