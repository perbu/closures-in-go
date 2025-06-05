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

type ConditionalHandler struct {
	skip      []string
	next      http.Handler
	validated http.Handler
}

func NewConditionalHandler(
	skip []string,
	validator func(http.Handler) http.Handler,
	next http.Handler,
) http.Handler {
	return &ConditionalHandler{
		skip:      skip,
		next:      next,
		validated: validator(next),
	}
}

func (h *ConditionalHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, path := range h.skip {
		if r.URL.Path == path {
			h.next.ServeHTTP(w, r)
			return
		}
	}
	h.validated.ServeHTTP(w, r)
}

// copy slice into a set for constant-time look-ups
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
	mux.HandleFunc("/skip", makeHandler("We're skipping auth!"))

	handler := NewConditionalHandler([]string{"/skip"},
		authMiddleware,
		loggingMiddleware(mux))
	log.Println("listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}
