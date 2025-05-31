package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

// ConditionalMiddleware is a function that returns a middleware function
// that applies to everything but the paths in the skipPaths list
func ConditionalMiddleware(skipPaths []string, validatorMiddleware func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, path := range skipPaths { // Check if the request path is in the list of paths to skip
				if r.URL.Path == path {
					next.ServeHTTP(w, r) // If it matches, bypass the validation middleware and serve the request directly
					return
				}
			}
			validatedHandler := validatorMiddleware(next) // If the path is not in the skip list, apply the validation middleware
			validatedHandler.ServeHTTP(w, r)
		})
	}
}

// how to inject a dependency into a function
func makeHandler(greeting string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "%s", greeting)
	}
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if v := r.Header.Get("Allowed"); strings.EqualFold(v, "true") {
			next.ServeHTTP(w, r)
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("unauthorized\n"))
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t0 := time.Now()
		next.ServeHTTP(w, r)
		slog.Info("Request served", "duration", time.Since(t0))
	})
}

func run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", makeHandler("Hello, world!"))
	mux.HandleFunc("/skip", makeHandler("This is skipped!"))

	wrapped := ConditionalMiddleware(
		[]string{"/skip"}, // paths to skip
		authMiddleware,    // the validator
	)

	log.Fatal(http.ListenAndServe(":8080", wrapped(loggingMiddleware(mux))))
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}
