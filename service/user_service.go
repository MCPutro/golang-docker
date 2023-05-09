package service

import (
	"context"
	"github.com/MCPutro/golang-docker/model"
)

type UserService interface {
	Create(ctx context.Context, req *model.User) (*model.User, error)
	Update(ctx context.Context, req *model.User) (*model.User, error)
	GetALl(ctx context.Context) ([]*model.User, error)
	GetById(ctx context.Context, id int) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	Remove(ctx context.Context, id int) error
}
