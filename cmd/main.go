package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"fuzzy-eureka_eafonso/config"
	"fuzzy-eureka_eafonso/internal/routes"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	config.InitDB()

	config.InitWorkerPool(10)

	routes.SetupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
