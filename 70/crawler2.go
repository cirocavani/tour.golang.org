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

func work(i int, t *trace, group *sync.WaitGroup, data_ch chan data, quit_ch chan bool) {
	var v data
	for {
		select {
		case v = <-data_ch:
			fmt.Println(i, "Work:", v.url)
		case <-quit_ch:
			fmt.Println(i, "Quitting.")
			return
		}

		if t.has(v.url) || v.depth <= 0 {
			fmt.Println(i, "Ignored:", v.url)
			group.Done()
			continue
		}

		t.add(v.url)

		body, urls, err := fetcher.Fetch(v.url)
		if err != nil {
			fmt.Println(i, "Error:", err)
			group.Done()
			continue
		}

		fmt.Printf("%v Found: %s %q\n", i, v.url, body)

		group.Add(len(urls))

		for _, u := range urls {
			data_ch <- data{u, v.depth - 1}
		}
		group.Done()
	}
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	// Fetch URLs in parallel.
	// Don't fetch the same URL twice.

	n := 3
	t := new_trace()
	g := &sync.WaitGroup{}
	data_ch := make(chan data, 10)
	quit_ch := make([]chan bool, n)

	g.Add(1)
	data_ch <- data{url, depth}

	for i := 0; i < n; i++ {
		quit := make(chan bool, 1)
		quit_ch[i] = quit
		go work(i+1, t, g, data_ch, quit)
	}
	g.Wait()
	for _, quit := range quit_ch {
		quit <- true
	}
	time.Sleep(10 * time.Millisecond)
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
