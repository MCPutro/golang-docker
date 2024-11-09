package test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MCPutro/golang-docker/internal/controller"
	"github.com/MCPutro/golang-docker/internal/model"
	"github.com/MCPutro/golang-docker/internal/repository/user"
	"github.com/MCPutro/golang-docker/internal/routes"
	user2 "github.com/MCPutro/golang-docker/internal/service/user"
	"github.com/MCPutro/golang-docker/internal/util"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var token string

func setup(db *sql.DB) *fiber.App {

	userRepository := user.NewUserRepository()
	userService := user2.NewUserService(userRepository, db)
	userController := controller.NewUserController(userService)

	router := routes.NewRouter(userController)

	token, _ = util.GenerateToken(&model.User{
		Id:       0,
		Username: "",
	})

	//fmt.Println(">>", token)

	return router
}

func Test_userController_Login(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("error init sql mock, error : %s", err)
	}
	defer db.Close()

	router := setup(db)

	tests := []struct {
		name               string
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
			name: "positive case",
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
			name: "incorrect password",
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
			name: "username not found",
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

func Test_userController_Registration(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("error init sql mock, error : %s", err)
	}
	defer db.Close()

	router := setup(db)

	tests := []struct {
		name               string
		mock               func()
		reqBody            io.Reader
		url                string
		method             string
		wantErr            bool
		expectedStatusCode int
		expectedError      error
	}{
		{
			name: "positive case",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(`select u.user_id, u.username, u.fullname, u.password, u.creation_date from public."users" u`).WithArgs("si_unyil").
					WillReturnRows(sqlmock.NewRows([]string{"user_id", "username", "fullname", "password", "u.creation_date"}))
				mock.ExpectQuery(`INSERT INTO public."users" (.+) RETURNING user_id`).
					WithArgs("si_unyil", "si unyil ke 1", sqlmock.AnyArg() /*request.Password*/).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
				mock.ExpectCommit()
			},
			reqBody:            strings.NewReader(`{"username" : "si_unyil", "fullname" : "si unyil ke 1", "password" : "123456789"}`),
			url:                "/user/registration",
			method:             http.MethodPost,
			wantErr:            false,
			expectedStatusCode: http.StatusCreated,
			expectedError:      nil,
		},
		{
			name: "username already used",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(`select (.+) from public."users" u`).WithArgs("si_unyil").
					WillReturnRows(sqlmock.NewRows([]string{"user_id", "username", "fullname", "password", "u.creation_date"}).
						AddRow(0, "si_unyil", "si unyil ke 1", "$2a$10$vzSUW9Zqo7O0UYrsSQE6LOs359dcuVPj6dlLPmOv4a4uwIQH5Ue0u", time.Now().String()))
				mock.ExpectRollback()
			},
			reqBody:            strings.NewReader(`{"username" : "si_unyil", "fullname" : "si unyil ke 1", "password" : "123456789"}`),
			url:                "/user/registration",
			method:             http.MethodPost,
			wantErr:            true,
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedError:      util.ErrAlreadyUsed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			httpRequest := httptest.NewRequest(tt.method, tt.url, tt.reqBody)
			httpRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

			test, _ := router.Test(httpRequest, 30000)

			if errMock := mock.ExpectationsWereMet(); errMock != nil {
				t.Errorf("there were unfulfilled expectations: %s", errMock)
			}

			body, _ := io.ReadAll(test.Body)
			var responseBody map[string]interface{}
			json.Unmarshal(body, &responseBody)

			assert.Equal(t, tt.expectedStatusCode, test.StatusCode)
			assert.Equal(t, tt.expectedStatusCode, int(responseBody["status"].(float64)))
			if tt.wantErr && tt.expectedError != nil {
				assert.Contains(t, responseBody["message"], tt.expectedError.Error())
			} else if !tt.wantErr {
				assert.Equal(t, "success", responseBody["message"])
			}

		})
	}
}

func Test_userController_ShowAllUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("error init sql mock, error : %s", err)
	}
	defer db.Close()

	router := setup(db)

	tests := []struct {
		name               string
		mock               func()
		reqBody            io.Reader
		url                string
		method             string
		wantErr            bool
		expectedStatusCode int
		expectedError      error
	}{
		{
			name: "positive case",
			mock: func() {
				rows := sqlmock.NewRows([]string{"user_id", "username", "fullname", "u.creation_date"}).
					AddRow(1, "user1", "name1", time.Now().String()).
					AddRow(2, "user2", "name2", time.Now().String()).
					AddRow(3, "user3", "name3", time.Now().String()).
					AddRow(4, "user4", "name4", time.Now().String())
				mock.ExpectBegin()
				mock.ExpectQuery(`select (.+) from public."users" u`).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			reqBody:            nil,
			url:                "/user",
			method:             http.MethodGet,
			wantErr:            false,
			expectedStatusCode: http.StatusOK,
			expectedError:      nil,
		},
		{
			name: "empty list",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(`select (.+) from public."users" u`).WillReturnRows(sqlmock.NewRows([]string{"user_id", "username", "fullname", "u.creation_date"}))
				mock.ExpectRollback()
			},
			reqBody:            nil,
			url:                "/user",
			method:             http.MethodGet,
			wantErr:            false,
			expectedStatusCode: http.StatusOK,
			expectedError:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			httpRequest := httptest.NewRequest(tt.method, tt.url, tt.reqBody)
			httpRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

			test, _ := router.Test(httpRequest, 30000)

			if errMock := mock.ExpectationsWereMet(); errMock != nil {
				t.Errorf("there were unfulfilled expectations: %s", errMock)
			}

			body, _ := io.ReadAll(test.Body)
			var responseBody map[string]interface{}
			json.Unmarshal(body, &responseBody)

			assert.Equal(t, tt.expectedStatusCode, test.StatusCode)
			assert.Equal(t, tt.expectedStatusCode, int(responseBody["status"].(float64)))
			if tt.wantErr && tt.expectedError != nil {
				assert.Contains(t, responseBody["message"], tt.expectedError.Error())
			} else {
				assert.Equal(t, "success", responseBody["message"])
			}

		})
	}
}

func Test_userController_ShowUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("error init sql mock, error : %s", err)
	}
	defer db.Close()

	router := setup(db)

	userId := 4

	tests := []struct {
		name               string
		mock               func()
		reqBody            io.Reader
		url                string
		method             string
		wantErr            bool
		expectedStatusCode int
		expectedError      error
	}{
		{
			name: "positive case",
			mock: func() {
				rows := sqlmock.NewRows([]string{"user_id", "username", "fullname", "u.creation_date"}).
					AddRow(userId, "user4", "name4", time.Now().String())
				mock.ExpectBegin()
				mock.ExpectQuery(`select (.+) from public."users" u`).WithArgs(userId).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			reqBody:            nil,
			url:                fmt.Sprintf("/user/%d", userId),
			method:             http.MethodGet,
			wantErr:            false,
			expectedStatusCode: http.StatusOK,
			expectedError:      nil,
		},
		{
			name: "user id not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"user_id", "username", "fullname", "u.creation_date"})
				mock.ExpectBegin()
				mock.ExpectQuery(`select (.+) from public."users" u`).WithArgs(userId).WillReturnRows(rows)
				mock.ExpectRollback()
			},
			reqBody:            nil,
			url:                fmt.Sprintf("/user/%d", userId),
			method:             http.MethodGet,
			wantErr:            true,
			expectedStatusCode: http.StatusNotFound,
			expectedError:      util.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			httpRequest := httptest.NewRequest(tt.method, tt.url, tt.reqBody)
			httpRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

			test, _ := router.Test(httpRequest, 30000)

			if errMock := mock.ExpectationsWereMet(); errMock != nil {
				t.Errorf("there were unfulfilled expectations: %s", errMock)
			}

			body, _ := io.ReadAll(test.Body)
			var responseBody map[string]interface{}
			json.Unmarshal(body, &responseBody)

			assert.Equal(t, tt.expectedStatusCode, test.StatusCode)
			assert.Equal(t, tt.expectedStatusCode, int(responseBody["status"].(float64)))
			if tt.wantErr && tt.expectedError != nil {
				assert.Contains(t, responseBody["message"], tt.expectedError.Error())
			} else {
				assert.Equal(t, "success", responseBody["message"])
			}

		})
	}
}

func Test_userController_UpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("error init sql mock, error : %s", err)
	}
	defer db.Close()

	router := setup(db)

	userId := 5

	tests := []struct {
		name               string
		mock               func()
		reqBody            io.Reader
		url                string
		method             string
		wantErr            bool
		expectedStatusCode int
		expectedError      error
	}{
		{
			name: "positive case",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE public."users"`).
					WithArgs(userId, "si unyil ke 0", sqlmock.AnyArg() /*,"password1-baru"*/, "si_unyil-0").
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			reqBody:            strings.NewReader(`{"username": "si_unyil-0", "fullname": "si unyil ke 0", "password": "123456789"}`),
			url:                fmt.Sprintf("/user/%d", userId),
			method:             http.MethodPut,
			wantErr:            false,
			expectedStatusCode: http.StatusOK,
			expectedError:      nil,
		},
		{
			name: "user id not found",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE public."users"`).
					WithArgs(userId, "si unyil ke 0", sqlmock.AnyArg() /*,"password1-baru"*/, "si_unyil-0").
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectRollback()
			},
			reqBody:            strings.NewReader(`{"username": "si_unyil-0", "fullname": "si unyil ke 0", "password": "123456789"}`),
			url:                fmt.Sprintf("/user/%d", userId),
			method:             http.MethodPut,
			wantErr:            true,
			expectedStatusCode: http.StatusNotFound,
			expectedError:      util.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			httpRequest := httptest.NewRequest(tt.method, tt.url, tt.reqBody)
			httpRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

			test, _ := router.Test(httpRequest, 30000)

			if errMock := mock.ExpectationsWereMet(); errMock != nil {
				t.Errorf("there were unfulfilled expectations: %s", errMock)
			}

			body, _ := io.ReadAll(test.Body)
			var responseBody map[string]interface{}
			json.Unmarshal(body, &responseBody)

			assert.Equal(t, tt.expectedStatusCode, test.StatusCode)
			assert.Equal(t, tt.expectedStatusCode, int(responseBody["status"].(float64)))
			if tt.wantErr && tt.expectedError != nil {
				assert.Contains(t, responseBody["message"], tt.expectedError.Error())
			} else {
				assert.Equal(t, "success", responseBody["message"])
				assert.Equal(t, "si_unyil-0", responseBody["data"].(map[string]interface{})["username"])
			}

		})

	}
}

func Test_userController_DeleteUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("error init sql mock, error : %s", err)
	}
	defer db.Close()

	userId := 20

	router := setup(db)

	tests := []struct {
		name               string
		mock               func()
		reqBody            io.Reader
		url                string
		method             string
		wantErr            bool
		expectedStatusCode int
		expectedError      error
	}{
		{
			name: "positive case",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM public."users"`).
					WithArgs(userId).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			reqBody:            nil,
			url:                fmt.Sprintf("/user/%d", userId),
			method:             http.MethodDelete,
			wantErr:            false,
			expectedStatusCode: http.StatusOK,
			expectedError:      nil,
		},
		{
			name: "user id not found",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM public."users"`).
					WithArgs(userId).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectRollback()
			},
			reqBody:            nil,
			url:                fmt.Sprintf("/user/%d", userId),
			method:             http.MethodDelete,
			wantErr:            true,
			expectedStatusCode: http.StatusNotFound,
			expectedError:      util.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			httpRequest := httptest.NewRequest(tt.method, tt.url, tt.reqBody)
			httpRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

			test, _ := router.Test(httpRequest, 30000)

			if errMock := mock.ExpectationsWereMet(); errMock != nil {
				t.Errorf("there were unfulfilled expectations: %s", errMock)
			}

			body, _ := io.ReadAll(test.Body)
			var responseBody map[string]interface{}
			json.Unmarshal(body, &responseBody)

			assert.Equal(t, tt.expectedStatusCode, test.StatusCode)
			assert.Equal(t, tt.expectedStatusCode, int(responseBody["status"].(float64)))
			if tt.wantErr && tt.expectedError != nil {
				assert.Contains(t, responseBody["message"], tt.expectedError.Error())
			} else {
				assert.Equal(t, "success", responseBody["message"])
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
