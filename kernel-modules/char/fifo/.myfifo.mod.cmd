cmd_/home/jaime/OtherLanguagesPlayground/kernel-modules/char/fifo/myfifo.mod := printf '%s\n'   myfifo.o | awk '!x[$$0]++ { print("/home/jaime/OtherLanguagesPlayground/kernel-modules/char/fifo/"$$0) }' > /home/jaime/OtherLanguagesPlayground/kernel-modules/char/fifo/myfifo.mod
