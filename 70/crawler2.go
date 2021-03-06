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
	sync.Mutex
	urls map[string]bool
}

func newTrace() *trace {
	return &trace{urls: make(map[string]bool)}
}

func (t *trace) has(url string) bool {
	t.Lock()
	defer t.Unlock()
	_, ok := t.urls[url]
	return ok
}

func (t *trace) add(url string) {
	t.Lock()
	defer t.Unlock()
	t.urls[url] = true
}

type data struct {
	url   string
	depth int
}

func work(i int, t *trace, dataWg *sync.WaitGroup, dataCh chan data, quitCh chan bool) {
	var v data
	for {
		select {
		case v = <-dataCh:
			fmt.Println(i, "Work:", v.url)
		case <-quitCh:
			fmt.Println(i, "Quitting.")
			return
		}

		if t.has(v.url) || v.depth <= 0 {
			fmt.Println(i, "Ignored:", v.url)
			dataWg.Done()
			continue
		}

		t.add(v.url)

		body, urls, err := fetcher.Fetch(v.url)
		if err != nil {
			fmt.Println(i, "Error:", err)
			dataWg.Done()
			continue
		}

		fmt.Printf("%v Found: %s %q\n", i, v.url, body)

		dataWg.Add(len(urls))

		for _, u := range urls {
			dataCh <- data{u, v.depth - 1}
		}
		dataWg.Done()
	}
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	// Fetch URLs in parallel.
	// Don't fetch the same URL twice.

	n := 3
	t := newTrace()
	dataWg := &sync.WaitGroup{}
	dataCh := make(chan data, 10)
	quitCh := make([]chan bool, n)

	dataWg.Add(1)
	dataCh <- data{url, depth}

	for i := 0; i < n; i++ {
		quit := make(chan bool, 1)
		quitCh[i] = quit
		go work(i+1, t, dataWg, dataCh, quit)
	}

	dataWg.Wait()

	for _, quit := range quitCh {
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
