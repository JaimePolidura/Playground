package main

import (
	"fmt"
)

func main() {
	x := [5]float64{98, 93, 77, 72, 83}
	var total float64 = 0
	for _, value := range x {
		total += value
	}
	fmt.Println(total / float64(len(x)))

	arr := [5]uint8{1, 2, 3, 4, 5}
	subArray := arr[0:2] //Devuelve [1, 2]
	fmt.Println(subArray)

	slice1 := []int{1,2,3}
	slice2 := append(slice1, 4, 5)
	fmt.Println(slice2)

	maps()
}

func maps()  {
	myMap := make(map[string]int)
	myMap["jaime"] = 1;
	myMap["molon"] = 2;

	delete(myMap, "jaime")

	fmt.Println(myMap["molon"])

	result, err := myMap["molon"]
	fmt.Println(result, err)

	elements := map[string]map[string]string{
		"H": map[string]string{
			"name":  "Hydrogen",
			"state": "gas",
		},
		"He": map[string]string{
			"name":  "Helium",
			"state": "gas",
		},
	}

	fmt.Println(elements)
}

func copyArrays()  {
	slice1 := []int{1,2,3}
	slice2 := make([]int, 2)
	copy(slice2, slice1)
	fmt.Println(slice1, slice2)
}