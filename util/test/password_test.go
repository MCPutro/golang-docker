package test

import (
	"github.com/MCPutro/golang-docker/util"
	"testing"
)

func TestPasswordHash(t *testing.T) {
	password := "1Abcdefgh"
	hashedPassword, err := util.EncryptPassword(password)

	if err != nil {
		t.Errorf("Error hashing password: %s", err)
	}

	if len(hashedPassword) == 0 {
		t.Errorf("Hashed password is empty")
	}

	if !util.ComparePassword(password, hashedPassword) {
		t.Errorf("Password and hash comparison failed")
	}

	if util.ComparePassword("132456789", hashedPassword) {
		t.Errorf("Incorrect password was matched")
	}
}
