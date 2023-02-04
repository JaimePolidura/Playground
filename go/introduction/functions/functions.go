package main

import (
	"fmt"
)

func main()  {
	defer funcionInutil()

	min, max := minMaxOf(4, 6, 1, 3, 7, 1, 1)

	fmt.Println("The min value =", min)
	fmt.Println("The max value =", max)

	lambda := func(a int) int{ return a + 1}
	fmt.Println(lambda(1))

	ptr := heapAllocation()
	fmt.Println("Heap allocation =", * ptr)
}

func heapAllocation() * int  {
	ptr := new(int)
	* ptr = 12

	return ptr
}

func funcionInutil() {
	fmt.Println("Un inutil ha llamado a la funcion inutil")
}

func minMaxOf(args ...int) (int, int)  {
	array := make([]int, len(args))
	for i, it := range args {
		array[i] = it;
	}

	return minMax(array)
}

func minMax(arr []int) (int, int) {
	actualMax := arr[0]
	actualMin := arr[0]

	for i := 0; i < len(arr); i++ {
		if arr[i] > actualMax{
			actualMax = arr[i]
		}
		if arr[i] < actualMin {
			 actualMin = arr[i]
		}
	}

	return actualMin, actualMax
}
