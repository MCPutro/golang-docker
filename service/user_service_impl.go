package service

import (
	"context"
	"database/sql"
	"github.com/MCPutro/golang-docker/model"
	"github.com/MCPutro/golang-docker/model/web"
	"github.com/MCPutro/golang-docker/repository"
	"log"
)

type userServiceImpl struct {
	repo repository.UserRepository
	db   *sql.DB
}

func NewUserService(repo repository.UserRepository, db *sql.DB) UserService {
	return &userServiceImpl{repo: repo, db: db}
}

func (u *userServiceImpl) Create(ctx context.Context, req *web.UserCreateRequest) (*model.User, error) {
	//Begin db transactional
	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			log.Println("Commit")
		} else {
			tx.Commit()
			log.Println("Rollback")
		}
	}()

	message := &model.User{
		Username: req.Username,
		FullName: req.FullName,
		Password: req.Password,
	}

	//call service
	message, err = u.repo.Save(ctx, tx, message)

	if err != nil {
		return nil, err
	} else {
		return message, nil
	}
}

func (u *userServiceImpl) Update(ctx context.Context, req *model.User) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userServiceImpl) GetALl(ctx context.Context) ([]*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userServiceImpl) GetById(ctx context.Context, id int) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userServiceImpl) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userServiceImpl) Remove(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}
