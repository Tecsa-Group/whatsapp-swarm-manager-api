package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/felipe-tecsa/whatsapp-swarm-manager-api/models"
	"github.com/felipe-tecsa/whatsapp-swarm-manager-api/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

func GetAllservers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var servers []models.Server
	models.DB.Find(&servers)

	json.NewEncoder(w).Encode(servers)
}

func GetServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	var server models.Server

	if err := models.DB.Where("id = ?", id).First(&server).Error; err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Instance not found")
		return
	}

	json.NewEncoder(w).Encode(server)
}

func CreateServer(w http.ResponseWriter, r *http.Request) {
	var input models.Server

	body, _ := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(body, &input)

	validate = validator.New()
	err := validate.Struct(input)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Validation Error")
		return
	}

	server := &models.Server{
		Name:             input.Name,
		IP:               input.IP,
		Port:             input.Port,
		Active:           input.Active,
		URL:              input.URL,
		InstanceQuantity: input.InstanceQuantity,
	}

	models.DB.Create(server)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(server)

}

type UpdateServerModel struct {
	Name             string `json:"name" validate:"omitempty"`
	IP               string `json:"ip" validate:"omitempty"`
	Port             int    `json:"port" validate:"omitempty"`
	Active           bool   `json:"active" validate:"omitempty"`
	URL              string `json:"url" validate:"omitempty"`
	InstanceQuantity int    `json:"instance_quantity" validate:"omitempty"`
}

func UpdateServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	var server models.Server

	if err := models.DB.Where("id = ?", id).First(&server).Error; err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Instance not found")
		return
	}

	var input UpdateServerModel

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := json.Unmarshal(body, &input); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	validate = validator.New()
	if err := validate.Struct(input); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Validation Error")
		return
	}

	// Atualize apenas os campos que não estão vazios na entrada
	if input.Name != "" {
		server.Name = input.Name
	}
	if input.IP != "" {
		server.IP = input.IP
	}
	if input.Active {
		fmt.Print(input)
		server.Active = input.Active
	}
	if input.URL != "" {
		server.URL = input.URL
	}
	if input.InstanceQuantity != 0 {
		server.InstanceQuantity = input.InstanceQuantity
	}

	if err := models.DB.Save(&server).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update server")
		return
	}

	json.NewEncoder(w).Encode(server)
}

func DeleteServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	var server models.Server

	if err := models.DB.Where("id = ?", id).First(&server).Error; err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "server not found")
		return
	}

	models.DB.Delete(&server)

	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(server)
}
