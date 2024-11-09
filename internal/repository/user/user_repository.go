package user

import (
	"context"
	"database/sql"
	"github.com/MCPutro/golang-docker/internal/model"
)

// Repository mendefinisikan func yg akan digunakan untuk oleh UserRepository
type Repository interface {
	// Save untuk menyimpan data user baru
	Save(ctx context.Context, tx *sql.Tx, newUser *model.User) (*model.User, error)

	// FindAll menampilan seluruh data user yang sudah terdaftar
	FindAll(ctx context.Context, tx *sql.Tx) ([]*model.User, error)

	// FindByID menampilkan user berdasarkan user ID
	FindByID(ctx context.Context, tx *sql.Tx, Id int) (*model.User, error)

	// FindByUsername menampilkan user berdasarkan username
	FindByUsername(ctx context.Context, tx *sql.Tx, Username string) (*model.User, error)

	// Update memperbaharui data user
	Update(ctx context.Context, tx *sql.Tx, newUser *model.User) error

	// Delete menghapus data user berdasarkan user Id
	Delete(ctx context.Context, tx *sql.Tx, Id int) error
}
