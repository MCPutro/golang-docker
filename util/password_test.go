package util

import (
	"testing"
)

func TestName(t *testing.T) {
	password := "1Abcdefgh"
	hashedPassword, err := EncryptPassword(password)

	if err != nil {
		t.Errorf("Error hashing password: %s", err)
	}

	if len(hashedPassword) == 0 {
		t.Errorf("Hashed password is empty")
	}

	if !Compare(password, hashedPassword) {
		t.Errorf("Password and hash comparison failed")
	}

	if Compare("132456789", hashedPassword) {
		t.Errorf("Incorrect password was matched")
	}
}
