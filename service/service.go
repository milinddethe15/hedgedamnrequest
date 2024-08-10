package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// logRequest is a helper function to log request details
func logRequest(r *http.Request, handlerName string) {
	log.Printf("[%s] %s %s", handlerName, r.Method, r.URL)
	defer func() {
		log.Printf("[%s] %s %s", handlerName, r.Method, r.URL)
	}()
}

func fastHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r, "fastHandler")
	// Responds quickly
	fmt.Fprintf(w, "Response from fast service")
}

func slowHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r, "slowHandler")
	// Simulates a slow response
	time.Sleep(32 * time.Second)
	fmt.Fprintf(w, "Response from slow service")
}

func mediumHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r, "mediumHandler")
	// Simulates a medium latency response
	time.Sleep(10 * time.Second)
	fmt.Fprintf(w, "Response from medium-latency service")
}

func main() {
	// Initialize logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Setting up the services on different ports
	http.HandleFunc("/fast", fastHandler)
	http.HandleFunc("/slow", slowHandler)
	http.HandleFunc("/medium", mediumHandler)

	// Start services
	fmt.Println("Starting services on :8081 (fast), :8082 (slow), :8083 (medium)")

	go func() {
		if err := http.ListenAndServe(":8081", nil); err != nil {
			log.Fatalf("Error starting fast service: %v", err)
		}
	}()
	go func() {
		if err := http.ListenAndServe(":8082", nil); err != nil {
			log.Fatalf("Error starting slow service: %v", err)
		}
	}()
	if err := http.ListenAndServe(":8083", nil); err != nil {
		log.Fatalf("Error starting medium service: %v", err)
	}
}
