package service

import (
	"context"
	"github.com/MCPutro/golang-docker/model"
	"github.com/MCPutro/golang-docker/model/web"
)

type UserService interface {
	Registration(ctx context.Context, req *web.UserCreateRequest) (*web.UserResponse, error)
	Update(ctx context.Context, req *model.User) (*model.User, error)
	GetAll(ctx context.Context) ([]*model.User, error)
	GetById(ctx context.Context, id int) (*model.User, error)
	Login(ctx context.Context, req *web.UserCreateRequest) (*web.UserResponse, error)
	Remove(ctx context.Context, id int) error
}
