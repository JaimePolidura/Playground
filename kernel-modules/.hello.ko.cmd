cmd_/home/jaime/OtherLanguagesPlayground/kernel-modules/hello.ko := ld -r -m elf_x86_64 -z noexecstack --build-id=sha1  -T scripts/module.lds -o /home/jaime/OtherLanguagesPlayground/kernel-modules/hello.ko /home/jaime/OtherLanguagesPlayground/kernel-modules/hello.o /home/jaime/OtherLanguagesPlayground/kernel-modules/hello.mod.o;  true