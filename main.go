package main

import (
    "fmt"
    "sync"
    "io"
    "net/http"
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
    for url := range jobs {
        resp, err := http.Get(url)
        if err != nil {
            results <- FetchResult{URL: url, Error: err}
            continue
        }
        defer resp.Body.Close()

        body, _ := io.ReadAll(resp.Body)
        results <- FetchResult{
            URL:        url,
            StatusCode: resp.StatusCode,
            Size:       len(body),
            Error:      nil,
        }
    }
}
func main() {
    urls := []string{
        "https://example.com",
        "https://golang.org",
        "https://uottawa.ca",
        "https://github.com",
        "https://httpbin.org/get",
    }

    numWorkers := 1

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

