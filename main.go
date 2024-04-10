package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const (
	defaultHost = ""
	defaultPort = "8080"
)

func main() {
	slog.Info("Starting up.")

	// Get the config...
	host := os.Getenv("APP_HOST")
	if host == "" {
		host = defaultHost
	}
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = defaultPort
	}
	addr := host + ":" + port

	// Define the routes...
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Incoming request", "method", r.Method, "path", r.URL.Path, "remote", r.RemoteAddr, "status", http.StatusNotFound)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found\n"))
	})
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Incoming request", "method", r.Method, "path", r.URL.Path, "remote", r.RemoteAddr)
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<h1>Hello, World!</h1>"))
	})
	mux.HandleFunc("GET /ping/{$}", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Incoming request", "method", r.Method, "path", r.URL.Path, "remote", r.RemoteAddr)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true, "msg": "pong"}\n`))
	})
	mux.HandleFunc("POST /echo/{$}", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Incoming request", "method", r.Method, "path", r.URL.Path, "remote", r.RemoteAddr)

		// Read the body...
		b, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("Error reading body", "err", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"success": false, "msg": "error reading request body"}\n`))
			return
		}

		// Try to parse the body as JSON...
		var data map[string]any
		if err := json.Unmarshal(b, &data); err != nil {
			slog.Error("Error parsing JSON", "err", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"success": false, "msg": "error parsing request body as JSON object"}\n`))
			return
		}

		// Encode the response data as json...
		b, err = json.Marshal(map[string]any{
			"success": true,
			"msg":     "echo",
			"data":    data,
		})
		if err != nil {
			slog.Error("Error encoding JSON", "err", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"success": false, "msg": "error encoding response as JSON"}\n`))
			return
		}

		// Write the response...
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	})

	// Start the server...
	slog.Info("Starting server", "addr", addr)
	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
		slog.Error("Error running server", "err", err)
	}
}
