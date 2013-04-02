package main

import "fmt"

func pow(b complex128, e int) complex128 {
	z := b
	for i := 1; i < e; i++ {
		z *= b
	}
	return z
}
 
func Cbrt(x complex128) complex128 {
	z := complex128(1.0);
	for i := 0; i < 10; i++ {
		z = (2 * pow(z, 3) + x ) / (3 * pow(z, 2))
	}
	return z
}

func main() {
	fmt.Println(Cbrt(2))
}
