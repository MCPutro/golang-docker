package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/MCPutro/golang-docker/model"
)

type userRepositoryImpl struct {
}

func NewUserRepositoryImpl() UserRepository {
	return &userRepositoryImpl{}
}

func (u *userRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, newUser *model.User) error {
	SQL := `INSERT INTO public."users" (username, fullname, password) VALUES ( $1, $2, $3);`

	_, err := tx.ExecContext(ctx, SQL, newUser.Username, newUser.FullName, newUser.Password)
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]*model.User, error) {
	SQL := `select u.user_id, u.username, u.fullname, u.creation_date from public."users" u;`

	rows, err := tx.QueryContext(ctx, SQL)
	if err != nil {
		return nil, err
	}

	var users []*model.User
	for rows.Next() {
		var temp model.User
		if err = rows.Scan(&temp.Id, &temp.Username, &temp.FullName, &temp.CreationDate); err != nil {
			return nil, err
		}
		users = append(users, &temp)
	}

	return users, nil
}

func (u *userRepositoryImpl) FindByID(ctx context.Context, tx *sql.Tx, Id int) (*model.User, error) {
	SQL := `select u.user_id, u.username, u.fullname, u.creation_date from public."users" u where u.user_id = $1;`

	row, err := tx.QueryContext(ctx, SQL, Id)
	if err != nil {
		return nil, err
	}

	if row.Next() {
		var temp model.User
		if err = row.Scan(&temp.Id, &temp.Username, &temp.FullName, &temp.CreationDate); err != nil {
			return nil, err
		}
		return &temp, nil
	}

	return nil, errors.New("no data found")
}

func (u *userRepositoryImpl) FindByUsername(ctx context.Context, tx *sql.Tx, Username string) (*model.User, error) {
	SQL := `select u.user_id, u.username, u.fullname, u.creation_date from public."users" u where u.username = $1;`

	row, err := tx.QueryContext(ctx, SQL, Username)
	if err != nil {
		return nil, err
	}

	if row.Next() {
		var temp model.User
		if err = row.Scan(&temp.Id, &temp.Username, &temp.FullName, &temp.CreationDate); err != nil {
			return nil, err
		}
		return &temp, nil
	}

	return nil, errors.New("no data found")
}

func (u *userRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, newUser *model.User) error {
	SQL := `UPDATE public."users" SET fullname = $2, password = $3 WHERE id = $1;`

	result, err := tx.ExecContext(ctx, SQL, newUser.Id, newUser.FullName, newUser.Password)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected > 0 {
		return nil
	} else {
		return errors.New("no data found")
	}

}

func (u *userRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, Id int) error {
	SQL := `DELETE FROM public."users" WHERE id = $1;`

	result, err := tx.ExecContext(ctx, SQL, Id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected > 0 {
		return nil
	} else {
		return errors.New("no data found")
	}
}
