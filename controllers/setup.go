package controllers

import (
	"net/http"

	"github.com/felipe-tecsa/whatsapp-swarm-manager-api/handlers"
	"github.com/gorilla/mux"
)

func New() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/instances", handlers.GetAllInstances).Methods("GET")
	router.HandleFunc("/instances/{id}", handlers.GetInstance).Methods("GET")
	router.HandleFunc("/instances", handlers.CreateInstance).Methods("POST")
	router.HandleFunc("/instances/{id}", handlers.UpdateInstance).Methods("PUT")
	router.HandleFunc("/instances/{id}", handlers.DeleteInstance).Methods("DELETE")

	return router
}
