package main

import "fmt"

// Filter remains as previously described.
func Filter(slice []int, predicate func(int) bool) []int {
	var result []int
	for _, value := range slice {
		if predicate(value) {
			result = append(result, value)
		}
	}
	return result
}

func main() {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	limit := 5 // This variable will be captured by the closure passed to Filter.

	// Use Filter with an anonymous function (closure) that captures 'limit'.
	belowLimitNumbers := Filter(numbers, func(n int) bool {
		return n < limit // 'limit' is captured from the surrounding scope.
	})

	fmt.Println("Numbers below limit:", belowLimitNumbers)
}
