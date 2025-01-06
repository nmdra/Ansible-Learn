package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Response struct {
	Message string `json:"message"`
}

func hello(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := Response{Message: "Hello, World!"}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		// If an error occurs during encoding, send an internal server error
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// headers function to handle the "/headers" route
func headers(w http.ResponseWriter, req *http.Request) {
	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Create a map to hold the headers
	headersMap := make(map[string][]string)

	// Populate the map with the headers from the incoming request
	for name, headers := range req.Header {
		headersMap[name] = headers
	}

	// Encode the headers map into JSON and send it to the client
	err := json.NewEncoder(w).Encode(headersMap)
	if err != nil {
		// If an error occurs during encoding, send an internal server error
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// loggingMiddleware logs each incoming request
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", req.Method, req.RequestURI)
		next.ServeHTTP(w, req)
		log.Printf("Completed %s %s in %v", req.Method, req.RequestURI, time.Since(start))
	})
}

func main() {
	PORT := "8080"

	if PORT == "" {
		log.Println("Port is not set")
		PORT = "8081"
	}

	// Handle the /hello route
	http.Handle("/hello", loggingMiddleware(http.HandlerFunc(hello)))

	// Handle the /headers route
	http.Handle("/headers", loggingMiddleware(http.HandlerFunc(headers)))

	// Start the HTTP server and listen on the specified port
	fmt.Println("Server starting on " + PORT)

	err := http.ListenAndServe(":"+PORT, nil)
	if err != nil {
		// If the server fails to start, print the error and exit
		log.Fatal("Error starting server:", err)
	}
}
