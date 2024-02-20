package controllers

import (
	"net/http"

	"github.com/felipe-tecsa/whatsapp-swarm-manager-api/handlers"
	"github.com/gorilla/mux"
)

func New() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/quests", handlers.GetAllInstances).Methods("GET")
	router.HandleFunc("/quest/{id}", handlers.GetInstance).Methods("GET")
	router.HandleFunc("/quest", handlers.CreateInstance).Methods("POST")
	router.HandleFunc("/quest/{id}", handlers.UpdateInstance).Methods("PUT")
	router.HandleFunc("/quest/{id}", handlers.DeleteInstance).Methods("DELETE")

	return router
}
