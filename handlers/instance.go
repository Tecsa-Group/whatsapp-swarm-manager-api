package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
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

func CreateInstance(input models.Instance) error {
	// Validação de entrada
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

var urlServer = "http://evolution.shub.tech"

func CreateInstanceEvolution(w http.ResponseWriter, r *http.Request) {
	// err := verifyServerAvailabity()
	// if err != nil {
	// 	http.Error(w, "Erro ao verificar a disponibilidade do servidor: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }

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

	url := verifyServerAvailabity() + "/instance/create"
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
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Erro ao ler a resposta da solicitação HTTP:"+err.Error(), http.StatusInternalServerError)
		return
	}
	envValueStr := os.Getenv("ID_SERVER")
	envValueInt, err := strconv.Atoi(envValueStr)
	if err != nil {
		fmt.Println("Erro ao converter variável de ambiente para inteiro:", err)
	}

	newInstance := models.Instance{
		Name:     payload.InstanceName,
		Status:   "connecting",
		ServerID: envValueInt,
		Apikey:   payload.Token,
	}

	err = CreateInstance(newInstance)
	if err != nil {
		http.Error(w, "Erro ao criar a instância: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Escreve a resposta no corpo da resposta HTTP
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func ConnectInstanceEvolution(w http.ResponseWriter, r *http.Request) {
	// Verifique se o método da solicitação é GET
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtenha o nome da instância do path da solicitação
	instanceName := mux.Vars(r)["instanceName"]

	// Verifique se o nome da instância não está vazio
	if instanceName == "" {
		http.Error(w, "Nome da instância não fornecido", http.StatusBadRequest)
		return
	}

	url := verifyServerAvailabity() + ":8080/instance/connect/" + instanceName

	// Cria a solicitação HTTP GET
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, "Erro ao criar a solicitação HTTP: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Define o cabeçalho da chave de API
	req.Header.Set("apikey", "e1998e715164141382c8d44434629632")

	// Faz a solicitação HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Erro ao enviar a solicitação HTTP: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Lê o corpo da resposta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Erro ao ler a resposta da solicitação HTTP: "+err.Error(), http.StatusInternalServerError)
		return
	}

	UpdateStatusInstance(instanceName, "open")

	// Define o cabeçalho Content-Type na resposta
	w.Header().Set("Content-Type", "application/json")

	// Define o código de status da resposta HTTP com base no status da resposta recebida
	w.WriteHeader(resp.StatusCode)

	// Escreve o corpo da resposta no corpo da resposta HTTP
	w.Write(body)
}

func FetchInstances() ([]models.ServerInstance, error) {
	var servers []models.Server

	err := models.DB.Table("servers s").
		Distinct("s.url").
		Find(&servers).Error

	if err != nil {
		// Lidar com o erro
		fmt.Println("Erro ao executar a consulta:", err)
		return nil, err
	}

	var allInstances []models.ServerInstance
	fmt.Print("servers", servers)

	for _, server := range servers {
		client := &http.Client{}

		// Crie uma solicitação HTTP GET
		req, err := http.NewRequest("GET", server.URL+"/instance/fetchInstances", nil)
		if err != nil {
			return nil, err
		}

		// Adicione o cabeçalho "apikey" à solicitação
		req.Header.Set("apikey", "e1998e715164141382c8d44434629632")

		// Faça a solicitação HTTP
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		// Decodifique os dados JSON da resposta
		var instances []models.ServerInstance
		if err := json.NewDecoder(resp.Body).Decode(&instances); err != nil {
			return nil, err
		}

		// Atualize o status das instâncias
		for _, instance := range instances {
			fmt.Print("instance", instance)
			err := UpdateStatusInstance(instance.APIKey, instance.Status)
			if err != nil {
				// Lidar com o erro, se necessário
				fmt.Println("Erro ao atualizar status da instância:", err)
			}
		}

		// Adicione as instâncias obtidas ao slice allInstances
		allInstances = append(allInstances, instances...)
	}

	return allInstances, nil
}

func UpdateStatusInstance(field string, status string) error {
	var instance models.Instance
	result := models.DB.First(&instance, field)
	if result.Error != nil {
		return result.Error
	}

	instance.Status = status

	if err := models.DB.Save(&instance).Error; err != nil {
		return err
	}

	return nil
}

func verifyServerAvailabity() string {
	var servers []models.Result

	// Realiza a consulta utilizando o GORM
	err := models.DB.Table("servers s").
		Select("s.url, COUNT(i.id) AS count_open").
		Joins("LEFT JOIN instances i ON s.id = i.server_id").
		Where("i.status = ?", "open").
		Group("s.url").
		Scan(&servers).Error

	if err != nil {
		// Lidar com o erro
		fmt.Println("Erro ao executar a consulta:", err)
	}

	// if len(servers) <= 0 {
	// 	return urlServer
	// }
	for _, server := range servers {
		if server.CountOpen <= 20 {
			return server.URL
		}
	}
	go CreateServerHetzner()

	return urlServer
}

// proximos passos criar um metodo q consome a api do guilherme
// integrar com o metodo acima
// utilizar o veridy em todos os metodos que utilizam a url do servidor do evolution
// utilizar o iota para nomer os servidores
// testar toda a api

// responseCreateServer, err := CreateServerHetzner()
// if err != nil {
// 	return "error at CreateServerHetzner: " + err.Error()
// }
// // cria o servidor e com os valores retornados eu preencho e insiro no banco
// serverIdGlobal = responseCreateServer.ID
// return responseCreateServer.URL
