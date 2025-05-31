package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// how to inject a dependency into a function
func makeHandler(greeting string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "%s", greeting)
	}
}

func timingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info("Request served", "duration", time.Since(start), "path", r.URL.Path)
	})
}

func verbotenMiddleware(verboten string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == verboten {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
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
	mux.HandleFunc("/goodbye", makeHandler("Goodbye, world!"))
	// apply the timingMiddleware to all requests
	err := http.ListenAndServe(":8080", loggingMiddleware(timingMiddleware(verbotenMiddleware("/", mux))))
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http.ListenAndServe")
	}
	return nil
}

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}
