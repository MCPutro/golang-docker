package test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MCPutro/golang-docker/controller"
	"github.com/MCPutro/golang-docker/repository"
	"github.com/MCPutro/golang-docker/service"
	"github.com/gofiber/fiber/v2"
	"testing"
)

func Test_userControllerImpl_Login(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository, db)

	type fields struct {
		service service.UserService
	}

	type args struct {
		c *fiber.Ctx
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "OK",
			fields:  fields{service: userService},
			args:    args{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := controller.NewUserController(tt.fields.service)
			if err := u.Login(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
