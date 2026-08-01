package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type Response struct {
	Message     string    `json:"message"`
	Environment string    `json:"environment"`
	Timestamp   time.Time `json:"timestamp"`
	Version     string    `json:"version"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	response := Response{
		Message:     "Cloud-Native DevOps Blueprint API is running!",
		Environment: env,
		Timestamp:   time.Now(),
		Version:     "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/v1/status", apiHandler)

	log.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
