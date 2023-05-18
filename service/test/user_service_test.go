package test

import (
	"context"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MCPutro/golang-docker/model"
	"github.com/MCPutro/golang-docker/model/web"
	"github.com/MCPutro/golang-docker/repository"
	"github.com/MCPutro/golang-docker/service"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestServiceUserCreate(t *testing.T) {

	request := &web.UserCreateRequest{
		Username: "unyil",
		Fullname: "unyil-unyilan",
		Password: "pass-unyil",
	}

	expectedId := 5

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO public."users" (.+) RETURNING user_id`).
		//WithArgs(request.Username, request.Fullname, request.Password).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedId))
	mock.ExpectCommit()

	ctx := context.Background()
	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository, db)
	resp, err := userService.Registration(ctx, request)

	//fmt.Println(">>", err)
	//fmt.Println(">>", resp)

	// we make sure that all expectations were met
	if errMock := mock.ExpectationsWereMet(); errMock != nil {
		t.Errorf("there were unfulfilled expectations: %s", errMock)
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedId, resp.Id)
}

//positive case
func TestServiceUpdateUser_Positive(t *testing.T) {

	request := &model.User{
		Id:       12,
		Username: "pa ogah",
		Fullname: "gope dulu",
		Password: "gope dulu",
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE public."users"`).
		WithArgs(request.Id, request.Fullname, sqlmock.AnyArg() /*request.Password*/, request.Username).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	ctx := context.Background()
	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository, db)
	resp, err := userService.Update(ctx, request)

	fmt.Println(">>", err)
	fmt.Println(">>", resp)

	// we make sure that all expectations were met
	if errMock := mock.ExpectationsWereMet(); errMock != nil {
		t.Errorf("there were unfulfilled expectations: %s", errMock)
	}

	assert.NoError(t, err)
	assert.Equal(t, request.Id, resp.Id)
	assert.Equal(t, request.Username, resp.Username)
}

//negative case
func TestServiceUpdateUser_Negative(t *testing.T) {

	request := &model.User{
		Id:       12,
		Username: "pa ogah",
		Fullname: "gope dulu",
		Password: "gope dulu",
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE public."users"`).
		WithArgs(request.Id, request.Fullname, sqlmock.AnyArg() /*request.Password*/, request.Username).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()

	ctx := context.Background()
	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository, db)
	resp, err := userService.Update(ctx, request)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.Error(t, err)
	assert.Nil(t, resp)

}

func TestGetAllUser(t *testing.T) {
	ctx := context.Background()
	users := []model.User{
		{Id: 1, Username: "user1", Fullname: "name1"},
		{Id: 2, Username: "user2", Fullname: "name2"},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock db: %v", err)
	}
	defer db.Close()

	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository, db)

	//positive case
	//set expect data
	rows := sqlmock.NewRows([]string{"user_id", "username", "fullname", "u.creation_date"})
	for _, user := range users {
		rows.AddRow(user.Id, user.Username, user.Fullname, user.CreationDate)
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`select u.user_id, u.username, u.fullname, u.creation_date from public."users" u`).
		WillReturnRows(rows)
	mock.ExpectCommit()

	getALl, err := userService.GetAll(ctx)

	// we make sure that all expectations were met
	if errMock := mock.ExpectationsWereMet(); errMock != nil {
		t.Errorf("there were unfulfilled expectations: %s", errMock)
	}

	assert.Len(t, getALl, len(users))

	//negative case
	mock.ExpectBegin()
	mock.ExpectQuery(`select u.user_id, u.username, u.fullname, u.creation_date from public."users" u`).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}))
	mock.ExpectRollback()

	list, err := userService.GetAll(ctx)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.Error(t, err)
	assert.Nil(t, list)
}

func TestGetUserById(t *testing.T) {
	ctx := context.Background()
	id := 1
	users := []model.User{
		{Id: 1, Username: "user1", Fullname: "name1"},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock db: %v", err)
	}
	defer db.Close()

	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository, db)

	//positive case
	//set expect data
	rows := sqlmock.NewRows([]string{"user_id", "username", "fullname", "u.creation_date"})
	for _, user := range users {
		rows.AddRow(user.Id, user.Username, user.Fullname, user.CreationDate)
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`select u.user_id, u.username, u.fullname, u.creation_date from public."users" u`).
		WithArgs(id).
		WillReturnRows(rows)
	mock.ExpectCommit()

	user, err1 := userService.GetById(ctx, id)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NoError(t, err1)
	assert.Equal(t, user.Id, 1)

	//negative case
	mock.ExpectBegin()
	mock.ExpectQuery(`select u.user_id, u.username, u.fullname, u.creation_date from public."users" u`).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}))
	mock.ExpectRollback()

	list, err2 := userService.GetAll(ctx)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.Error(t, err2)
	assert.Nil(t, list)
}

func TestLogin(t *testing.T) {
	ctx := context.Background()

	request := web.UserCreateRequest{
		Username: "user1",
		Password: "user1",
	}

	users := []model.User{
		{Id: 1, Username: "user1", Fullname: "name1", Password: "$2a$10$5CSbYma21UNed8iAkhTnh.RDobwn5dYPflW/oQb/1sVSJPOv7M9Pe"},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock db: %v", err)
	}
	defer db.Close()

	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository, db)

	//positive case
	//set expect data
	rows := sqlmock.NewRows([]string{"user_id", "username", "fullname", "password", "u.creation_date"})
	for _, user := range users {
		rows.AddRow(user.Id, user.Username, user.Fullname, user.Password, user.CreationDate)
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`select u.user_id, u.username, u.fullname, u.password, u.creation_date from public."users" u`).
		WithArgs(request.Username).
		WillReturnRows(rows)
	mock.ExpectCommit()

	login1, err1 := userService.Login(ctx, &request)

	// we make sure that all expectations were met
	if errMock := mock.ExpectationsWereMet(); errMock != nil {
		t.Errorf("there were unfulfilled expectations: %s", errMock)
	}

	assert.NoError(t, err1)
	assert.Equal(t, login1.Id, 1)

	//negative case username not found
	mock.ExpectBegin()
	mock.ExpectQuery(`select u.user_id, u.username, u.fullname, u.password, u.creation_date from public."users" u`).
		WithArgs(request.Username).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}))
	mock.ExpectRollback()

	login2, err2 := userService.Login(ctx, &request)

	// we make sure that all expectations were met
	if errMock2 := mock.ExpectationsWereMet(); errMock2 != nil {
		t.Errorf("there were unfulfilled expectations: %s", errMock2)
	}

	assert.Error(t, err2)
	assert.Nil(t, login2)
}

func TestDeleteUser(t *testing.T) {
	ctx := context.Background()
	id := 1

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock db: %v", err)
	}
	defer db.Close()

	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository, db)

	//positive case
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM public."users"`).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err1 := userService.Remove(ctx, id)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	assert.NoError(t, err1)

	//negative case
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM public."users"`).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()

	err2 := userService.Remove(ctx, id)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	assert.Error(t, err2)
}
