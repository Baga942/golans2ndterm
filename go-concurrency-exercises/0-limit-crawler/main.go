package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	Fetch(url string) (body string, urls []string, err error)
}
type SafeCounter struct {
	v   map[string]bool
	mux sync.Mutex
	wg  sync.WaitGroup
}

var c SafeCounter = SafeCounter{v: make(map[string]bool)}

func (s SafeCounter) checkvisited(url string) bool {
	s.mux.Lock()
	defer s.mux.Unlock()
	_, ok := s.v[url]
	if ok == false {
		s.v[url] = true
		return false
	}
	return true

}

func Crawl(url string, depth int, fetcher Fetcher) {
	defer c.wg.Done()
	if depth <= 0 {
		return
	}
	if c.checkvisited(url) {
		return
	}
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		c.wg.Add(1)
		go Crawl(u, depth-1, fetcher)
	}
	return
}

func main() {
	c.wg.Add(1)
	Crawl("http://golang.org/", 4, fetcher)
	c.wg.Wait()
}