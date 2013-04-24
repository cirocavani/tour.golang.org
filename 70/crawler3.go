package main

import (
	"fmt"
	"time"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type data struct {
	url    string
	depth  int
	result []string
}

func broker(workToDoCh, workDoneCh chan data, doneCh chan bool) {
	workDoing := make(map[string]bool)
	workDone := make(map[string]bool)

	doing := func(url string) bool {
		_, ok := workDoing[url]
		return ok
	}

	done := func(url string) bool {
		_, ok := workDone[url]
		return ok
	}

	old := func(url string) bool {
		return done(url) || doing(url)
	}

	for {
		v := <-workDoneCh
		workDone[v.url] = true
		delete(workDoing, v.url)
		if len(v.result) > 0 && v.depth > 0 {
			for _, url := range v.result {
				if old(url) {
					fmt.Println("Ignored:", url)
					continue
				}
				workDoing[url] = true
				workToDoCh <- data{url: url, depth: v.depth - 1}
			}
		}
		if len(workDoing) == 0 {
			break
		}
	}

	doneCh <- true
}

func worker(i int, workToDoCh, workDoneCh chan data, quitCh chan bool) {
	var v data
	for {
		select {
		case v = <-workToDoCh:
			fmt.Println(i, "Work:", v.url)
		case <-quitCh:
			fmt.Println(i, "Quitting.")
			return
		}

		body, urls, err := fetcher.Fetch(v.url)
		if err != nil {
			fmt.Println(i, "Error:", err)
			workDoneCh <- v
			continue
		}

		fmt.Printf("%v Found: %s %q\n", i, v.url, body)

		workDoneCh <- data{v.url, v.depth, urls}
	}
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	// Fetch URLs in parallel.
	// Don't fetch the same URL twice.

	n := 3
	workToDoCh := make(chan data, 10)
	workDoneCh := make(chan data, 10)
	doneCh := make(chan bool, 1)
	quitCh := make(chan bool, 1)

	go broker(workToDoCh, workDoneCh, doneCh)

	for i := 0; i < n; i++ {
		go worker(i+1, workToDoCh, workDoneCh, quitCh)
	}

	workToDoCh <- data{url: url, depth: depth}

	<-doneCh
	for i := 0; i < n; i++ {
		quitCh <- true
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
