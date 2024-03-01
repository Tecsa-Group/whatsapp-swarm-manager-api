package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/felipe-tecsa/whatsapp-swarm-manager-api/models"
	"github.com/felipe-tecsa/whatsapp-swarm-manager-api/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

func HandleProxy(w http.ResponseWriter, r *http.Request) {
	path := removeLastItemAfterLastSlash(string(r.URL.Path))

	switch r.Method {
	case http.MethodGet:
		switch path {
		case "/instance/connectionState":
			ConnectionStateInstanceEvolution(w, r)
		case "/instance/connect":
			ConnectInstanceEvolution(w, r)
		default:
			http.NotFound(w, r)
		}
	case http.MethodPost:
		switch path {
		case "/instance":
			CreateInstanceEvolution(w, r)
		default:
			http.NotFound(w, r)
		}
	case http.MethodDelete:
		switch path {
		case "/instance/logout":
			LogoutInstanceEvolution(w, r)
		case "/instance/delete":
			DeleteInstanceEvolution(w, r)
		default:
			http.NotFound(w, r)
		}
	case http.MethodPut:
		switch path {
		case "/instance/restart":
			LogoutInstanceEvolution(w, r)
		default:
			http.NotFound(w, r)
		}
	default:
		http.Error(w, "Método não suportado", http.StatusMethodNotAllowed)
	}
}

func removeLastItemAfterLastSlash(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) > 1 {
		parts = parts[:len(parts)-1]
	}
	return strings.Join(parts, "/")
}

func getLastItemAfterLastSlash(url string) string {
	parts := strings.Split(url, "/")
	lastItem := parts[len(parts)-1]
	return lastItem
}

func GetAllInstances(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var intances []models.Instance
	models.DB.Find(&intances)

	json.NewEncoder(w).Encode(intances)
}

func GetInstance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	id := mux.Vars(r)["id"]
	var instance models.Instance

	if err := models.DB.Where("id = ?", id).First(&instance).Error; err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Instance not found")
		return
	}

	json.NewEncoder(w).Encode(instance)
}

var validate *validator.Validate

func CreateInstance(input models.Instance) error {
	validate := validator.New()
	err := validate.Struct(input)
	if err != nil {
		return err
	}

	instance := &models.Instance{
		Name:     input.Name,
		Status:   input.Status,
		ServerID: input.ServerID,
		Apikey:   input.Apikey,
	}

	if result := models.DB.Create(instance); result.Error != nil {
		return result.Error
	}

	return nil
}

type UpdateInstanceModel struct {
	Name      string     `json:"name" validate:"omitempty"`
	Status    string     `json:"status" validate:"omitempty"`
	UpdatedAt *time.Time `json:"updated_at" validate:"omitempty"`
}

func UpdateInstance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

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
	w.Header().Set("Access-Control-Allow-Origin", "*")

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
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	serverID := vars["server_id"]

	var instances []models.Instance
	if err := models.DB.Preload("Server").Where("server_id = ?", serverID).Find(&instances).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve instances by server ID")
		return
	}

	json.NewEncoder(w).Encode(instances)
}

func CreateInstanceEvolution(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var payload models.InstanceRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Erro ao decodificar o corpo da solicitação JSON", http.StatusBadRequest)
		return
	}
	serverUrl, serverId := verifyServerAvailability()

	url := serverUrl + r.URL.Path
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Erro ao codificar o payload:"+err.Error(), http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest(r.Method, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		http.Error(w, "Erro ao criar a solicitação HTTP:"+err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", os.Getenv("EVOLUTION_APIKEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Erro ao enviar a solicitação HTTP:"+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Erro ao ler a resposta da solicitação HTTP:"+err.Error(), http.StatusInternalServerError)
		return
	}

	newInstance := models.Instance{
		Name:     payload.InstanceName,
		Status:   "open",
		ServerID: serverId,
		Apikey:   payload.Token,
	}

	err = CreateInstance(newInstance)
	if err != nil {
		http.Error(w, "Erro ao criar a instância: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func DeleteInstanceEvolution(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	instanceName := getLastItemAfterLastSlash(r.URL.Path)

	if instanceName == "" {
		http.Error(w, "Nome da instância não fornecido", http.StatusBadRequest)
		return
	}

	serverUrl, error := findIntanceByIntanceName(instanceName)
	if error != nil {
		http.Error(w, "Servidor não encontrado", http.StatusNotFound)
		return
	}

	url := serverUrl + r.URL.Path
	req, err := http.NewRequest(r.Method, url, nil)
	if err != nil {
		http.Error(w, "Erro ao criar a solicitação HTTP: "+err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Set("apikey", os.Getenv("EVOLUTION_APIKEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Erro ao enviar a solicitação HTTP: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Erro ao ler a resposta da solicitação HTTP: "+err.Error(), http.StatusInternalServerError)
		return
	}

	models.DB.Table("instances").Where("instances.name = ?", instanceName).Delete(&models.Instance{})

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func ConnectionStateInstanceEvolution(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	instanceName := getLastItemAfterLastSlash(r.URL.Path)

	if instanceName == "" {
		http.Error(w, "Nome da instância não fornecido", http.StatusBadRequest)
		return
	}

	serverUrl, error := findIntanceByIntanceName(instanceName)
	if error != nil {
		http.Error(w, "Servidor não encontrado", http.StatusNotFound)
		return
	}

	url := serverUrl + r.URL.Path
	req, err := http.NewRequest(r.Method, url, nil)
	if err != nil {
		http.Error(w, "Erro ao criar a solicitação HTTP: "+err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Set("apikey", os.Getenv("EVOLUTION_APIKEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Erro ao enviar a solicitação HTTP: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Erro ao ler a resposta da solicitação HTTP: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func findIntanceByIntanceName(instanceName string) (string, error) {
	serverUrl := models.UrlServer

	error := models.DB.Table("servers").
		Select("servers.url").
		Joins("LEFT JOIN instances on servers.id = instances.server_id").
		Where("instances.name = ?", instanceName).
		Group("servers.url").
		First(&serverUrl).
		Error

	if error != nil {
		return "", error
	}
	return serverUrl.URL, nil
}

func RestartInstanceEvolution(w http.ResponseWriter, r *http.Request) {
	instanceName := getLastItemAfterLastSlash(r.URL.Path)

	if instanceName == "" {
		http.Error(w, "Nome da instância não fornecido", http.StatusBadRequest)
		return
	}

	serverUrl, error := findIntanceByIntanceName(instanceName)
	if error != nil {
		http.Error(w, "Servidor não encontrado", http.StatusNotFound)
		return
	}

	url := serverUrl + r.URL.Path
	req, err := http.NewRequest(r.Method, url, nil)
	if err != nil {
		http.Error(w, "Erro ao criar a solicitação HTTP: "+err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Set("apikey", os.Getenv("EVOLUTION_APIKEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Erro ao enviar a solicitação HTTP: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Erro ao ler a resposta da solicitação HTTP: "+err.Error(), http.StatusInternalServerError)
		return
	}

	UpdateStatusInstance("name", instanceName, "close")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func LogoutInstanceEvolution(w http.ResponseWriter, r *http.Request) {
	instanceName := getLastItemAfterLastSlash(r.URL.Path)

	if instanceName == "" {
		http.Error(w, "Nome da instância não fornecido", http.StatusBadRequest)
		return
	}

	serverUrl, error := findIntanceByIntanceName(instanceName)
	if error != nil {
		http.Error(w, "Servidor não encontrado", http.StatusNotFound)
		return
	}

	url := serverUrl + r.URL.Path
	req, err := http.NewRequest(r.Method, url, nil)
	if err != nil {
		http.Error(w, "Erro ao criar a solicitação HTTP: "+err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Set("apikey", os.Getenv("EVOLUTION_APIKEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Erro ao enviar a solicitação HTTP: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Erro ao ler a resposta da solicitação HTTP: "+err.Error(), http.StatusInternalServerError)
		return
	}

	UpdateStatusInstance("name", instanceName, "close")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func ConnectInstanceEvolution(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	instanceName := getLastItemAfterLastSlash(r.URL.Path)

	if instanceName == "" {
		http.Error(w, "Nome da instância não fornecido", http.StatusBadRequest)
		return
	}

	serverUrl, error := findIntanceByIntanceName(instanceName)
	if error != nil {
		http.Error(w, "Servidor não encontrado", http.StatusNotFound)
		return
	}

	url := serverUrl + r.URL.Path

	req, err := http.NewRequest(r.Method, url, nil)
	if err != nil {
		http.Error(w, "Erro ao criar a solicitação HTTP: "+err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Set("apikey", os.Getenv("EVOLUTION_APIKEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Erro ao enviar a solicitação HTTP: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Erro ao ler a resposta da solicitação HTTP: "+err.Error(), http.StatusInternalServerError)
		return
	}

	UpdateStatusInstance("name", instanceName, "open")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.WriteHeader(resp.StatusCode)

	w.Write(body)
}

func FetchInstances() {
	fmt.Println("Start Cron")
	var servers []models.Server

	err := models.DB.Table("servers s").
		Distinct("s.url").
		Find(&servers).Error

	if err != nil {
		fmt.Println("Erro ao executar a consulta:", err)
	}

	for _, server := range servers {
		client := &http.Client{}

		req, err := http.NewRequest("GET", server.URL+"/instance/fetchInstances", nil)
		if err != nil {
			fmt.Println("Erro ao criar requisição HTTP:", err)
			continue // Continue para o próximo servidor em caso de erro na requisição
		}
		req.Header.Set("apikey", os.Getenv("EVOLUTION_APIKEY"))

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Erro ao fazer requisição HTTP:", err)
			continue // Continue para o próximo servidor em caso de erro na requisição
		}
		defer resp.Body.Close()

		var instances []models.ServerInstance
		if err := json.NewDecoder(resp.Body).Decode(&instances); err != nil {
			fmt.Println("Erro ao decodificar resposta JSON:", err)
			continue // Continue para o próximo servidor em caso de erro na decodificação JSON
		}
		fmt.Print("instance", instances)

		for _, instance := range instances {
			err := UpdateStatusInstance("apikey", instance.Instance.ApiKey, instance.Instance.Status)
			if err != nil {
				fmt.Println("Erro ao atualizar status da instância:", err)
				// Você pode decidir continuar para o próximo servidor ou retornar o erro, dependendo dos requisitos do seu aplicativo
			}
		}
	}
	fmt.Println("End Cron")
}

func DeleteAllInstances() {
	var servers []models.Server

	err := models.DB.Table("servers s").
		Distinct("s.url").
		Find(&servers).Error

	if err != nil {
		fmt.Println("Erro ao executar a consulta:", err)
	}

	for _, server := range servers {
		client := &http.Client{}

		req, err := http.NewRequest("GET", server.URL+"/instance/fetchInstances", nil)
		if err != nil {
			fmt.Println("Erro ao criar requisição HTTP:", err)
			continue // Continue para o próximo servidor em caso de erro na requisição
		}
		fmt.Print("url", server.URL)
		req.Header.Set("apikey", os.Getenv("EVOLUTION_APIKEY"))

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Erro ao fazer requisição HTTP:", err)
			continue // Continue para o próximo servidor em caso de erro na requisição
		}
		defer resp.Body.Close()

		var instances []models.ServerInstance
		if err := json.NewDecoder(resp.Body).Decode(&instances); err != nil {
			fmt.Println("Erro ao decodificar resposta JSON:", err)
			go FetchInstances()
			return // Continue para o próximo servidor em caso de erro na decodificação JSON
		}

		for _, instance := range instances {
			replaced := strings.Replace(instance.Instance.InstanceName, " ", "%20", -1)
			fmt.Println("replaced", replaced)
			url := server.URL + "/instance/delete/" + replaced
			req, _ := http.NewRequest("DELETE", url, nil)

			req.Header.Set("apikey", os.Getenv("EVOLUTION_APIKEY"))

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("erro no delete request", err)
			}
			defer resp.Body.Close()
		}
	}
}

func UpdateStatusInstance(field string, fieldValue string, status string) error {
	var instance models.Instance
	result := models.DB.Where(field+" = ?", fieldValue).First(&instance)
	if result.Error != nil {
		return result.Error
	}

	instance.Status = status

	if err := models.DB.Save(&instance).Error; err != nil {
		return err
	}

	return nil
}

func verifyServerAvailability() (string, int) {
	var servers []models.Result

	err := models.DB.Table("servers").
		Select("servers.url, servers.id, COUNT(instances.id) AS count_open").
		Joins("LEFT JOIN instances ON servers.id = instances.server_id AND instances.status = ?", "open").
		Group("servers.url, servers.id").
		Order("count_open DESC").
		Scan(&servers).Error

	if err != nil {
		fmt.Println("Erro ao executar a consulta:", err)
	}
	isFull := false
	for _, server := range servers {
		if server.CountOpen == 20 {
			isFull = true
		} else {
			isFull = false
		}
	}

	if !isFull {
		for _, server := range servers {
			// if server.CountOpen == 10 {
			// 	go CreateServerHetzner()
			// }
			if server.CountOpen < 20 {
				return server.URL, server.ID
			}
		}
	} else {
		models.DB.Table("instances").Delete(&models.Instance{})
		go DeleteAllInstances()
	}

	return servers[0].URL, servers[0].ID
}
