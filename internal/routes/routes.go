package routes

import (
	"net/http"

	"fuzzy-eureka_eafonso/internal/handlers"
)

func SetupRoutes() {
	http.HandleFunc("/request", handlers.RequestHandler)
}
