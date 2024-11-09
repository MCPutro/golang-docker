package repository

import (
	"github.com/MCPutro/golang-docker/internal/repository/user"
	"sync"
)

var (
	userRepository     user.Repository
	userRepositoryOnce sync.Once
)

type Repository interface {
	UserRepository() user.Repository
}

type repositoryManager struct{}

func NewRepositoryManager() Repository {
	return &repositoryManager{}
}

func (r *repositoryManager) UserRepository() user.Repository {
	userRepositoryOnce.Do(func() {
		userRepository = user.NewRepository()
	})

	return userRepository
}
