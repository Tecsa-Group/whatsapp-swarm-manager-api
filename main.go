package main

import (
	"fmt"
	"net/http"

	"github.com/felipe-tecsa/whatsapp-swarm-manager-api/controllers"
	"github.com/felipe-tecsa/whatsapp-swarm-manager-api/handlers"
	"github.com/felipe-tecsa/whatsapp-swarm-manager-api/models"
	"github.com/joho/godotenv"
	"gopkg.in/robfig/cron.v2"
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

	// Configurar e iniciar a cron job
	c := cron.New()
	c.AddFunc("@every 2m", func() {
		// Chamada da função FetchInstances a cada 10 minutos
		handlers.FetchInstances()
	})
	c.Start()

	// Iniciar o servidor HTTP
	err = server.ListenAndServe()
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
		return
	}
}
