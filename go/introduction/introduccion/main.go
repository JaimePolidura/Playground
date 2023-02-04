package main

import "fmt"

const e float64 = 2.71

func main() {
	fmt.Println("La suma de 1 + 2 =", 1+2)
	fmt.Println(len("tres"))
	fmt.Println("ABC"[1])

	var nombre string = "jaime"
	fmt.Println("Mi nombre es " + nombre)

	number := 112
	fmt.Println(number)

	printNumberE()

	var (
		a = 1
		b = 2
	)

	fmt.Println(a + b)
}

func printNumberE() {
	fmt.Println(e)
}
