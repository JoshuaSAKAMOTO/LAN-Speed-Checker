package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	defaultServerURL   = "http://localhost:8080"
	parallelConnections = 4 // Number of parallel connections for testing
)

// printProgress prints a simple progress indicator
func printProgress(message string, done chan bool) {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	spinChars := []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
	i := 0

	for {
		select {
		case <-done:
			fmt.Print("\r" + message + " ✓\n")
			return
		case <-ticker.C:
			fmt.Printf("\r%s %s", message, spinChars[i%len(spinChars)])
			i++
		}
	}
}

// measureDownloadSpeed measures the download speed from the server
func measureDownloadSpeed(serverURL string) (float64, error) {
	url := serverURL + "/download"

	done := make(chan bool)
	go printProgress("  Downloading test data...", done)

	// Start timer
	startTime := time.Now()

	// Make HTTP request
	resp, err := http.Get(url)
	if err != nil {
		close(done)
		return 0, fmt.Errorf("failed to connect to server: %w", err)
	}
	defer resp.Body.Close()

	// Read all data
	totalBytes := int64(0)
	buffer := make([]byte, 32*1024) // 32KB buffer

	for {
		n, err := resp.Body.Read(buffer)
		totalBytes += int64(n)

		if err != nil {
			if err == io.EOF {
				break
			}
			close(done)
			return 0, fmt.Errorf("error reading data: %w", err)
		}
	}

	// Calculate duration and speed
	duration := time.Since(startTime)
	speedMbps := (float64(totalBytes) * 8) / duration.Seconds() / 1000000 // Convert to Mbps

	close(done)
	time.Sleep(100 * time.Millisecond) // Wait for progress indicator to finish

	return speedMbps, nil
}

// measureUploadSpeed measures the upload speed to the server
func measureUploadSpeed(serverURL string) (float64, error) {
	url := serverURL + "/upload"

	done := make(chan bool)
	go printProgress("  Uploading test data...", done)

	// Create data to upload (50MB)
	dataSize := 50 * 1024 * 1024

	// Start timer
	startTime := time.Now()

	// Create HTTP request
	resp, err := http.Post(url, "application/octet-stream", &io.LimitedReader{
		R: &infiniteReader{},
		N: int64(dataSize),
	})
	if err != nil {
		close(done)
		return 0, fmt.Errorf("failed to upload data: %w", err)
	}
	defer resp.Body.Close()

	// Calculate duration and speed
	duration := time.Since(startTime)
	speedMbps := (float64(dataSize) * 8) / duration.Seconds() / 1000000 // Convert to Mbps

	close(done)
	time.Sleep(100 * time.Millisecond) // Wait for progress indicator to finish

	return speedMbps, nil
}

// infiniteReader implements io.Reader that returns zero bytes indefinitely
type infiniteReader struct{}

func (r *infiniteReader) Read(p []byte) (n int, err error) {
	return len(p), nil
}

// measureDownloadSpeedParallel measures download speed using multiple parallel connections
func measureDownloadSpeedParallel(serverURL string, connections int) (float64, error) {
	done := make(chan bool)
	go printProgress(fmt.Sprintf("  Downloading with %d parallel connections...", connections), done)

	var wg sync.WaitGroup
	speedChan := make(chan float64, connections)
	errorChan := make(chan error, connections)

	// Start timer
	startTime := time.Now()

	// Launch parallel downloads
	for i := 0; i < connections; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			url := serverURL + "/download"
			resp, err := http.Get(url)
			if err != nil {
				errorChan <- fmt.Errorf("connection %d failed: %w", id, err)
				return
			}
			defer resp.Body.Close()

			// Read all data
			totalBytes := int64(0)
			buffer := make([]byte, 32*1024)

			for {
				n, err := resp.Body.Read(buffer)
				totalBytes += int64(n)

				if err != nil {
					if err == io.EOF {
						break
					}
					errorChan <- fmt.Errorf("connection %d read error: %w", id, err)
					return
				}
			}
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(speedChan)
	close(errorChan)

	// Check for errors
	if len(errorChan) > 0 {
		close(done)
		return 0, <-errorChan
	}

	// Calculate total duration and speed
	duration := time.Since(startTime)

	// Approximate total data transferred (100MB per connection)
	totalBytes := int64(connections) * 100 * 1024 * 1024
	speedMbps := (float64(totalBytes) * 8) / duration.Seconds() / 1000000

	close(done)
	time.Sleep(100 * time.Millisecond) // Wait for progress indicator to finish

	return speedMbps, nil
}

// measureUploadSpeedParallel measures upload speed using multiple parallel connections
func measureUploadSpeedParallel(serverURL string, connections int) (float64, error) {
	done := make(chan bool)
	go printProgress(fmt.Sprintf("  Uploading with %d parallel connections...", connections), done)

	var wg sync.WaitGroup
	errorChan := make(chan error, connections)

	// Data size per connection
	dataSize := 50 * 1024 * 1024 // 50MB per connection

	// Start timer
	startTime := time.Now()

	// Launch parallel uploads
	for i := 0; i < connections; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			url := serverURL + "/upload"
			resp, err := http.Post(url, "application/octet-stream", &io.LimitedReader{
				R: &infiniteReader{},
				N: int64(dataSize),
			})
			if err != nil {
				errorChan <- fmt.Errorf("connection %d failed: %w", id, err)
				return
			}
			defer resp.Body.Close()
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(errorChan)

	// Check for errors
	if len(errorChan) > 0 {
		close(done)
		return 0, <-errorChan
	}

	// Calculate total duration and speed
	duration := time.Since(startTime)
	totalBytes := int64(connections) * int64(dataSize)
	speedMbps := (float64(totalBytes) * 8) / duration.Seconds() / 1000000

	close(done)
	time.Sleep(100 * time.Millisecond) // Wait for progress indicator to finish

	return speedMbps, nil
}

func main() {
	serverURL := defaultServerURL

	// Print header
	fmt.Println("\n╔═══════════════════════════════════╗")
	fmt.Println("║   LAN Speed Tester Client        ║")
	fmt.Println("╚═══════════════════════════════════╝")
	fmt.Printf("Server: %s\n\n", serverURL)

	// Test server connectivity
	done := make(chan bool)
	go printProgress("Connecting to server...", done)
	resp, err := http.Get(serverURL)
	if err != nil {
		close(done)
		log.Fatalf("\n✗ Failed to connect to server: %v", err)
	}
	resp.Body.Close()
	close(done)
	time.Sleep(100 * time.Millisecond)
	fmt.Println()

	// Single connection tests
	fmt.Println("┌─────────────────────────────────┐")
	fmt.Println("│  Single Connection Test         │")
	fmt.Println("└─────────────────────────────────┘")

	// Measure download speed (single connection)
	downloadSpeedSingle, err := measureDownloadSpeed(serverURL)
	if err != nil {
		log.Fatalf("✗ Download test failed: %v", err)
	}
	fmt.Printf("  Download: %.2f Mbps\n\n", downloadSpeedSingle)

	// Measure upload speed (single connection)
	uploadSpeedSingle, err := measureUploadSpeed(serverURL)
	if err != nil {
		log.Fatalf("✗ Upload test failed: %v", err)
	}
	fmt.Printf("  Upload:   %.2f Mbps\n\n", uploadSpeedSingle)

	// Parallel connection tests
	fmt.Println("┌─────────────────────────────────┐")
	fmt.Printf("│  Parallel Test (%d connections) │\n", parallelConnections)
	fmt.Println("└─────────────────────────────────┘")

	// Measure download speed (parallel)
	downloadSpeedParallel, err := measureDownloadSpeedParallel(serverURL, parallelConnections)
	if err != nil {
		log.Fatalf("✗ Parallel download test failed: %v", err)
	}
	fmt.Printf("  Download: %.2f Mbps\n\n", downloadSpeedParallel)

	// Measure upload speed (parallel)
	uploadSpeedParallel, err := measureUploadSpeedParallel(serverURL, parallelConnections)
	if err != nil {
		log.Fatalf("✗ Parallel upload test failed: %v", err)
	}
	fmt.Printf("  Upload:   %.2f Mbps\n\n", uploadSpeedParallel)

	// Display results summary
	fmt.Println("╔═══════════════════════════════════╗")
	fmt.Println("║      Test Results Summary        ║")
	fmt.Println("╠═══════════════════════════════════╣")
	fmt.Println("║ Single Connection:               ║")
	fmt.Printf("║   Download: %9.2f Mbps      ║\n", downloadSpeedSingle)
	fmt.Printf("║   Upload:   %9.2f Mbps      ║\n", uploadSpeedSingle)
	fmt.Println("║                                  ║")
	fmt.Printf("║ Parallel (%d connections):       ║\n", parallelConnections)
	fmt.Printf("║   Download: %9.2f Mbps      ║\n", downloadSpeedParallel)
	fmt.Printf("║   Upload:   %9.2f Mbps      ║\n", uploadSpeedParallel)
	fmt.Println("╚═══════════════════════════════════╝\n")
}
