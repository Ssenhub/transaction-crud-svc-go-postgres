package storage

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DockerHost   string
	InternalHost string
	LocalHost    string
	Port         string
	Password     string
	User         string
	DbName       string
	SslMode      string
}

func NewConnection(config *Config) (*gorm.DB, error) {
	docker_dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.DockerHost, config.Port, config.User, config.Password, config.DbName, config.SslMode)

	internal_dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.InternalHost, config.Port, config.User, config.Password, config.DbName, config.SslMode)

	local_dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.LocalHost, config.Port, config.User, config.Password, config.DbName, config.SslMode)

	db, err := gorm.Open(postgres.Open(local_dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Cannot connect to local DB")
	} else {
		fmt.Println("Connected to local DB")
		return db, nil
	}

	db, err = gorm.Open(postgres.Open(docker_dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Cannot connect to composed docker DB")
	} else {
		fmt.Println("Connected to composed docker DB")
		return db, nil
	}

	db, err = gorm.Open(postgres.Open(internal_dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Cannot connect to internal DB")
	} else {
		fmt.Println("Connected to internal DB")
		return db, nil
	}

	return nil, err
}
