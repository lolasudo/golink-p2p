package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Service2Response struct {
	Timestamp string   `json:"timestamp"`
	Hostname  string   `json:"hostname"`
	Service   string   `json:"service"`
	Data      []string `json:"data"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]string{
			"service":  "User Service",
			"status":   "running",
			"language": "Golang",
		}
		json.NewEncoder(w).Encode(response)
	})

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		users := []User{
			{ID: 1, Name: "Alice", Email: "alice@example.com"},
			{ID: 2, Name: "Bob", Email: "bob@example.com"},
		}
		json.NewEncoder(w).Encode(users)
	})

	http.HandleFunc("/user-info", func(w http.ResponseWriter, r *http.Request) {
		// Вызов второго сервиса с таймаутом
		client := &http.Client{
			Timeout: 5 * time.Second,
		}

		resp, err := client.Get("http://service2:8081/info")
		if err != nil {
			errorMsg := fmt.Sprintf("Error calling service2: %v", err)
			http.Error(w, errorMsg, http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Error reading response from service2", http.StatusInternalServerError)
			return
		}

		var service2Resp Service2Response
		if err := json.Unmarshal(body, &service2Resp); err != nil {
			http.Error(w, "Error parsing service2 response", http.StatusInternalServerError)
			return
		}

		// Возвращаем комбинированный ответ
		combinedResponse := map[string]interface{}{
			"user_service":        "Service 1 (Go)",
			"info_from_service2":  service2Resp,
			"combined_timestamp":  time.Now().Format(time.RFC3339),
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(combinedResponse)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Service 1 is healthy at %v", time.Now())
	})

	fmt.Println("Service 1 starting on :8080")
	http.ListenAndServe(":8080", nil)
}