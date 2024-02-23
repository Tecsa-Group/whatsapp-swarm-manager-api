package handlers

import (
	"bytes"
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

func CreateInstanceEvolution(w http.ResponseWriter, r *http.Request) {
	// Verifique se o método da solicitação é POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Decodifique o corpo da solicitação JSON em uma struct InstanceRequest
	var payload models.InstanceRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Erro ao decodificar o corpo da solicitação JSON", http.StatusBadRequest)
		return
	}

	url := "http://evolution.shub.tech/instance/create"

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Erro ao codificar o payload:"+err.Error(), http.StatusInternalServerError)
		return
	}

	// Cria a solicitação HTTP POST com o payload JSON
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		http.Error(w, "Erro ao criar a solicitação HTTP:"+err.Error(), http.StatusInternalServerError)
		return
	}

	// Define os cabeçalhos da solicitação
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", "e1998e715164141382c8d44434629632")

	// Faz a solicitação HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Erro ao enviar a solicitação HTTP:"+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Lê a resposta da solicitação
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Erro ao ler a resposta da solicitação HTTP:"+err.Error(), http.StatusInternalServerError)
		return
	}

	// Escreve a resposta no corpo da resposta HTTP
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}
