package test

import (
	"database/sql"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MCPutro/golang-docker/config"
	"github.com/MCPutro/golang-docker/controller"
	"github.com/MCPutro/golang-docker/repository"
	"github.com/MCPutro/golang-docker/service"
	"github.com/MCPutro/golang-docker/util"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func setup(db *sql.DB) *fiber.App {

	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository, db)
	userController := controller.NewUserController(userService)

	router := config.NewRouter(userController)

	return router
}

func Test_userControllerImpl_Login(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("error init sql mock, error : %s", err)
	}
	defer db.Close()

	router := setup(db)

	tests := []struct {
		name               string
		router             *fiber.App
		mock               func()
		reqBody            io.Reader
		url                string
		method             string
		wantErr            bool
		expectedStatusCode int
		expectedError      error
	}{
		// TODO: test case
		{
			name: "login - positive case",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(`select (.+) from public."users" u`).WithArgs("user1").
					WillReturnRows(sqlmock.NewRows([]string{"user_id", "username", "fullname", "password", "u.creation_date"}).
						AddRow(0, "user1", "name1", "$2a$10$vzSUW9Zqo7O0UYrsSQE6LOs359dcuVPj6dlLPmOv4a4uwIQH5Ue0u", time.Now().String()))
				mock.ExpectCommit()
			},
			reqBody:            strings.NewReader(`{"username" : "user1", "password" : "admin123"}`),
			url:                "/login",
			method:             http.MethodPost,
			wantErr:            false,
			expectedStatusCode: http.StatusOK,
			expectedError:      nil,
		},
		{
			name: "login - incorrect password",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(`select (.+) from public."users" u`).WithArgs("user1").
					WillReturnRows(sqlmock.NewRows([]string{"user_id", "username", "fullname", "password", "u.creation_date"}).
						AddRow(0, "user1", "name1", "$2a$10$vzSUW9Zqo7O0UYrsSQE6LOs359dcuVPj6dlLPmOv4a4uwIQH5Ue0u", time.Now().String()))
				mock.ExpectRollback()
			},
			reqBody:            strings.NewReader(`{"username" : "user1", "password" : "1123456"}`),
			url:                "/login",
			method:             http.MethodPost,
			wantErr:            true,
			expectedStatusCode: http.StatusUnauthorized,
			expectedError:      util.ErrNotMatch,
		},
		{
			name: "login - username not found",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(`select (.+) from public."users" u`).WithArgs("user1").
					WillReturnRows(sqlmock.NewRows([]string{"user_id", "username", "fullname", "password", "u.creation_date"}))
				mock.ExpectRollback()
			},
			reqBody:            strings.NewReader(`{"username" : "user1", "password" : "1123456"}`),
			url:                "/login",
			method:             http.MethodPost,
			wantErr:            true,
			expectedStatusCode: http.StatusNotFound,
			expectedError:      util.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			httpRequest := httptest.NewRequest(tt.method, tt.url, tt.reqBody)
			test, _ := router.Test(httpRequest, 30000)

			if errMock := mock.ExpectationsWereMet(); errMock != nil {
				t.Errorf("there were unfulfilled expectations: %s", errMock)
			}

			body, _ := io.ReadAll(test.Body)
			var responseBody map[string]interface{}
			json.Unmarshal(body, &responseBody)

			assert.Equal(t, tt.expectedStatusCode, test.StatusCode)
			assert.Equal(t, tt.expectedStatusCode, int(responseBody["status"].(float64)))
			if tt.wantErr {
				assert.Contains(t, responseBody["message"], tt.expectedError.Error())
			} else {
				assert.Equal(t, "success", responseBody["message"])
				assert.NotEmpty(t, responseBody["data"].(map[string]interface{})["token"])
			}

		})
	}
}

//func setupData() (*fiber.App, sqlmock.Sqlmock) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		return nil, nil
//	}
//	//defer db.Close()
//
//	userRepository := repository.NewUserRepository()
//	userService := service.NewUserService(userRepository, db)
//	userController := controller.NewUserController(userService)
//
//	router := config.NewRouter(userController)
//
//	return router, mock
//}
//
//func Test_userController_Login(t *testing.T) {
//
//	router, mock := setupData()
//	//mock data resp
//	users := []model.User{
//		{Id: 1, Username: "user1", Fullname: "name1", Password: "$2a$10$vzSUW9Zqo7O0UYrsSQE6LOs359dcuVPj6dlLPmOv4a4uwIQH5Ue0u"},
//	}
//	//set expect data
//	rows := sqlmock.NewRows([]string{"user_id", "username", "fullname", "password", "u.creation_date"})
//	for _, user := range users {
//		rows.AddRow(user.Id, user.Username, user.Fullname, user.Password, user.CreationDate)
//	}
//	mock.ExpectBegin()
//	mock.ExpectQuery(`select (.+) from public."users" u`).WithArgs(users[0].Username).WillReturnRows(rows)
//	mock.ExpectCommit()
//
//	requestBody := strings.NewReader(`{"username" : "user1", "password" : "admin123"}`)
//	req := httptest.NewRequest(fiber.MethodPost, "/login", requestBody)
//	test, err := router.Test(req, 30000)
//
//	fmt.Println(err)
//	fmt.Println(test)
//
//	body, _ := io.ReadAll(test.Body)
//	var responseBody map[string]interface{}
//	json.Unmarshal(body, &responseBody)
//
//	assert.NoError(t, err)
//	assert.Equal(t, 200, test.StatusCode)
//	assert.Equal(t, 200, int(responseBody["status"].(float64)))
//	assert.Equal(t, "success", responseBody["message"])
//	assert.NotEmpty(t, responseBody["data"].(map[string]interface{})["token"])
//}
