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

func (u *userServiceImpl) Registration(ctx context.Context, req *web.UserCreateRequest) (*web.UserResponse, error) {
	//Begin db transactional
	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() { util.CommitOrRollback(err, tx) }()

	//check username already exist or not
	existingUsername, err := u.repo.FindByUsername(ctx, tx, req.Username)
	// if username is exists return error
	if existingUsername != nil {
		err = fmt.Errorf("username %w", util.ErrAlreadyUsed)
		return nil, err
	}

	password, err := util.EncryptPassword(req.Password)
	if err != nil {
		return nil, err
	}

	message := &model.User{
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

		return &web.UserResponse{
			Id:           message.Id,
			Username:     message.Username,
			Fullname:     message.Fullname,
			Token:        token,
			CreationDate: message.CreationDate,
		}, nil
	}
}

func (u *userServiceImpl) Update(ctx context.Context, req *model.User) (*web.UserResponse, error) {
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
		return &web.UserResponse{
			Id:       req.Id,
			Username: req.Username,
			Fullname: req.Fullname,
		}, nil
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

func (u *userServiceImpl) Login(ctx context.Context, req *web.UserCreateRequest) (*web.UserResponse, error) {
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

		return &web.UserResponse{
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
