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
	slog.Info("ConditionalMiddleware outer created")
	return func(next http.Handler) http.Handler {
		slog.Info("ConditionalMiddleware inner created")
		validatedHandler := validatorMiddleware(next) // premade middleware.
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, path := range skipPaths { // Check if the request path is in the list of paths to skip
				if r.URL.Path == path {
					slog.Info("ConditionalMiddleware skipping path", "path", path)
					next.ServeHTTP(w, r) // If it matches, bypass the validation middleware and serve the request directly
					return
				}
			}
			slog.Info("ConditionalMiddleware validating path", "path", r.URL.Path)
			validatedHandler.ServeHTTP(w, r)
		})
	}
}

func makeHandler(greeting string) http.HandlerFunc {
	slog.Info("makeHandler", "greeting", greeting)
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("handler invoked", "greeting", greeting, "path", r.URL.Path)
		_, _ = fmt.Fprintf(w, "%s", greeting)
	}
}

func authMiddleware(next http.Handler) http.Handler {
	slog.Info("authMiddleware created")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if v := r.Header.Get("Allowed"); strings.EqualFold(v, "true") {
			slog.Info("authMiddleware invoked", "verdict", "allowed")
			next.ServeHTTP(w, r)
			return
		}
		slog.Info("authMiddleware invoked", "verdict", "denied")
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("unauthorized\n"))
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	slog.Info("loggingMiddleware created")
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
