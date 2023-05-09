package model

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserModel(t *testing.T) {

	//data dummy
	newUser := User{
		Id:           1,
		Username:     "emchepe",
		FullName:     "Mu'ti C Putro",
		CreationDate: "2023-12-12 10:12:13",
	}

	expectResult := `{"id":1,"username":"emchepe","full_name":"Mu'ti C Putro","creation_date":"2023-12-12 10:12:13"}`

	//meng-convert struct user menjadi json dengan tipe data byte
	result, err := json.Marshal(newUser)

	//membuat temp variable
	var temp User

	//meng-convert string json dalam bentuk byte menjadi struct dan menyimpannya ke dalam variable temp
	err1 := json.Unmarshal(result, &temp)

	assert.Nil(t, err)
	assert.Nil(t, err1)
	assert.Equal(t, expectResult, string(result))
	assert.Equal(t, temp.Username, "emchepe")
}
