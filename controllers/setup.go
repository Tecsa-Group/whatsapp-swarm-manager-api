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
	router.HandleFunc("/instances-by-server-id/{server_id}", handlers.GetInstancesByServerID).Methods("GET")
	router.HandleFunc("/instances", handlers.CreateInstance).Methods("POST")
	router.HandleFunc("/instances/{id}", handlers.UpdateInstance).Methods("PUT")
	router.HandleFunc("/instances/{id}", handlers.DeleteInstance).Methods("DELETE")

	router.HandleFunc("/servers", handlers.GetAllservers).Methods("GET")
	router.HandleFunc("/servers/{id}", handlers.GetServer).Methods("GET")
	router.HandleFunc("/servers", handlers.CreateServer).Methods("POST")
	router.HandleFunc("/servers/{id}", handlers.UpdateServer).Methods("PUT")
	router.HandleFunc("/servers/{id}", handlers.DeleteServer).Methods("DELETE")

	return router
}
