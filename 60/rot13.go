package main

import (
	"io"
	"os"
	"strings"
)

type rot13Reader struct {
	r io.Reader
}

func (w rot13Reader) Read(p []byte) (n int, err error) {
	n, err = w.r.Read(p)
	if err != nil {
		return
	}
	a_up := byte('A')
	a_low := byte('a')

	for i, v := range p {
		k := -1
		if up := int(v - a_up); up >= 0 && up < 26 {
			k = up
		} else if low := int(v - a_low); low >= 0 && low < 26 {
			k = low
		} else {
			continue
		}

		if k < 13 {
			p[i] = v + 13
		} else {
			p[i] = v - 13
		}
	}
	return
}

func main() {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}
