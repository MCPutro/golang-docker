package user

import (
	"context"
	"github.com/MCPutro/golang-docker/internal/entity"
	"github.com/MCPutro/golang-docker/internal/web/request"
	"github.com/MCPutro/golang-docker/internal/web/response"
)

type Service interface {
	Registration(ctx context.Context, req *request.UserCreate) (*response.UserResponse, error)
	Update(ctx context.Context, req *entity.User) (*response.UserResponse, error)
	GetAll(ctx context.Context) ([]*entity.User, error)
	GetById(ctx context.Context, id int) (*entity.User, error)
	Login(ctx context.Context, req *request.UserCreate) (*response.UserResponse, error)
	Remove(ctx context.Context, id int) error
}
