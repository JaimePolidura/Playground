package main

import (
	"fmt"
	"math"
	"strings"
)

type GeometricalShape interface {
	volume() float64
}

type Sphere struct {
	x float64
	y float64
	z float64
	r float64
}

func (sphere Sphere) volume() float64 {
	return 4/3 * math.Pi * math.Pow(sphere.r, 2)
}

func main()  {
	sphere := Sphere{x: 0, y: 0, z: 0, r: 1}

	fmt.Println("The volue is =", calculateVolume(sphere))
	fmt.Println(strings.Join([]string{"exec", "-it", "sh"}, ", "))
}

func calculateVolume(shape GeometricalShape) float64 {
	return shape.volume()
}
