package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]string{
			"service":  "Info Service",
			"status":   "running", 
			"language": "Golang",
		}
		json.NewEncoder(w).Encode(response)
	})

	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		hostname, _ := os.Hostname()
		info := map[string]interface{}{
			"timestamp": time.Now().Format(time.RFC3339),
			"hostname":  hostname,
			"service":   "Info Service",
			"data":      []string{"item1", "item2", "item3"},
		}
		json.NewEncoder(w).Encode(info)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Service 2 is healthy at %v", time.Now())
	})

	fmt.Println("Service 2 starting on :8081")
	http.ListenAndServe(":8081", nil)
}