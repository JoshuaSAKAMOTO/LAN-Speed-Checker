package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
)

const (
	defaultPort = "8080"
)

// downloadHandler handles download speed test requests
func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// Set headers
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Cache-Control", "no-cache")

	// Generate random data (100MB)
	dataSize := 100 * 1024 * 1024 // 100MB
	buffer := make([]byte, 1024*1024) // 1MB buffer

	written := 0
	for written < dataSize {
		// Generate random data
		if _, err := rand.Read(buffer); err != nil {
			log.Printf("Error generating random data: %v", err)
			return
		}

		// Write to response
		n, err := w.Write(buffer)
		if err != nil {
			log.Printf("Error writing data: %v", err)
			return
		}
		written += n
	}

	log.Printf("Download completed: %d bytes sent", written)
}

// uploadHandler handles upload speed test requests
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read all uploaded data
	buffer := make([]byte, 1024*1024) // 1MB buffer
	totalReceived := 0

	for {
		n, err := r.Body.Read(buffer)
		totalReceived += n
		if err != nil {
			break
		}
	}

	log.Printf("Upload completed: %d bytes received", totalReceived)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Received %d bytes", totalReceived)
}

// healthHandler handles health check requests
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "LAN Speed Tester Server is running")
}

func main() {
	// Register handlers
	http.HandleFunc("/", healthHandler)
	http.HandleFunc("/download", downloadHandler)
	http.HandleFunc("/upload", uploadHandler)

	// Start server
	addr := ":" + defaultPort
	log.Printf("Starting LAN Speed Tester Server on %s", addr)
	log.Printf("Endpoints:")
	log.Printf("  - Health Check: http://localhost%s/", addr)
	log.Printf("  - Download Test: http://localhost%s/download", addr)
	log.Printf("  - Upload Test: http://localhost%s/upload", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
