package user

import (
	"context"
	"database/sql"
	"github.com/MCPutro/golang-docker/internal/entity"
)

// Repository mendefinisikan func yg akan digunakan untuk oleh UserRepository
type Repository interface {
	// Save untuk menyimpan data user baru
	Save(ctx context.Context, tx *sql.Tx, newUser *entity.User) (*entity.User, error)

	// FindAll menampilan seluruh data user yang sudah terdaftar
	FindAll(ctx context.Context, tx *sql.Tx) ([]*entity.User, error)

	// FindByID menampilkan user berdasarkan user ID
	FindByID(ctx context.Context, tx *sql.Tx, Id int) (*entity.User, error)

	// FindByUsername menampilkan user berdasarkan username
	FindByUsername(ctx context.Context, tx *sql.Tx, Username string) (*entity.User, error)

	// Update memperbaharui data user
	Update(ctx context.Context, tx *sql.Tx, newUser *entity.User) error

	// Delete menghapus data user berdasarkan user Id
	Delete(ctx context.Context, tx *sql.Tx, Id int) error
}
