package main

import "code.google.com/p/go-tour/pic"

func Pow(b, e int) int {
	v := 1
	for i := 0; i < e; i++ {
		v *= b
	}
	return v
}

func Pic(dx, dy int) [][]uint8 {
	m := make([][]uint8, dy)
	for y := 0; y < dy; y++ {
		m[y] = make([]uint8, dx)
		for x := 0; x < dx; x++ {
			// m[y][x] = uint8((x + y) / 2)
			// m[y][x] = uint8(x * y)
			m[y][x] = uint8(Pow(x, y))
		}
	}
	return m
}

func main() {
	pic.Show(Pic)
}
