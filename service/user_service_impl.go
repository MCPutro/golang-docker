package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/MCPutro/golang-docker/model"
	"github.com/MCPutro/golang-docker/model/web"
	"github.com/MCPutro/golang-docker/repository"
	"github.com/MCPutro/golang-docker/util"
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
			log.Println("Rollback")
		} else {
			tx.Commit()
			log.Println("Commit")
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
	//Begin db transactional
	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}
	defer util.CommitOrRollback(err, tx)

	//call service
	err = u.repo.Update(ctx, tx, req)

	if err != nil {
		return nil, err
	} else {
		return req, nil
	}
}

func (u *userServiceImpl) GetALl(ctx context.Context) ([]*model.User, error) {
	//Begin db transactional
	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}
	defer util.CommitOrRollback(err, tx)

	//call service
	users, err := u.repo.FindAll(ctx, tx)

	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, errors.New("no data found")
	}

	return users, nil
}

func (u *userServiceImpl) GetById(ctx context.Context, id int) (*model.User, error) {
	//Begin db transactional
	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}
	defer util.CommitOrRollback(err, tx)

	//call service
	findByID, err := u.repo.FindByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return findByID, nil
}

func (u *userServiceImpl) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	//Begin db transactional
	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}
	defer util.CommitOrRollback(err, tx)

	//call service
	findByUsername, err := u.repo.FindByUsername(ctx, tx, username)

	if err != nil {
		return nil, err
	}

	return findByUsername, nil
}

func (u *userServiceImpl) Remove(ctx context.Context, id int) error {
	//Begin db transactional
	tx, err := u.db.Begin()
	if err != nil {
		return err
	}
	defer util.CommitOrRollback(err, tx)

	//call service
	err = u.repo.Delete(ctx, tx, id)

	if err != nil {
		return err
	}

	return nil
}
