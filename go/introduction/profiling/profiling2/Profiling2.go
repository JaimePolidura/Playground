package main

import (
	"fmt"
	"github.com/pkg/profile"
	"math"
)

func main()  {
	defer profile.Start(profile.ProfilePath("/tmp")).Stop()

	total := 0
	for i := 2; i < 200000; i++ {
		n := N1(i)
		if n {
			total++
		}
	}
	fmt.Println("Total: ", total)
	total = 0
	k := 0
	for i := 0; i < 5000; i++ {
		for j := 0; j < 400; j++ {
			k = Multiply(i, j)
			total++
		}
	}
	fmt.Println("Total: ", total, k)
}

func N1(n int) bool {
	k := math.Floor(float64(n/2 + 1))
	for i := 2; i < int(k); i++ {
		if (n % i) == 0 {
			return false
		}
	}
	return true
}

func Multiply(a, b int) int {
	if a == 1 {
		return b
	}
	if a == 0 || b == 0 {
		return 0
	}
	if a < 0 {
		return -Multiply(-a, b)
	}

	return b + Multiply(a-1, b)
}