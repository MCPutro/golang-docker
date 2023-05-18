package test

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MCPutro/golang-docker/model"
	"github.com/MCPutro/golang-docker/model/web"
	"github.com/MCPutro/golang-docker/repository"
	"github.com/MCPutro/golang-docker/service"
	"github.com/MCPutro/golang-docker/util"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

var creationDate = time.Now().String()

func Test_userServiceImpl_GetAll(t *testing.T) {
	db, mock, err1 := sqlmock.New()
	if err1 != nil {
		t.Errorf("Error sql mock, %s", err1)
	}
	defer db.Close()

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		mock    func()
		repo    repository.UserRepository
		args    args
		want    []*model.User
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "get all user - positive case",
			mock: func() {
				rows := sqlmock.NewRows([]string{"user_id", "username", "fullname", "u.creation_date"}).
					AddRow(1, "user1", "name1", creationDate).
					AddRow(2, "user2", "name2", creationDate).
					AddRow(3, "user3", "name3", creationDate).
					AddRow(4, "user4", "name4", creationDate)
				mock.ExpectBegin()
				mock.ExpectQuery(`select (.+) from public."users" u`).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			repo: repository.NewUserRepository(),
			args: args{ctx: context.Background()},
			want: []*model.User{
				{Id: 1, Username: "user1", Fullname: "name1", CreationDate: creationDate},
				{Id: 2, Username: "user2", Fullname: "name2", CreationDate: creationDate},
				{Id: 3, Username: "user3", Fullname: "name3", CreationDate: creationDate},
				{Id: 4, Username: "user4", Fullname: "name4", CreationDate: creationDate},
			},
			wantErr: false,
		},
		{
			name: "get all user - list is empty",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(`select (.+) from public."users" u`).WillReturnRows(sqlmock.NewRows([]string{"user_id"}))
				mock.ExpectRollback()
			},
			repo:    repository.NewUserRepository(),
			args:    args{ctx: context.Background()},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			u := service.NewUserService(tt.repo, db)
			got, err := u.GetAll(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				//t.Errorf("GetAll(%v)", tt.args.ctx)
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAll() = %v, want %v", got, tt.want)
			}
			if errMock := mock.ExpectationsWereMet(); errMock != nil {
				t.Errorf("there were unfulfilled expectations: %s", errMock)
			}
			assert.Equalf(t, tt.want, got, "GetAll(%v)", tt.args.ctx)
		})
	}
}

func Test_userServiceImpl_GetById(t *testing.T) {
	db, mock, err1 := sqlmock.New()
	if err1 != nil {
		t.Errorf("Error sql mock, %s", err1)
	}
	defer db.Close()

	userRepository := repository.NewUserRepository()

	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		name    string
		mock    func()
		repo    repository.UserRepository
		args    args
		want    *model.User
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "get user by id - positive case",
			mock: func() {
				rows := sqlmock.NewRows([]string{"user_id", "username", "fullname", "u.creation_date"}).
					AddRow(1, "user1", "name1", creationDate)
				mock.ExpectBegin()
				mock.ExpectQuery(`select u.user_id, u.username, u.fullname, u.creation_date from public."users" u`).WithArgs(1).
					WillReturnRows(rows)
				mock.ExpectCommit()
			},
			repo: userRepository,
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want: &model.User{
				Id:           1,
				Username:     "user1",
				Fullname:     "name1",
				Password:     "",
				CreationDate: creationDate,
			},
			wantErr: false,
		},
		{
			name: "get user by id - data not found",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(`select u.user_id, u.username, u.fullname, u.creation_date from public."users" u`).WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"user_id", "username", "fullname", "creation_date"}))
				mock.ExpectRollback()
			},
			repo: userRepository,
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			u := service.NewUserService(tt.repo, db)
			got, err := u.GetById(tt.args.ctx, tt.args.id)
			//if !tt.wantErr(t, err1, fmt.Sprintf("GetById(%v, %v)", tt.args.ctx, tt.args.id)) {
			//	return
			//}

			if (err != nil) != tt.wantErr {
				//t.Errorf("GetById(%v, %v)", tt.args.ctx, tt.args.id)
				t.Errorf("GetById() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetById() = %v, want %v", got, tt.want)
			}
			if errMock := mock.ExpectationsWereMet(); errMock != nil {
				t.Errorf("there were unfulfilled expectations: %s", errMock)
			}
			//assert.Equalf(t, tt.want, got, "GetById(%v, %v)", tt.args.ctx, tt.args.id)
			assert.Equalf(t, tt.want, got, "GetById(%v, %v)", tt.args.ctx, tt.args.id)
		})
	}
}

func Test_userServiceImpl_Login(t *testing.T) {
	db, mock, err1 := sqlmock.New()
	if err1 != nil {
		t.Errorf("Error sql mock, %s", err1)
	}
	defer db.Close()

	request := web.UserCreateRequest{
		Username: "admin.support",
		Password: "admin123",
	}

	type args struct {
		ctx context.Context
		req *web.UserCreateRequest
	}
	tests := []struct {
		name               string
		mock               func()
		repo               repository.UserRepository
		args               args
		want               *web.UserResponse
		wantErr            bool
		expectErrorMessage error
	}{
		// TODO: Add test cases.
		{
			name: "login - positive case",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(`select u.user_id, u.username, u.fullname, u.password, u.creation_date from public."users" u`).
					WithArgs(request.Username).
					WillReturnRows(sqlmock.NewRows([]string{"user_id", "username", "fullname", "password", "u.creation_date"}).
						AddRow(0, "admin.support", "Administrator", "$2a$10$vzSUW9Zqo7O0UYrsSQE6LOs359dcuVPj6dlLPmOv4a4uwIQH5Ue0u", creationDate))
				mock.ExpectCommit()
			},
			repo: repository.NewUserRepository(),
			args: args{ctx: context.Background(), req: &request},
			want: &web.UserResponse{
				Id:           0,
				Username:     request.Username,
				Fullname:     "Administrator",
				CreationDate: creationDate,
			},
			wantErr:            false,
			expectErrorMessage: nil,
		},
		{
			name: "login - password incorrect",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(`select u.user_id, u.username, u.fullname, u.password, u.creation_date from public."users" u`).
					WithArgs(request.Username).
					WillReturnRows(sqlmock.NewRows([]string{"user_id", "username", "fullname", "password", "u.creation_date"}).
						AddRow(0, "admin.support", "Administrator", "$2a$10$vzSUW9Zqo7O0UYrsSQE6LOs359dcuVPj6dlLPmOv4a4uwIQH5Ue0u", creationDate))
				mock.ExpectRollback()
			},
			repo: repository.NewUserRepository(),
			args: args{
				ctx: context.Background(),
				req: &web.UserCreateRequest{
					Username: request.Username,
					Password: "request.Password",
				},
			},
			want:               nil,
			wantErr:            true,
			expectErrorMessage: util.ErrNotMatch,
		},
		{
			name: "login - username not found",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(`select u.user_id, u.username, u.fullname, u.password, u.creation_date from public."users" u`).
					WithArgs(request.Username).
					WillReturnRows(sqlmock.NewRows([]string{"user_id", "username", "fullname", "password", "u.creation_date"}))
				mock.ExpectRollback()
			},
			repo: repository.NewUserRepository(),
			args: args{
				ctx: context.Background(),
				req: &web.UserCreateRequest{
					Username: request.Username,
					Password: "request.Password",
				},
			},
			want:               nil,
			wantErr:            true,
			expectErrorMessage: util.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			u := service.NewUserService(tt.repo, db)
			got, err := u.Login(tt.args.ctx, tt.args.req)
			//if !tt.wantErr(t, err, fmt.Sprintf("Login(%v, %v)", tt.args.ctx, tt.args.req)) {
			//	return
			//}
			if !tt.wantErr && got != nil {
				tt.want.Token = got.Token
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Login() = %v, want %v", got, tt.want)
			}
			if errMock := mock.ExpectationsWereMet(); errMock != nil {
				t.Errorf("there were unfulfilled expectations: %s", errMock)
			}
			assert.Equalf(t, tt.want, got, "Login(%v, %v)", tt.args.ctx, tt.args.req)
			if tt.wantErr && tt.expectErrorMessage != nil {
				assert.ErrorIs(t, err, tt.expectErrorMessage)
				assert.True(t, errors.Is(err, tt.expectErrorMessage))
			} else if tt.wantErr {
				assert.Error(t, err)
			} else if !tt.wantErr {
				assert.NoError(t, err)
				assert.Equal(t, err, tt.expectErrorMessage)
			}
		})
	}
}

func Test_userServiceImpl_Registration(t *testing.T) {
	db, mock, err1 := sqlmock.New()
	if err1 != nil {
		t.Errorf("Error sql mock, %s", err1)
	}
	defer db.Close()

	//dummy data
	dummy := model.User{
		Id:           1,
		Username:     "si.unyil",
		Fullname:     "si unyil",
		Password:     "123456",
		CreationDate: creationDate,
	}

	resp := &web.UserResponse{
		Id:       dummy.Id,
		Username: dummy.Username,
		Fullname: dummy.Fullname,
		Token:    "",
	}

	type args struct {
		ctx context.Context
		req *web.UserCreateRequest
	}
	tests := []struct {
		name               string
		mock               func()
		repo               repository.UserRepository
		args               args
		want               *web.UserResponse
		wantErr            bool
		expectErrorMessage error
	}{
		// TODO: Add test cases.
		{
			name: "registration - positive case",
			mock: func() {
				mock.ExpectBegin()

				mock.ExpectQuery(`select u.user_id, u.username, u.fullname, u.password, u.creation_date from public."users" u`).
					WithArgs(dummy.Username).
					WillReturnRows(sqlmock.NewRows([]string{"user_id", "username", "fullname", "password", "u.creation_date"}))

				mock.ExpectQuery(`INSERT INTO public."users" (.+) RETURNING user_id`).
					WithArgs(dummy.Username, dummy.Fullname, sqlmock.AnyArg() /*request.Password*/).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(dummy.Id))
				mock.ExpectCommit()
			},
			repo: repository.NewUserRepository(),
			args: args{
				ctx: context.Background(),
				req: &web.UserCreateRequest{
					Username: dummy.Username,
					Fullname: dummy.Fullname,
					Password: dummy.Password,
				},
			},
			want:               resp,
			wantErr:            false,
			expectErrorMessage: nil,
		},
		{
			name: "registration - username already use",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(`select u.user_id, u.username, u.fullname, u.password, u.creation_date from public."users" u`).
					WithArgs(dummy.Username).
					WillReturnRows(sqlmock.NewRows([]string{"user_id", "username", "fullname", "password", "u.creation_date"}).
						AddRow(dummy.Id, dummy.Username, dummy.Fullname, dummy.Password, dummy.CreationDate))
				mock.ExpectRollback()
			},
			repo: repository.NewUserRepository(),
			args: args{
				ctx: context.Background(),
				req: &web.UserCreateRequest{
					Username: dummy.Username,
					Fullname: dummy.Fullname,
					Password: dummy.Password,
				},
			},
			want:               nil,
			wantErr:            true,
			expectErrorMessage: util.ErrAlreadyUsed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			u := service.NewUserService(tt.repo, db)
			got, err := u.Registration(tt.args.ctx, tt.args.req)
			//if !tt.wantErr(t, err, fmt.Sprintf("Registration(%v, %v)", tt.args.ctx, tt.args.req)) {
			//	return
			//}

			if !tt.wantErr && got != nil {
				tt.want.Token = got.Token
			}
			if (err != nil) != tt.wantErr {
				//t.Errorf("GetById(%v, %v)", tt.args.ctx, tt.args.id)
				t.Errorf("Registration() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Registration() = %v, want %v", got, tt.want)
			}
			if errMock := mock.ExpectationsWereMet(); errMock != nil {
				t.Errorf("there were unfulfilled expectations: %s", errMock)
			}
			assert.Equalf(t, tt.want, got, "Registration(%v, %v)", tt.args.ctx, tt.args.req)
			if tt.wantErr && tt.expectErrorMessage != nil {
				assert.ErrorIs(t, err, tt.expectErrorMessage)
				assert.True(t, errors.Is(err, tt.expectErrorMessage))
			} else if tt.wantErr {
				assert.Error(t, err)
			} else if !tt.wantErr {
				assert.NoError(t, err)
				assert.Equal(t, err, tt.expectErrorMessage)
			}
		})
	}
}

func Test_userServiceImpl_Remove(t *testing.T) {
	db, mock, err1 := sqlmock.New()
	if err1 != nil {
		t.Errorf("Error sql mock, %s", err1)
	}
	defer db.Close()

	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		name               string
		mock               func()
		repo               repository.UserRepository
		args               args
		wantErr            bool
		expectErrorMessage error
	}{
		// TODO: Add test cases.
		{
			name: "delete - positive case",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM public."users"`).
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			repo: repository.NewUserRepository(),
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			wantErr:            false,
			expectErrorMessage: nil,
		},
		{
			name: "delete - user id not found",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM public."users"`).
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectRollback()
			},
			repo: repository.NewUserRepository(),
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			wantErr:            true,
			expectErrorMessage: util.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			u := service.NewUserService(tt.repo, db)
			//tt.wantErr(t, u.Remove(tt.args.ctx, tt.args.id), fmt.Sprintf("Remove(%v, %v)", tt.args.ctx, tt.args.id))
			err := u.Remove(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
			if errMock := mock.ExpectationsWereMet(); errMock != nil {
				t.Errorf("there were unfulfilled expectations: %s", errMock)
			}

			if tt.wantErr && tt.expectErrorMessage != nil {
				assert.ErrorIs(t, err, tt.expectErrorMessage)
				assert.True(t, errors.Is(err, tt.expectErrorMessage))
				assert.Error(t, err)
			} else if tt.wantErr {
				assert.Error(t, err)
			} else if !tt.wantErr {
				assert.NoError(t, err)
				assert.Equal(t, err, tt.expectErrorMessage)
			}
		})
	}
}

func Test_userServiceImpl_Update(t *testing.T) {
	db, mock, err1 := sqlmock.New()
	if err1 != nil {
		t.Errorf("Error sql mock, %s", err1)
	}
	defer db.Close()

	type args struct {
		ctx context.Context
		req *model.User
	}
	tests := []struct {
		name               string
		mock               func()
		repo               repository.UserRepository
		args               args
		want               *web.UserResponse
		wantErr            bool
		expectErrorMessage error
	}{
		// TODO: Add test cases.
		{
			name: "update - positive case",
			mock: func() {

				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE public."users"`).
					WithArgs(1, "fullname1-baru", sqlmock.AnyArg() /*,"password1-baru"*/, "username1-baru").
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			repo: repository.NewUserRepository(),
			args: args{
				ctx: context.Background(),
				req: &model.User{
					Id:       1,
					Username: "username1-baru",
					Fullname: "fullname1-baru",
					Password: "password1-baru",
				},
			},
			want: &web.UserResponse{
				Id:       1,
				Username: "username1-baru",
				Fullname: "fullname1-baru",
			},
			wantErr:            false,
			expectErrorMessage: nil,
		},
		{
			name: "update - user id not found",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE public."users"`).
					//WithArgs(1, "fullname1-baru", "password1-baru", "username1-baru").
					WithArgs(1, "fullname1-baru", sqlmock.AnyArg() /*,"password1-baru"*/, "username1-baru").
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectRollback()
			},
			repo: repository.NewUserRepository(),
			args: args{
				ctx: context.Background(),
				req: &model.User{
					Id:       1,
					Username: "username1-baru",
					Fullname: "fullname1-baru",
					Password: "password1-baru",
				},
			},
			want:               nil,
			wantErr:            true,
			expectErrorMessage: util.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			u := service.NewUserService(tt.repo, db)
			got, err := u.Update(tt.args.ctx, tt.args.req)
			//if !tt.wantErr(t, err, fmt.Sprintf("Update(%v, %v)", tt.args.ctx, tt.args.req)) {
			//	return
			//}

			if errMock := mock.ExpectationsWereMet(); errMock != nil {
				t.Errorf("there were unfulfilled expectations: %s", errMock)
			}
			assert.Equalf(t, tt.want, got, "Update(%v, %v)", tt.args.ctx, tt.args.req)
			if tt.wantErr && tt.expectErrorMessage != nil {
				assert.ErrorIs(t, err, tt.expectErrorMessage)
				assert.True(t, errors.Is(err, tt.expectErrorMessage))
			} else if tt.wantErr {
				assert.Error(t, err)
			} else if !tt.wantErr {
				assert.NoError(t, err)
				assert.Equal(t, err, tt.expectErrorMessage)
			}
		})
	}
}
