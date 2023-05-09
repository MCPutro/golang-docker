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
)

func mustGetEnv(k string) string {
	/* load from file .env */
	err := godotenv.Load(".env")
	if err != nil {
		return ""
	}

	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("Warning: %s environment variable not set.", k)
	}

	log.Println("env", k, "is OK")
	return v
}

//type config struct {
//	DB_Host     string
//	DB_Pass     string
//	DB_Username string
//	DB_Name     string
//	DB_Port     string
//}
//
//var cfg *config
//var once sync.Once
//
//func Get() *config {
//	/* load from file .env */
//	err := godotenv.Load(".env")
//	if err != nil {
//		log.Fatal("can't load from .env error:", err)
//		return nil
//	}
//
//	once.Do(func() {
//		cfg = &config{
//			DB_Host:     os.Getenv("POSTGRES_HOSTNAME"),
//			DB_Pass:     os.Getenv("POSTGRES_PASSWORD"),
//			DB_Username: os.Getenv("POSTGRES_USER"),
//			DB_Name:     os.Getenv("POSTGRES_DB"),
//			DB_Port:     os.Getenv("POSTGRES_DB_PORT"),
//		}
//		fmt.Println("--test singleton--")
//	})
//
//	return cfg
//}
