package main

import (
	"fmt"
)

type CrawlResult struct {
	depth int
	url   string
	body  string
	urls  []string
	err   error
}

func Worker(url string, depth int, fetcher Fetcher, out chan CrawlResult) {
	body, urls, err := fetcher.Fetch(url)
	out <- CrawlResult{
		depth: depth,
		url:   url,
		body:  body,
		urls:  urls,
		err:   err,
	}
}

func Crawl1(url string, depth int, fetcher Fetcher) {
	if depth <= 0 {
		return
	}
	resultChannel := make(chan CrawlResult)
	visited := make(map[string]bool)

	visited[url] = true
	workerCount := 0
	workerCount++
	go Worker(url, depth, fetcher, resultChannel)

	for workerCount > 0 {
		res := <-resultChannel
		workerCount--

		if res.err != nil {
			fmt.Println(res.err)
			continue
		}
		fmt.Printf("found: %s %q\n", res.url, res.body)

		if res.depth > 0 {
			for _, resurl := range res.urls {
				if !visited[resurl] {
					visited[resurl] = true
					workerCount++
					go Worker(resurl, res.depth-1, fetcher, resultChannel)
				}
			}
		}
	}
}

// Strategy 2: Channel Coordinator
func test2() {
	Crawl1("https://golang.org/", 4, fetcher)
}
