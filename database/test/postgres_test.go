package test

import (
	conf "github.com/MCPutro/golang-docker/config"
	"github.com/MCPutro/golang-docker/database"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestCallDatabase(t *testing.T) {
	//load enviroment variable from .env file
	if err := godotenv.Load("../../.env"); err != nil {
		t.Error("Error loading .env file")
	}

	conf.DB_Host = os.Getenv("POSTGRES_HOSTNAME")
	conf.DB_Pass = os.Getenv("POSTGRES_PASSWORD")
	conf.DB_Username = os.Getenv("POSTGRES_USER")
	conf.DB_Name = os.Getenv("POSTGRES_DB")
	conf.DB_Port = os.Getenv("POSTGRES_DB_PORT")

	db, err := database.InitDatabase()
	if err != nil {
		log.Fatalln(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("fail ping db, error : %s", err)

	}

	assert.NoError(t, err)

}
