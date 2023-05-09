package database

import (
	conf "github.com/MCPutro/golang-docker/config"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestCallDatabase(t *testing.T) {
	//load enviroment variable from .env file
	if err := godotenv.Load("../.env"); err != nil {
		t.Error("Error loading .env file")
	}

	conf.DB_Host = os.Getenv("POSTGRES_HOSTNAME")
	conf.DB_Pass = os.Getenv("POSTGRES_PASSWORD")
	conf.DB_Username = os.Getenv("POSTGRES_USER")
	conf.DB_Name = os.Getenv("POSTGRES_DB")
	conf.DB_Port = os.Getenv("POSTGRES_DB_PORT")

	db, err := InitDatabase()
	if err != nil {
		log.Fatalln(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("gagal ping db, error : %s", err)

	}

	log.Fatalln("ssss")

	assert.NoError(t, err)

}
