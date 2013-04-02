package main

import "fmt"

func fibonacci() func() int {
	i := 0
	last1 := 0
	last2 := 1
	return func() int {
		if i == 0 {
			i++
			return last1
		}
		if i == 1 {
			i++
			return last2
		}
		v := last1 + last2
		last1 = last2
		last2 = v
		return v
	}
}

func main() {
    f := fibonacci()
    for i := 0; i < 10; i++ {
        fmt.Println(f())
    }
}
