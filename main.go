package main

import (
    "fmt"
    "sync"
)

// WaitGroup is used to wait for all workers to finish.
var wg sync.WaitGroup

// Structure to store results
type FetchResult struct {
    URL        string
    StatusCode int
    Size       int
    Error      error
}

// Worker function
func worker(id int, jobs <-chan string, results chan<- FetchResult) {
    defer wg.Done()
    // TODO: fetch the URL
    // TODO: send Result struct to results channel
    // hint: use resp, err := http.Get(url)
}

func main() {
    urls := []string{
        "https://example.com",
        "https://golang.org",
        "https://uottawa.ca",
        "https://github.com",
        "https://httpbin.org/get",
    }

    numWorkers := 3

    jobs := make(chan string, len(urls))
    results := make(chan FetchResult, len(urls))

    // Start workers
    wg.Add(numWorkers)
    for id := 1; id <= numWorkers; id++ {
        go worker(id, jobs, results)
    }

    // Send jobs
    for _, url := range urls {
        jobs <- url
    }
    close(jobs)

    // Collect results
    for i := 0; i < len(urls); i++ {
        res := <-results
        if res.Error != nil {
            fmt.Printf("URL: %s | error: %v\n", res.URL, res.Error)
        } else {
            fmt.Printf("URL: %s | status: %d | size: %d bytes\n",
                res.URL, res.StatusCode, res.Size)
        }
    }

    wg.Wait()
    fmt.Println("\nScraping complete!")
}

