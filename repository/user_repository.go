package repository

import (
	"context"
	"database/sql"
	"github.com/MCPutro/golang-docker/model"
)

//mendefinisikan func yg akan digunakan untuk oleh UserRepository
type UserRepository interface {
	//untuk menyimpan data user baru
	Save(ctx context.Context, tx *sql.Tx, newUser *model.User) (*model.User, error)

	//menampilan seluruh data user yang sudah terdaftar
	FindAll(ctx context.Context, tx *sql.Tx) ([]*model.User, error)

	//menampilkan user berdasarkan user ID
	FindByID(ctx context.Context, tx *sql.Tx, Id int) (*model.User, error)

	//menampilkan user berdasarkan username
	FindByUsername(ctx context.Context, tx *sql.Tx, Username string) (*model.User, error)

	//memperbaharui data user
	Update(ctx context.Context, tx *sql.Tx, newUser *model.User) error

	//menghapus data user berdasarkan user Id
	Delete(ctx context.Context, tx *sql.Tx, Id int) error
}
