package main

import (
	"fmt"
	"strings"
)

func main() {
	animals := "dog fish cat"
	for i, animal := range superSplit(animals, " ") {
		fmt.Println(i, animal)
	}
}

func superSplit(s, sep string) func(func(int, string) bool) {
	return func(yield func(int, string) bool) {
		for i := 0; len(s) > 0; i++ {
			j := strings.Index(s, sep)
			if j == -1 {
				yield(i, s)
				return
			}
			if !yield(i, s[:j]) {
				return
			}
			s = s[j+len(sep):]
		}
	}
}
