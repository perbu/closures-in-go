package main

import (
	"fmt"
)

// NewCounter returns a closure function that increments and returns its internal state (count) each time it is called.
func NewCounter(initial int) (func(int) int, func()) {
	// 'count' is a variable encapsulated by the closure returned below.
	count := initial
	// The closure that is being returned.
	return func(inc int) int {
			count = count + inc // Increment the count.
			return count        // Return the updated count.
		}, func() {
			count = initial
		}
}

func main() {
	counterA, resetA := NewCounter(0)  // Create a new counter starting at 0.
	counterB, resetB := NewCounter(10) // Create another counter starting at 10.

	// Demonstrate that each counter maintains its own state.
	fmt.Println(counterA(2)) // Outputs: 1
	fmt.Println(counterA(2)) // Outputs: 2
	fmt.Println(counterB(2)) // Outputs: 11
	fmt.Println(counterB(2)) // Outputs: 12
	fmt.Println(counterA(2)) // Outputs: 3
	fmt.Println("=====")
	resetA()
	fmt.Println(counterA(2)) // Outputs: 1
	resetB()
	fmt.Println(counterB(2)) // Outputs: 11
}
