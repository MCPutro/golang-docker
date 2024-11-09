package config

import (
	"os"
	"sync"

	"github.com/MCPutro/golang-docker/util/logger"
	"github.com/joho/godotenv"
)

var (
	DB_Host     string
	DB_Pass     string
	DB_Username string
	DB_Name     string
	DB_Port     string
	App_Port    string

	configOnce sync.Once
)

func NewConfig() {
	configOnce.Do(func() {
		DB_Host = mustGetEnv("POSTGRES_HOSTNAME")
		DB_Pass = mustGetEnv("POSTGRES_PASSWORD")
		DB_Username = mustGetEnv("POSTGRES_USER")
		DB_Name = mustGetEnv("POSTGRES_DB")
		DB_Port = mustGetEnv("POSTGRES_DB_PORT")
		App_Port = mustGetEnv("APP_PORT")
	})

}

func mustGetEnv(k string) string {
	/* load from file .env */
	err := godotenv.Load(".env")
	if err != nil {
		// log.Println("Error loading .env file")
		logger.GetLogger().Fatalf("Error loading .env file")
		return ""
	}

	v := os.Getenv(k)
	if v == "" {
		// log.Fatalf("Warning: %s environment variable not set.", k)
		logger.GetLogger().Fatalf("Warning: %s environment variable not set.", k)
	}

	logger.GetLogger().Infoln("env", k, "is OK")
	return v
}
