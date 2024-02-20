package models

import "time"

type Instance struct {
	ID        int        `gorm:"primary_key" json:"id"`
	Name      string     `db:"name" json:"name"`
	Status    string     `db:"status" json:"status"`
	Server_id string     `db:"server_id" json:"server_id"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
}
