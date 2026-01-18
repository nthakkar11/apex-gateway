package main

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestIdempotency(t *testing.T) {
	// Add a timestamp to the key so every test run starts "fresh"
	key := fmt.Sprintf("test-key-%d", time.Now().UnixNano())
	url := "http://localhost:8080/v1/transaction?user_id=tester"

	// First Call: Should be new (201)
	resp1, _ := sendReq(key, url)
	if resp1 != 201 {
		t.Errorf("Expected 201, got %d. (Tip: If you get 200, the key already exists in Redis)", resp1)
	}

	// Second Call: Should be cached (200)
	resp2, _ := sendReq(key, url)
	if resp2 != 200 {
		t.Errorf("Idempotency failed: Expected 200, got %d", resp2)
	}
}

// TestRateLimitConcurrent simulates 110 people clicking 'Pay' at the exact same time
func TestRateLimitConcurrent(t *testing.T) {
	url := "http://localhost:8080/v1/transaction?user_id=tester_concurrent"
	var wg sync.WaitGroup

	// Thread-safe counters for the test
	var successCount int
	var blockedCount int
	var mu sync.Mutex // Prevents race conditions on our counters during the test

	totalRequests := 110
	wg.Add(totalRequests)

	fmt.Printf("🚀 Launching %d concurrent requests...\n", totalRequests)

	for i := 0; i < totalRequests; i++ {
		// Add a tiny delay to prevent local socket exhaustion
		time.Sleep(10 * time.Millisecond)

		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("concurrent-key-%d-%d", time.Now().UnixNano(), id) // Unique key per run
			status, err := sendReq(key, url)

			if err != nil {
				return // Skip if network failed entirely
			}

			mu.Lock()
			if status == 201 {
				successCount++
			} else if status == 429 {
				blockedCount++
			}
			mu.Unlock()
		}(i)
	}
	wg.Wait() // Wait for all 110 Goroutines to finish

	fmt.Printf("🏁 Results -> Allowed: %d, Rate-Limited: %d\n", successCount, blockedCount)

	if successCount != 100 {
		t.Errorf("Rate limit failed logic: expected 100 allowed, got %d", successCount)
	}
}

func sendReq(key, url string) (int, error) {
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("X-Idempotency-Key", key)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}
