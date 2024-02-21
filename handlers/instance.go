package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/felipe-tecsa/whatsapp-swarm-manager-api/models"
	"github.com/felipe-tecsa/whatsapp-swarm-manager-api/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

func GetAllInstances(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var intances []models.Instance
	models.DB.Find(&intances)

	json.NewEncoder(w).Encode(intances)
}

func GetInstance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	var instance models.Instance

	if err := models.DB.Where("id = ?", id).First(&instance).Error; err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Instance not found")
		return
	}

	json.NewEncoder(w).Encode(instance)
}

var validate *validator.Validate

func CreateInstance(w http.ResponseWriter, r *http.Request) {
	var input models.Instance

	body, _ := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(body, &input)

	validate = validator.New()
	err := validate.Struct(input)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Validation Error")
		return
	}

	instance := &models.Instance{
		Name:     input.Name,
		Status:   input.Status,
		ServerID: input.ServerID,
	}

	models.DB.Create(instance)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(instance)

}

type UpdateInstanceModel struct {
	Name      string     `json:"name" validate:"omitempty"`
	Status    string     `json:"status" validate:"omitempty"`
	UpdatedAt *time.Time `json:"updated_at" validate:"omitempty"`
}

func UpdateInstance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	var instance models.Instance

	if err := models.DB.Where("id = ?", id).First(&instance).Error; err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Instance not found")
		return
	}

	var input UpdateInstanceModel

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
		instance.Name = input.Name
	}
	if input.Status != "" {
		instance.Status = input.Status
	}
	if input.UpdatedAt != nil {
		instance.UpdatedAt = input.UpdatedAt
	}

	if err := models.DB.Save(&instance).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update instance")
		return
	}

	json.NewEncoder(w).Encode(instance)
}

func DeleteInstance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	var instance models.Instance

	if err := models.DB.Where("id = ?", id).First(&instance).Error; err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "instance not found")
		return
	}

	models.DB.Delete(&instance)

	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(instance)
}

func GetInstancesByServerID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extrai o ID do servidor da solicitação
	vars := mux.Vars(r)
	serverID := vars["server_id"]

	// Realiza a consulta no banco de dados para obter instâncias com o ID do servidor fornecido
	var instances []models.Instance
	if err := models.DB.Preload("Server").Where("server_id = ?", serverID).Find(&instances).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve instances by server ID")
		return
	}

	// Responde com as instâncias encontradas
	json.NewEncoder(w).Encode(instances)
}
