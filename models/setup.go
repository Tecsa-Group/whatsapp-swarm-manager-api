package models

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
	)
	fmt.Print("aqui", dsn)
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	err = database.AutoMigrate(&Server{}, &Instance{})
	if err != nil {
		return fmt.Errorf("failed to auto migrate tables: %w", err)
	}

	newServer := Server{
		Name:             "primeiro servidor",
		IP:               "http://5.161.71.166/",
		Port:             8080,
		Active:           true,
		CreatedAt:        time.Now(),
		URL:              "http://evolution.shub.tech",
		InstanceQuantity: 0,
	}
	if err := database.Create(&newServer).Error; err != nil {
		return fmt.Errorf("failed to create new server: %w", err)
	}

	DB = database
	return nil
}
