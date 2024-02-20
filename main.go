package main

import (
	"net/http"

	"github.com/felipe-tecsa/whatsapp-swarm-manager-api/controllers"
	"github.com/felipe-tecsa/whatsapp-swarm-manager-api/models"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	handler := controllers.New()

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: handler,
	}

	models.ConnectDatabase()

	server.ListenAndServe()
}
