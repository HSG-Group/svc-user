package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"user-svc/application/query"
	infrahttp "user-svc/infrastructure/http"
	"user-svc/infrastructure/persistence"
)

func main() {
	port := os.Getenv("PORT")
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
