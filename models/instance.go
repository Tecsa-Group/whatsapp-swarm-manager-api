package models

import "time"

type Instance struct {
	ID        int        `gorm:"primary_key" json:"id"`
	Name      string     `gorm:"column:name" json:"name"`
	Status    string     `gorm:"column:status" json:"status"`
	ServerID  int        `gorm:"column:server_id" json:"server_id"`
	Server    Server     `gorm:"foreignkey:ServerID" json:"server"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
}
