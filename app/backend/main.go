package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb *redis.Client

type Task struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

func initRedis() {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "redis:6379"
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("Failed to connect to Redis at %s: %v", redisHost, err)
	} else {
		log.Printf("Successfully connected to Redis at %s", redisHost)
	}
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	tasks, err := rdb.LRange(ctx, "tasks", 0, -1).Result()
	if err != nil {
		http.Error(w, `{"error": "Failed to fetch tasks"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if len(tasks) == 0 {
		w.Write([]byte("[]"))
		return
	}

	// Create JSON array string
	jsonString := "["
	for i, task := range tasks {
		jsonString += task
		if i < len(tasks)-1 {
			jsonString += ","
		}
	}
	jsonString += "]"
	
	w.Write([]byte(jsonString))
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req map[string]string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
		return
	}

	content, ok := req["content"]
	if !ok || content == "" {
		http.Error(w, `{"error": "Content is required"}`, http.StatusBadRequest)
		return
	}

	task := Task{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Content:   content,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	taskJSON, _ := json.Marshal(task)
	
	err := rdb.RPush(ctx, "tasks", taskJSON).Err()
	if err != nil {
		http.Error(w, `{"error": "Failed to save task"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(taskJSON)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	initRedis()

	http.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getTasksHandler(w, r)
		} else if r.Method == "POST" || r.Method == "OPTIONS" {
			addTaskHandler(w, r)
		} else {
			http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		}
	})
	
	http.HandleFunc("/health", healthHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Backend API server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
