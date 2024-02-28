package models

import "time"

type Instance struct {
	ID        int        `gorm:"primary_key" json:"id"`
	Name      string     `gorm:"column:name" json:"name"`
	Status    string     `gorm:"column:status" json:"status"`
	ServerID  int        `gorm:"column:server_id" json:"server_id"`
	Server    Server     `gorm:"foreignkey:ServerID" json:"server"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	Apikey    string     `gorm:"uniqueIndex"`
}

type InstanceRequest struct {
	InstanceName string `json:"instanceName"`
	Token        string `json:"token"`
	QRCode       bool   `json:"qrcode"`
}

type ServerInstance struct {
	InstanceName      string `json:"instanceName"`
	Owner             string `json:"owner"`
	ProfileName       string `json:"profileName"`
	ProfilePictureURL string `json:"profilePictureUrl"`
	ProfileStatus     string `json:"profileStatus"`
	Status            string `json:"status"`
	ServerURL         string `json:"serverUrl"`
	APIKey            string `json:"apikey"`
}

type Result struct {
	URL       string
	CountOpen int
	ID        int
}

var UrlServer struct {
	URL string `gorm:"column:url"`
}
