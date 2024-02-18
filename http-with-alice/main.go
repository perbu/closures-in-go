package main

import (
	"fmt"
	"github.com/justinas/alice"
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

// makeVerbotenMiddleware creates a verboten Middleware. This returns a 403 Forbidden for a specific path.
func makeVerbotenMiddleware(verboten string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == verboten {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", makeHandler("Hello, world!"))
	mux.HandleFunc("/goodbye", makeHandler("Goodbye, world!"))
	// apply the timingMiddleware to all requests
	// note how we create the verboten middleware.
	chain := alice.New(timingMiddleware, makeVerbotenMiddleware("/")).Then(mux)
	err := http.ListenAndServe(":8080", chain)
	if err != nil && err != http.ErrServerClosed {
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
