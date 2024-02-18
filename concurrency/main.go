package main

import (
	"fmt"
	"sync"
	"time"
)

// processItem represents the logic to process each item.
func processItem(item int) func() {
	return func() {
		// Simulate processing.
		time.Sleep(time.Second)
		// 20% of the time, return an error.
		if item%5 == 0 {
			fmt.Printf("Error processing item: %d\n", item)
			return
		}
		fmt.Printf("Processed item: %d\n", item)
		return
	}
}

func main() {
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var wg sync.WaitGroup
	for _, item := range items {
		wg.Add(1)
		go func(item int) {
			defer wg.Done()
			process := processItem(item)
			process()
		}(item)
	}
	wg.Wait() // Wait for all goroutines to complete.
}
