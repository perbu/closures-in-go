package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

// ConditionalMiddleware applies validator to the nextHandler
// for all paths except those in skipPaths.
func ConditionalMiddleware(
	nextHandler http.Handler, // The handler this middleware will wrap or bypass
	skipPaths []string, // A list of URL paths to bypass the validator for
	validator func(http.Handler) http.Handler, // The middleware to conditionally apply
) http.Handler {
	slog.Info("ConditionalMiddleware applying to handler", "next_handler_type", fmt.Sprintf("%T", nextHandler))

	// Prepare the handler that has the validator applied to nextHandler.
	// This is done once when ConditionalMiddleware is called to set up the chain.
	// Notice how it happens in the outer closure.
	validatedHandler := validator(nextHandler)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currentPath := r.URL.Path
		for _, skipPath := range skipPaths { // Iterate over the paths to skip
			if currentPath == skipPath {
				slog.Info("ConditionalMiddleware: skipping validation", "path", currentPath)
				// If the current request path is in skipPaths, serve the request
				// using the original nextHandler, bypassing the validator.
				nextHandler.ServeHTTP(w, r)
				return
			}
		}
		// If the path is not in skipPaths, validation is required.
		slog.Info("ConditionalMiddleware: applying validation", "path", currentPath)
		// Serve the request using the validatedHandler (validator applied to nextHandler).
		validatedHandler.ServeHTTP(w, r)
	})
}

func makeHandler(greeting string) http.HandlerFunc {
	slog.Info("makeHandler created", "greeting", greeting)
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("handler invoked", "greeting", greeting, "path", r.URL.Path)
		_, _ = fmt.Fprintf(w, "%s: %s\n", greeting, r.URL.Path)
	}
}

func authMiddleware(next http.Handler) http.Handler {
	slog.Info("authMiddleware created, will wrap", "next_handler_type", fmt.Sprintf("%T", next))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if v := r.Header.Get("Allowed"); strings.EqualFold(v, "true") {
			slog.Info("authMiddleware: allowed", "path", r.URL.Path)
			next.ServeHTTP(w, r) // Proceed to the next handler
			return
		}
		slog.Info("authMiddleware: denied", "path", r.URL.Path)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("unauthorized\n"))
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	slog.Info("loggingMiddleware created, will wrap", "next_handler_type", fmt.Sprintf("%T", next))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t0 := time.Now()
		slog.Info("loggingMiddleware: request received", "path", r.URL.Path, "method", r.Method)
		next.ServeHTTP(w, r) // Proceed to the next handler
		slog.Info("loggingMiddleware: request completed", "path", r.URL.Path, "duration", time.Since(t0))
	})
}

func run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", makeHandler("Hello, world!"))
	mux.HandleFunc("/skip", makeHandler("This is a public endpoint"))

	handlerWithLogging := loggingMiddleware(mux)

	finalHandler := ConditionalMiddleware(
		handlerWithLogging, // The handler to process
		[]string{"/skip"},  // Paths to bypass authMiddleware for
		authMiddleware,     // The authentication middleware
	)
	slog.Info("run: finalHandler created", "type", fmt.Sprintf("%T", finalHandler))

	slog.Info("Starting server", "address", ":8080")
	// http.ListenAndServe expects an http.Handler, which finalHandler is.
	if err := http.ListenAndServe(":8080", finalHandler); err != nil {
		// log.Fatal will call os.Exit(1) after printing,
		// so we use slog for consistency if we want to return the error.
		slog.Error("ListenAndServe failed", "error", err)
		return err // Return the error for main to handle
	}
	return nil
}

func main() {
	// Configure a global structured logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Customize time format for cleaner logs
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(a.Value.Time().Format("15:04:05.000"))
			}
			return a
		},
	}))
	slog.SetDefault(logger)

	if err := run(); err != nil {
		// The error is already logged by run(), so just exit.
		os.Exit(1)
	}
}
