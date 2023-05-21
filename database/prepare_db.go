package database

import (
	"context"
	"database/sql"
	"github.com/MCPutro/golang-docker/util"
)

func Create_table(db *sql.DB) error {
	Sql1 := `create table if not exists users
	(
		user_id       serial not null
			constraint users_pkey
				primary key,
		username      varchar(15)                                                        not null
			constraint users_username_key
				unique,
		fullname      varchar(200),
		password      varchar,
		creation_date timestamp with time zone default now()                             not null
	);`

	Sql2 := `INSERT INTO public.users (user_id, username, fullname, password) VALUES (0, 'admin.support', 'administrator', '$2a$10$/LsESEBsKbfWx1rIMPwuGekMKdHE2S169tiuJb3mMtiqj6crJn8Da') 
					ON CONFLICT (user_id) 
-- 					DO UPDATE set user_id = 0, username = excluded.username, fullname = excluded.fullname, password = excluded.password;
					DO NOTHING;
	`

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() { util.CommitOrRollback(err, tx) }()

	_, err = tx.ExecContext(context.Background(), Sql1)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(context.Background(), Sql2)
	if err != nil {
		return err
	}

	return nil

}
