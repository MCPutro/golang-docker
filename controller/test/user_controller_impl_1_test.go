package test

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MCPutro/golang-docker/config"
	"github.com/MCPutro/golang-docker/controller"
	"github.com/MCPutro/golang-docker/model"
	"github.com/MCPutro/golang-docker/repository"
	"github.com/MCPutro/golang-docker/service"
	"github.com/gofiber/fiber/v2"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupData() (*fiber.App, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil
	}
	//defer db.Close()

	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository, db)
	userController := controller.NewUserController(userService)

	router := config.NewRouter(userController)

	return router, mock
}

func Test_userControllerImpl_Login1(t *testing.T) {

	router, mock := setupData()
	//mock data resp
	users := []model.User{
		{Id: 1, Username: "user1", FullName: "name1", Password: "$2a$10$vzSUW9Zqo7O0UYrsSQE6LOs359dcuVPj6dlLPmOv4a4uwIQH5Ue0u"},
	}
	//set expect data
	rows := sqlmock.NewRows([]string{"user_id", "username", "fullname", "password", "u.creation_date"})
	for _, user := range users {
		rows.AddRow(user.Id, user.Username, user.FullName, user.Password, user.CreationDate)
	}
	mock.ExpectBegin()
	mock.ExpectQuery(`select (.+) from public."users" u`).WithArgs(users[0].Username).WillReturnRows(rows)
	mock.ExpectCommit()

	requestBody := strings.NewReader(`{"username" : "user1", "password" : "admin123"}`)
	req := httptest.NewRequest(fiber.MethodPost, "/login", requestBody)
	test, err := router.Test(req, 30000)

	fmt.Println(err)
	fmt.Println(test.Body)
	body, _ := io.ReadAll(test.Body)
	fmt.Println(string(body))

}
