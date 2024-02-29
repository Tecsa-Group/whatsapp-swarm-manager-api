package controllers

import (
	"net/http"

	"github.com/felipe-tecsa/whatsapp-swarm-manager-api/handlers"
	"github.com/gorilla/mux"
)

func New() http.Handler {

	router := mux.NewRouter()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")

	router.PathPrefix("/").HandlerFunc(handlers.HandleProxy)
	http.Handle("/", router)

	// router.HandleFunc("/instances", handlers.GetAllInstances).Methods("GET")
	// router.HandleFunc("/instances/{id}", handlers.GetInstance).Methods("GET")
	// router.HandleFunc("/instances-by-server-id/{server_id}", handlers.GetInstancesByServerID).Methods("GET")
	// router.HandleFunc("/instances/{id}", handlers.UpdateInstance).Methods("PUT")
	// router.HandleFunc("/instances/{id}", handlers.DeleteInstance).Methods("DELETE")

	// router.HandleFunc("/create-instance", handlers.CreateInstanceEvolution).Methods("POST")
	// router.HandleFunc("/connect-instance/{instanceName}", handlers.ConnectInstanceEvolution).Methods("GET")
	// router.HandleFunc("/delete-instance/{instanceName}", handlers.DeleteInstanceEvolution).Methods("DELETE")
	// router.HandleFunc("/logout-instance/{instanceName}", handlers.LogoutInstanceEvolution).Methods("DELETE")
	// router.HandleFunc("/teste/connect-state/{instanceName}", handlers.ConnectionStateInstanceEvolution).Methods("GET")

	// router.HandleFunc("/servers", handlers.GetAllservers).Methods("GET")
	// router.HandleFunc("/servers/{id}", handlers.GetServer).Methods("GET")
	// router.HandleFunc("/servers/{id}", handlers.UpdateServer).Methods("PUT")
	// router.HandleFunc("/servers/{id}", handlers.DeleteServer).Methods("DELETE")
	return router
}
