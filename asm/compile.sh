as $1 -o output.o || exit
gcc -o output output.o -nostdlib -static || exit
./output

rm output
rm output.o
