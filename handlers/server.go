package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"

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

func CreateServer(input models.Server) error {
	validate := validator.New()
	err := validate.Struct(input)
	if err != nil {
		return err
	}

	server := &models.Server{
		Name: input.Name,
		IP:   input.IP,
		URL:  input.URL,
	}

	models.DB.Create(server)

	return nil
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
	if input.URL != "" {
		server.URL = input.URL
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

type ServerCreationPayload struct {
	Firewalls        []Firewall `json:"firewalls"`
	Image            string     `json:"image"`
	Name             string     `json:"name"`
	ServerType       string     `json:"server_type"`
	SSHKeys          []int      `json:"ssh_keys"`
	StartAfterCreate bool       `json:"start_after_create"`
}

type Firewall struct {
	Firewall int `json:"firewall"`
}

func CreateServerHetzner() (models.Server, error) {
	url := "https://api.hetzner.cloud/v1/servers"
	now := time.Now()
	nameServer := fmt.Sprintf("eapi%s", now.Format("20060102150405"))

	payload := ServerCreationPayload{
		Firewalls: []Firewall{
			{Firewall: 1205992},
			{Firewall: 1207630},
		},
		Image:      "ubuntu-22.04",
		Name:       nameServer,
		ServerType: "cx11",
		SSHKeys: []int{
			19697489, 19698323,
		},
		StartAfterCreate: true,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return models.Server{}, err
	} // Cria a solicitação HTTP POST com o payload JSON

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return models.Server{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer 5dZtXxaD8DbEzT3OqzMlcQCGvtydifY6OkQzBVtjiDZrDuPeI6tnqDRb3hTG8fa3")
	// Faz a solicitação HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.Server{}, err
	}
	defer resp.Body.Close()

	// Lê a resposta da solicitação
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Server{}, err
	}
	var responseBody models.ResponseServerHetzner
	err = json.Unmarshal(body, &responseBody)
	fmt.Print("body", responseBody)

	if err != nil {
		return models.Server{}, err
	}

	dns, err := CreateDNSRecord(nameServer+".shub.tech", "A", responseBody.Server.PublicNet.IPv4.IP, 120, false)
	if err != nil {
		return models.Server{}, err
	}
	fmt.Println("dns", dns)

	newServer := models.Server{
		Name:      responseBody.Server.Name,
		IP:        responseBody.Server.PublicNet.IPv4.IP,
		CreatedAt: time.Now(),
		URL:       "https://" + nameServer + ".shub.tech",
	}

	if err := models.DB.Create(&newServer).Error; err != nil {
		return models.Server{}, err
	}

	serverIdGlobal <- newServer.ID
	cmd := exec.Command("sh", "../stacks/deploy_stack.sh", responseBody.Server.PublicNet.IPv4.IP, nameServer)

	error := cmd.Run()
	if error != nil {
		fmt.Print("Erro no script sh", error)
	}
	return newServer, nil
}

func CreateDNSRecord(name string, recordType string, content string, ttl int, proxied bool) (*http.Response, error) {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", os.Getenv("CLOUDFLARE_ZONE_ID"))

	record := models.DNSRecord{
		Type:    recordType,
		Name:    name,
		Content: content,
		TTL:     ttl,
		Proxied: proxied,
	}

	body, err := json.Marshal(record)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("CLOUDFLARE_API_TOKEN")))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error criar dns", err)
		return nil, err
	}

	return resp, nil
}
