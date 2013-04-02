package main

import (
	"code.google.com/p/go-tour/wc"
	"strings"
)

func WordCount(s string) map[string]int {
	m := make(map[string]int)
	for _, v := range strings.Fields(s) {
		n := m[v] + 1
		m[v] = n
	}
	return m
}

func main() {
	wc.Test(WordCount)
}
