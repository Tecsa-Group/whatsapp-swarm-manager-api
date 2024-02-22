package main

import (
	"fmt"
	"net/http"

	"github.com/felipe-tecsa/whatsapp-swarm-manager-api/controllers"
	"github.com/felipe-tecsa/whatsapp-swarm-manager-api/models"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Erro ao carregar o arquivo .env:", err)
		return
	}

	handler := controllers.New()

	server := &http.Server{
		Addr:    "0.0.0.0:5000",
		Handler: handler,
	}

	// Conectar ao banco de dados
	err = models.ConnectDatabase()
	if err != nil {
		fmt.Println("Erro ao conectar ao banco de dados:", err)
		return
	}

	err = server.ListenAndServe()
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
		return
	}
}
