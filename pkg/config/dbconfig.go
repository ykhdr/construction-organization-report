package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

func NewDBConfig() (*DBConfig, error) {
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, errors.Join(errors.New("DB_PORT isn't passed or not integer"), err)
	}

	dbConfig := &DBConfig{
		Host:     dbHost,
		Port:     dbPort,
		User:     dbUser,
		Password: dbPassword,
		Name:     dbName,
	}

	return dbConfig, nil
}

func (c *DBConfig) ConnectionInfo() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.Name,
	)
}
