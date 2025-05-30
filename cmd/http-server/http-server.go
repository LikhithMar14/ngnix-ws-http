package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "pong", "timestamp": "%s"}`, time.Now().Format(time.RFC3339))
	})


	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			name = "World"
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Hello, %s!", "timestamp": "%s"}`, name, time.Now().Format(time.RFC3339))
	})


	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status": "HTTP server running", "port": 8080, "timestamp": "%s"}`, time.Now().Format(time.RFC3339))
	})


	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "HTTP Server is running on port 8080!")
	})

	fmt.Println("HTTP Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}