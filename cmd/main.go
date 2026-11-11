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

	//------------ STUDENT HANDLER ----------------
	// 1. Initialize Infrastructure adapters
	studentRepo := persistence.NewMemoryStudentRepo()

	// 2. Initialize Application use cases
	listStudents := query.NewListStudentsHandler(studentRepo)
	getStudentByID := query.NewGetStudentByIDHandler(studentRepo)
	listStudentCourses := query.NewListStudentCoursesHandler(studentRepo)

	// 3. Initialize HTTP handlers
	studentHandler := infrahttp.NewStudentHandler(listStudents, getStudentByID, listStudentCourses)

	// 4. Build the router (all routes registered inside)
	router := infrahttp.NewRouterWithMux(infrahttp.Handlers{
		Hello:   helloHandler,
		Student: studentHandler,
	}, infrahttp.FakeValidator{})

	fmt.Printf("User Service is starting on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
