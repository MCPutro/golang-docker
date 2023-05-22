package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var (
	DB_Host     = mustGetEnv("POSTGRES_HOSTNAME")
	DB_Pass     = mustGetEnv("POSTGRES_PASSWORD")
	DB_Username = mustGetEnv("POSTGRES_USER")
	DB_Name     = mustGetEnv("POSTGRES_DB")
	DB_Port     = mustGetEnv("POSTGRES_DB_PORT")
	App_Port    = mustGetEnv("APP_PORT")
)

func mustGetEnv(k string) string {
	/* load from file .env */
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
		return ""
	}

	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("Warning: %s environment variable not set.", k)
	}

	log.Println("env", k, "is OK")
	return v
}
