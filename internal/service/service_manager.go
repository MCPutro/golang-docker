package service

import (
	"database/sql"
	"github.com/MCPutro/golang-docker/internal/repository"
	"github.com/MCPutro/golang-docker/internal/service/user"
	"sync"
)

var (
	userService     user.Service
	userServiceOnce sync.Once
)

type Service interface {
	UserService() user.Service
}

type serviceManager struct {
	repoManager repository.Repository
	db          *sql.DB
}

func (s *serviceManager) UserService() user.Service {
	userServiceOnce.Do(func() {
		userService = user.NewService(s.repoManager.UserRepository(), s.db)
	})

	return userService
}

func NewServiceManager(repoManager repository.Repository, db *sql.DB) Service {
	return &serviceManager{repoManager: repoManager, db: db}
}
