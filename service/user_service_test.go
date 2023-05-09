package service

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MCPutro/golang-docker/model/web"
	"github.com/MCPutro/golang-docker/repository"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func Test_service_user_create(t *testing.T) {

	request := &web.UserCreateRequest{
		Username: "unyil",
		FullName: "unyil-unyilan",
		Password: "pass-unyil",
	}

	expectedId := 5

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO public."users" (.+) RETURNING id`).
		WithArgs(request.Username, request.FullName, request.Password).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedId))
	mock.ExpectCommit()

	ctx := context.Background()
	userRepository := repository.NewUserRepository()
	userService := NewUserService(userRepository, db)
	resp, err := userService.Create(ctx, request)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedId, resp.Id)
}
