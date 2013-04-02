package main

import (
	"fmt"
	"math"
)

func Sqrt(x, d float64) (z, i float64) {
	z = 1.0
	for i = 1;; i++ {
		t := (z * z + x) / ( 2 * z)
		if math.Abs(t - z) < d {
			break
		}
		z = t
	}
	return
}

func out(n, d float64) (x, y, i float64) {
	x, i = Sqrt(n, d)
	y = math.Sqrt(n)
	return
}

func main() {
	const d = 1e-10
	fmt.Println(out(1, d))
	fmt.Println(out(2, d))
	fmt.Println(out(3, d))
	fmt.Println(out(5, d))
	fmt.Println(out(7, d))
	fmt.Println(out(11, d))
}