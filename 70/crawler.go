package main

import (
	"fmt"
	"sync"
	"time"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type pool struct {
	group *sync.WaitGroup
	mutex sync.Mutex
	n     int
}

func (p *pool) lock() {
	p.mutex.Lock()
}

func (p *pool) unlock() {
	p.mutex.Unlock()
}

func (p *pool) quit() {
	p.group.Done()
}

func (p *pool) work() {
	p.lock()
	p.n++
	p.unlock()
}
func (p *pool) done() {
	p.lock()
	p.n--
	p.unlock()
}
func (p *pool) empty() bool {
	p.lock()
	empty := p.n == 0
	p.unlock()
	return empty
}

type trace struct {
	mutex sync.Mutex
	urls  map[string]bool
}

func new_trace() *trace {
	return &trace{urls: make(map[string]bool)}
}

func (t *trace) lock() {
	t.mutex.Lock()
}

func (t *trace) unlock() {
	t.mutex.Unlock()
}

func (t *trace) has(url string) bool {
	t.lock()
	_, ok := t.urls[url]
	t.unlock()
	return ok
}

func (t *trace) add(url string) {
	t.lock()
	t.urls[url] = true
	t.unlock()
}

type data struct {
	url   string
	depth int
}

func work(i int, p *pool, t *trace, ch chan data) {
	var v data
	for {
		select {
		case v = <-ch:
			p.work()
			fmt.Println(i, "Work:", v.url)
		default:
			if p.empty() {
				fmt.Println(i, "Quitting.")
				p.quit()
				return
			} else {
				time.Sleep(10 * time.Millisecond)
				continue
			}
		}

		if t.has(v.url) || v.depth <= 0 {
			fmt.Println(i, "Ignored:", v.url)
			p.done()
			continue
		}

		t.add(v.url)

		body, urls, err := fetcher.Fetch(v.url)
		if err != nil {
			fmt.Println(i, "Error:", err)
			p.done()
			continue
		}

		fmt.Printf("%v Found: %s %q\n", i, v.url, body)
		for _, u := range urls {
			ch <- data{u, v.depth - 1}
		}
		p.done()
	}
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	// Fetch URLs in parallel.
	// Don't fetch the same URL twice.

	n := 3
	g := &sync.WaitGroup{}
	p := &pool{group: g}
	t := new_trace()
	ch := make(chan data, 10)

	ch <- data{url, depth}

	g.Add(n)
	for i := 0; i < n; i++ {
		go work(i+1, p, t, ch)
	}
	g.Wait()
}

func main() {
	Crawl("http://golang.org/", 4, fetcher)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f *fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := (*f)[url]; ok {
		time.Sleep(10 * time.Millisecond)
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = &fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}
