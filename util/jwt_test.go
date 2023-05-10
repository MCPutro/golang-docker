package util

import (
	"github.com/MCPutro/golang-docker/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateToken(t *testing.T) {

	user := model.User{
		Id:       123569065,
		Username: "2",
		FullName: "3",
		Password: "4",
	}

	token, err := GenerateToken(&user)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	validateToken, err := ValidateToken(token)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	//if claims, ok := validateToken.Claims.(jwt.MapClaims); ok && validateToken.Valid {
	//	fmt.Println(claims["Username"], claims["id"])
	//}

	claims, ok := validateToken.Claims.(jwt.MapClaims)

	assert.True(t, ok)
	assert.Equal(t, claims["Username"], user.Username)
	assert.Equal(t, claims["Id"].(float64), float64(user.Id))

}

func TestValidateToken(t *testing.T) {
	//case expired token
	token := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJJZCI6MSwiVXNlcm5hbWUiOiIyIiwiZXhwIjoxNjgzNzE2MDYwLCJqdGkiOiIxYzlmMmJhOC05OTFmLTRkOWUtOTRmYi0xNjBiNjdhODk4MTIiLCJzdWIiOiJsb2dpbiJ9.5bI1T_GFL0xB_FaiTPljUF1SGFAxjRES16wmvjJ_gQc2JwyGjGerO-ICahJEYBS8QR6NxxlVkM_V3-MLe28e6w"

	validateToken, err := ValidateToken(token)

	assert.Error(t, err)
	assert.Nil(t, validateToken)
	assert.Equal(t, err, jwt.ErrTokenExpired)

	//case signature invalid
	token = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJJZCI6MSwiVXNlcm5hbWUiOiIyIiwiZXhwIjoxNjgzNzE2MDYwLCJqdGkiOiIxYzlmMmJhOC05OTFmLTRkOWUtOTRmYi0xNjBiNjdhODk4MTIiLCJzdWIiOiJsb2dpbngifQ.5bI1T_GFL0xB_FaiTPljUF1SGFAxjRES16wmvjJ_gQc2JwyGjGerO-ICahJEYBS8QR6NxxlVkM_V3-MLe28e6w"

	validateToken2, err := ValidateToken(token)

	assert.Error(t, err)
	assert.Nil(t, validateToken2)
	assert.Equal(t, err, jwt.ErrTokenSignatureInvalid)

}
