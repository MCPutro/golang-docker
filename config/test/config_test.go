package test

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoadEnvFromFile(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		t.Error("Error loading .env file")
	}

	Host := os.Getenv("POSTGRES_HOSTNAME")
	Pass := os.Getenv("POSTGRES_PASSWORD")
	Username := os.Getenv("POSTGRES_USER")
	Name := os.Getenv("POSTGRES_DB")
	Port := os.Getenv("POSTGRES_DB_PORT")

	assert.NotEmpty(t, Host)
	assert.NotEmpty(t, Pass)
	assert.NotEmpty(t, Username)
	assert.NotEmpty(t, Name)
	assert.NotEmpty(t, Port)
}
