package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/MCPutro/golang-docker/internal/entity"
	"github.com/MCPutro/golang-docker/internal/repository/user"
	"github.com/MCPutro/golang-docker/internal/util"
	"github.com/MCPutro/golang-docker/internal/web/request"
	"github.com/MCPutro/golang-docker/internal/web/response"
)

type serviceImpl struct {
	repo user.Repository
	db   *sql.DB
}

func NewService(repo user.Repository, db *sql.DB) Service {
	return &serviceImpl{repo: repo, db: db}
}

func (u *serviceImpl) Registration(ctx context.Context, req *request.UserCreate) (*response.UserResponse, error) {
	//Begin db transactional
	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() { util.CommitOrRollback(err, tx) }()

	//check username already exist or not
	existingUsername, err := u.repo.FindByUsername(ctx, tx, req.Username)
	// if username is existing return error
	if existingUsername != nil {
		err = fmt.Errorf("username %w", util.ErrAlreadyUsed)
		return nil, err
	}

	password, err := util.EncryptPassword(req.Password)
	if err != nil {
		return nil, err
	}

	message := &entity.User{
		Username: req.Username,
		Fullname: req.Fullname,
		Password: password,
	}

	//call service
	message, err = u.repo.Save(ctx, tx, message)

	if err != nil {
		return nil, err
	} else {
		//generate token
		token, err := util.GenerateToken(message)
		if err != nil {
			return nil, err
		}

		return &response.UserResponse{
			Id:           message.Id,
			Username:     message.Username,
			Fullname:     message.Fullname,
			Token:        token,
			CreationDate: message.CreationDate,
		}, nil
	}
}

func (u *serviceImpl) Update(ctx context.Context, req *entity.User) (*response.UserResponse, error) {
	//Begin db transactional
	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() { util.CommitOrRollback(err, tx) }()

	//hash password
	password, err := util.EncryptPassword(req.Password)
	if err != nil {
		return nil, err
	}
	req.Password = password

	//call service
	err = u.repo.Update(ctx, tx, req)

	if err != nil {
		return nil, err
	} else {
		return &response.UserResponse{
			Id:       req.Id,
			Username: req.Username,
			Fullname: req.Fullname,
		}, nil
	}
}

func (u *serviceImpl) GetAll(ctx context.Context) ([]*entity.User, error) {
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

func (u *serviceImpl) GetById(ctx context.Context, id int) (*entity.User, error) {
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

func (u *serviceImpl) Login(ctx context.Context, req *request.UserCreate) (*response.UserResponse, error) {
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
		//generate token
		token, err := util.GenerateToken(findByUsername)
		if err != nil {
			return nil, err
		}

		return &response.UserResponse{
			Id:           findByUsername.Id,
			Username:     findByUsername.Username,
			Fullname:     findByUsername.Fullname,
			Token:        token,
			CreationDate: findByUsername.CreationDate,
		}, nil
	}

	err = fmt.Errorf("username and password %w", util.ErrNotMatch)
	return nil, err
}

func (u *serviceImpl) Remove(ctx context.Context, id int) error {
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
