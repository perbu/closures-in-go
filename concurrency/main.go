package main

import (
	"fmt"
	"sync"
	"time"
)

// processItem represents the logic to process each item.
func processItem(item int) func() error {
	return func() error {
		// Simulate processing.
		time.Sleep(time.Second * 2)
		// 20% of the time, return an error.
		if item%5 == 0 {
			fmt.Printf("Error processing item: %d\n", item)
			return fmt.Errorf("error processing item: %d", item)
		}
		fmt.Printf("Processed item: %d\n", item)
		return nil
	}
}

func main() {
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var wg sync.WaitGroup
	var errCh = make(chan error, len(items))
	concurrencyGuardChannel := make(chan bool, 3)
	for _, item := range items {
		wg.Add(1)
		go func(item int) {
			defer wg.Done()
			concurrencyGuardChannel <- true
			defer func() {
				_ = <-concurrencyGuardChannel
			}()
			process := processItem(item)
			err := process()
			if err != nil {
				errCh <- err
				return
			}
		}(item)
	}
	wg.Wait() // Wait for all goroutines to complete.
	close(errCh)
	for err := range errCh {
		fmt.Println(err)
	}
}
