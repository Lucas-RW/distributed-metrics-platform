package main 

import (
	"log"
	"net/http"
	"github.com/Lucas-RW/distributed-metrics-platform/internal/handlers"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/metrics", handlers.MetricsHandler)

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", mux)

	if err != nil { 
		log.Fatal(err)
	}
}