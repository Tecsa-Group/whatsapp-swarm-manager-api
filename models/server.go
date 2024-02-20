package models

import "time"

type Server struct {
	ID               int       `gorm:"primary_key" json:"id"`
	Name             string    `json:"name"`
	IP               string    `json:"ip"`
	Port             int       `json:"port"`
	Active           bool      `json:"active"`
	CreatedAt        time.Time `json:"created_at"`
	URL              string    `json:"url"`
	InstanceQuantity int       `json:"instance_quantity"`
}
