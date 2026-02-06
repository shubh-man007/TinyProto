package main

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	url := "http://localhost:8080"
	totalRequests := 10000
	concurrency := 100

	var successful, failed atomic.Int64
	var totalLatency atomic.Int64

	start := time.Now()
	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrency)

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		sem <- struct{}{}

		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			reqStart := time.Now()
			resp, err := http.Get(url)
			latency := time.Since(reqStart)

			if err != nil || resp.StatusCode != http.StatusOK {
				failed.Add(1)
				if resp != nil && resp.Body != nil {
					resp.Body.Close()
				}
				return
			}

			successful.Add(1)
			totalLatency.Add(int64(latency))
			resp.Body.Close()
		}()
	}

	wg.Wait()
	duration := time.Since(start)

	successCount := successful.Load()
	failCount := failed.Load()

	var avgLatency time.Duration
	if successCount > 0 {
		avgLatency = time.Duration(totalLatency.Load() / successCount)
	}

	fmt.Println("=== Load Test Results ===")
	fmt.Printf("Total Requests: %d\n", totalRequests)
	fmt.Printf("Successful: %d\n", successCount)
	fmt.Printf("Failed: %d\n", failCount)
	fmt.Printf("Duration: %v\n", duration)
	fmt.Printf("RPS: %.2f\n", float64(totalRequests)/duration.Seconds())
	fmt.Printf("Average Latency: %v\n", avgLatency)
	fmt.Printf("Concurrency: %d\n", concurrency)
}

