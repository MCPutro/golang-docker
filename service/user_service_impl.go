package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/MCPutro/golang-docker/model"
	"github.com/MCPutro/golang-docker/model/web"
	"github.com/MCPutro/golang-docker/repository"
	"github.com/MCPutro/golang-docker/util"
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
	defer func() { util.CommitOrRollback(err, tx) }()

	password, err := util.EncryptPassword(req.Password)
	if err != nil {
		return nil, err
	}

	message := &model.User{
		Username: req.Username,
		FullName: req.FullName,
		Password: password,
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
	defer func() { util.CommitOrRollback(err, tx) }()

	//call service
	err = u.repo.Update(ctx, tx, req)

	fmt.Println(">>", err)

	if err != nil {
		return nil, err
	} else {
		return req, nil
	}
}

func (u *userServiceImpl) GetAll(ctx context.Context) ([]*model.User, error) {
	//Begin db transactional
	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() { util.CommitOrRollback(err, tx) }()

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
	defer func() { util.CommitOrRollback(err, tx) }()

	//call service
	findByID, err := u.repo.FindByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return findByID, nil
}

func (u *userServiceImpl) Login(ctx context.Context, req *web.UserCreateRequest) (*model.User, error) {
	//Begin db transactional
	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() { util.CommitOrRollback(err, tx) }()

	//call service
	findByUsername, err := u.repo.FindByUsername(ctx, tx, req.Username)

	if err != nil {
		return nil, err
	}

	//validation password
	if util.ComparePassword(req.Password, findByUsername.Password) {
		return findByUsername, nil
	}

	return nil, errors.New("username and password not match")
}

func (u *userServiceImpl) Remove(ctx context.Context, id int) error {
	//Begin db transactional
	tx, err := u.db.Begin()
	if err != nil {
		return err
	}
	defer func() { util.CommitOrRollback(err, tx) }()

	//call service
	err = u.repo.Delete(ctx, tx, id)

	if err != nil {
		return err
	}

	return nil
}
