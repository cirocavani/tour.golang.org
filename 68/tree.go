package main

import (
	"code.google.com/p/go-tour/tree"
	"container/list"
	"fmt"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	m := list.New()
	m.PushBack(t)
	for e:= m.Front(); e != nil; e = e.Next() {
		n := e.Value.(*tree.Tree)
		ch <- n.Value
		if n.Left != nil {
			m.PushBack(n.Left)
		}
		if n.Right != nil {
			m.PushBack(n.Right)
		}
	}
	close(ch)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int, 10)
	ch2 := make(chan int, 10)

	go Walk(t1, ch1)
	go Walk(t2, ch2)

	eval := func (m, n *map[int]bool, k int, ok bool) (closed, end bool) {
		_m, _n := *m, *n
		end  = false
		if closed = !ok; closed {
			end = len(_n) != 0
		} else if _, nok := _n[k]; nok {
			delete(_n, k)
		} else {
			_m[k] = true
		}
		return
	}

	m1 := make(map[int]bool)
	m2 := make(map[int]bool)

	for closed1, closed2, end := false, false, false; !closed1 || !closed2; {
		select {
		case k, ok := <- ch1:
			if closed1, end = eval(&m1, &m2, k, ok); end {
				return false
			}
		case k, ok := <- ch2:
			if closed2, end = eval(&m2, &m1, k, ok); end {
				return false
			}
		}
	}
	return len(m1) == 0 && len(m2) == 0
}

func main() {
	ch := make(chan int)
	go Walk(tree.New(1), ch)
	for i := range ch {
		fmt.Println(i)
	}
	fmt.Println(Same(tree.New(1), tree.New(1)))
	fmt.Println(Same(tree.New(2), tree.New(1)))
}
