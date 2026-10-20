package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"user-svc/application/query"
	infrahttp "user-svc/infrastructure/http"
	"user-svc/infrastructure/persistence"
)

func main() {
	// port := os.Getenv("PORT")
	// Load .env file BEFORE reading env vars
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env vars")
	}

	port := os.Getenv("PORT")
	log.Println("PORT from env file: ", port)

	if port == "" {
		port = "8083"
	}

	// 1. Initialize Infrastructure adapters
	helloRepo := persistence.NewMemoryHelloRepo()

	// 2. Initialize Application use cases
	getHello := query.NewGetHelloHandler(helloRepo)

	// 3. Initialize HTTP handlers
	helloHandler := infrahttp.NewHelloHandler(getHello)

	// 4. Register Routes
	http.Handle("/hello", helloHandler)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	fmt.Printf("User Service is starting on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
