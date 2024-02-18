package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
)

type tx struct {
}

func (t tx) work() error {
	log.Println("Doing some work")
	// 50% of the time, return an error.
	if rand.Intn(2) == 0 {
		return errors.New("error")
	}
	return nil
}

// txHandler returns a function that commits or rolls back the transaction based on the error.
func txHandler(err *error) (tx, func()) {

	return tx{}, func() {
		if *err != nil {
			log.Println("TX Rolled back")
		} else {
			log.Println("TX committed")
		}
	}
}

func run() error {
	var err error
	tx, closer := txHandler(&err)
	defer closer()
	// Do some work.
	err = tx.work()
	if err != nil {
		return fmt.Errorf("error doing work: %w", err)
	}
	return nil
}

func main() {
	err := run()
	if err != nil {
		log.Println("Error:", err)
		os.Exit(1)
	}
}
