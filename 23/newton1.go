package main

import (
	"fmt"
	"math"
)

func Sqrt(x float64) float64 {
	z := 1.0
	for i := 0; i < 10; i++ {
		z = (z * z + x) / ( 2 * z) 
	}
	return z
}

func out(n float64) (x, y, d float64) {
	x = Sqrt(n)
	y = math.Sqrt(n)
	d = x - y
	return
}

func main() {
	fmt.Println(out(1))
	fmt.Println(out(2))
	fmt.Println(out(3))
	fmt.Println(out(5))
	fmt.Println(out(7))
	fmt.Println(out(11))
}