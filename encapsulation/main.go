package main

import (
	"fmt"
)

// NewCounter returns a closure function that increments and returns its internal state (count) each time it is called.
func NewCounter(initial int) func() int {
	// 'count' is a variable encapsulated by the closure returned below.
	count := initial
	// The closure that is being returned.
	return func() int {
		count++      // Increment the count.
		return count // Return the updated count.
	}
}

func main() {
	counterA := NewCounter(0)  // Create a new counter starting at 0.
	counterB := NewCounter(10) // Create another counter starting at 10.

	// Demonstrate that each counter maintains its own state.
	fmt.Println(counterA()) // Outputs: 1
	fmt.Println(counterA()) // Outputs: 2
	fmt.Println(counterB()) // Outputs: 11
	fmt.Println(counterB()) // Outputs: 12
	fmt.Println(counterA()) // Outputs: 3
}
