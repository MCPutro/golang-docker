package user

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/MCPutro/golang-docker/internal/model"
	"github.com/MCPutro/golang-docker/internal/util"
)

type repositoryImpl struct {
}

func NewRepository() Repository {
	return &repositoryImpl{}
}

func (u *repositoryImpl) Save(ctx context.Context, tx *sql.Tx, newUser *model.User) (*model.User, error) {
	SQL := `INSERT INTO public."users" (username, fullname, password) VALUES ( $1, $2, $3) RETURNING user_id;`

	//_, err := tx.ExecContext(ctx, SQL, newUser.Username, newUser.Fullname, newUser.Password)
	var id int
	err := tx.QueryRowContext(ctx, SQL, newUser.Username, newUser.Fullname, newUser.Password).Scan(&id)
	if err != nil {
		return nil, err
	}
	newUser.Id = id

	return newUser, nil
}

func (u *repositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]*model.User, error) {
	SQL := `select u.user_id, u.username, u.fullname, u.creation_date from public."users" u;`

	rows, err := tx.QueryContext(ctx, SQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var temp model.User
		if err = rows.Scan(&temp.Id, &temp.Username, &temp.Fullname, &temp.CreationDate); err != nil {
			return nil, err
		}
		users = append(users, &temp)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("%w", util.ErrNotFound)
	}
	return users, nil
}

func (u *repositoryImpl) FindByID(ctx context.Context, tx *sql.Tx, Id int) (*model.User, error) {
	SQL := `select u.user_id, u.username, u.fullname, u.creation_date from public."users" u where u.user_id = $1;`

	row, err := tx.QueryContext(ctx, SQL, Id)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	if row.Next() {
		var temp model.User
		if err = row.Scan(&temp.Id, &temp.Username, &temp.Fullname, &temp.CreationDate); err != nil {
			return nil, err
		}
		return &temp, nil
	}

	return nil, fmt.Errorf("user id %w", util.ErrNotFound)
}

func (u *repositoryImpl) FindByUsername(ctx context.Context, tx *sql.Tx, Username string) (*model.User, error) {
	SQL := `select u.user_id, u.username, u.fullname, u.password, u.creation_date from public."users" u where u.username = $1;`

	row, err := tx.QueryContext(ctx, SQL, Username)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	if row.Next() {
		var temp model.User
		if err = row.Scan(&temp.Id, &temp.Username, &temp.Fullname, &temp.Password, &temp.CreationDate); err != nil {
			return nil, err
		}
		return &temp, nil
	}

	return nil, fmt.Errorf("username %w", util.ErrNotFound)
}

func (u *repositoryImpl) Update(ctx context.Context, tx *sql.Tx, newUser *model.User) error {
	SQL := `UPDATE public."users" SET fullname = $2, password = $3, username = $4 WHERE user_id = $1;`

	result, err := tx.ExecContext(ctx, SQL, newUser.Id, newUser.Fullname, newUser.Password, newUser.Username)
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
		return fmt.Errorf("user id %w", util.ErrNotFound)
	}

}

func (u *repositoryImpl) Delete(ctx context.Context, tx *sql.Tx, Id int) error {
	SQL := `DELETE FROM public."users" WHERE user_id = $1;`

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
		return fmt.Errorf("user id %w", util.ErrNotFound)
	}
}
