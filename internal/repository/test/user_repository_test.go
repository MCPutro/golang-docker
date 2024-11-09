package test

import (
	"context"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MCPutro/golang-docker/internal/model"
	"github.com/MCPutro/golang-docker/internal/repository"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestSaveUser(t *testing.T) {

	ctx := context.Background()

	repoManager := repository.NewRepositoryManager()
	userRepository := repoManager.UserRepository()

	newUser := model.User{
		Username: "kuro",
		Fullname: "kurokami",
		Password: "asd123",
	}

	expectedId := 50

	//create database mock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock db: %v", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO public."users" (.+) RETURNING user_id`).
		WithArgs(newUser.Username, newUser.Fullname, newUser.Password).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedId))
	mock.ExpectCommit()

	//membuat Database transaction
	tx, err := db.Begin()

	//call save func from UserRepository
	userSaved, err1 := userRepository.Save(ctx, tx, &newUser)

	if err1 != nil {
		tx.Rollback()
		log.Println("rollback", err1)
	} else {
		tx.Commit()
		log.Println("commit")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NoError(t, err1)
	assert.Equal(t, expectedId, userSaved.Id)

}

func TestFindAll(t *testing.T) {
	ctx := context.Background()

	repoManager := repository.NewRepositoryManager()
	userRepository := repoManager.UserRepository()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock db: %v", err)
	}
	defer db.Close()

	users := []model.User{
		{Id: 1, Username: "user1", Fullname: "name1", CreationDate: time.Now().String()},
		{Id: 2, Username: "user2", Fullname: "name2", CreationDate: time.Now().String()},
	}

	//set expect data
	rows := sqlmock.NewRows([]string{"user_id", "username", "fullname", "u.creation_date"})
	for _, us := range users {
		rows.AddRow(us.Id, us.Username, us.Fullname, us.CreationDate)
	}

	mock.ExpectBegin()

	mock.ExpectQuery(`select u.user_id, u.username, u.fullname, u.creation_date from public."users" u`).WillReturnRows(rows)

	mock.ExpectCommit()

	//membuat Database transaction
	tx, err := db.Begin()

	//call save func from UserRepository
	findAll, err := userRepository.FindAll(ctx, tx)

	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NoError(t, err)
	assert.Len(t, findAll, 2, "invalid length data")
}

func TestFindByID(t *testing.T) {
	ctx := context.Background()

	repoManager := repository.NewRepositoryManager()
	userRepository := repoManager.UserRepository()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock db: %v", err)
	}
	defer db.Close()

	//mock data resp
	users := []model.User{
		{Id: 1, Username: "user1", Fullname: "name1"},
	}

	//set expect data
	rows := sqlmock.NewRows([]string{"user_id", "username", "Fullname", "u.creation_date"})
	for _, us := range users {
		rows.AddRow(us.Id, us.Username, us.Fullname, us.CreationDate)
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`select u.user_id, u.username, u.fullname, u.creation_date from public."users" u`).
		//select (.+) from public."users" u
		WithArgs(1).
		WillReturnRows(rows)
	mock.ExpectCommit()

	//membuat Database transaction
	tx, err := db.Begin()

	//call save func from UserRepository
	findByID, err := userRepository.FindByID(ctx, tx, 1)

	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, findByID.Id, 1)
	assert.Equal(t, findByID.Username, "user1")
	assert.Equal(t, findByID.Fullname, "name1")
}

func TestFindByUsername(t *testing.T) {
	ctx := context.Background()

	repoManager := repository.NewRepositoryManager()
	userRepository := repoManager.UserRepository()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock db: %v", err)
	}
	defer db.Close()

	//mock data resp
	users := []model.User{
		{Id: 1, Username: "user1", Fullname: "name1"},
	}

	//set expect data
	rows := sqlmock.NewRows([]string{"user_id", "username", "fullname", "password", "u.creation_date"})
	for _, us := range users {
		rows.AddRow(us.Id, us.Username, us.Fullname, us.Password, us.CreationDate)
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`select u.user_id, u.username, u.fullname, u.password, u.creation_date from public."users" u`).
		WithArgs(users[0].Username).
		WillReturnRows(rows)
	mock.ExpectCommit()

	//membuat Database transaction
	tx, err := db.Begin()

	//call save func from UserRepository
	findByUsername, err := userRepository.FindByUsername(ctx, tx, "user1")

	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, findByUsername.Id, 1)
	assert.Equal(t, findByUsername.Username, "user1")
	assert.Equal(t, findByUsername.Fullname, "name1")
}

func TestUpdate(t *testing.T) {

	newUser := model.User{
		Id:       4,
		Fullname: "empat",
		Password: "3mpat",
	}

	/*init database*/
	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	//set mock
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE public."users"`).
		WithArgs(newUser.Id, newUser.Fullname, newUser.Password, newUser.Username).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	//membuat Database transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}

	repoManager := repository.NewRepositoryManager()
	userRepository := repoManager.UserRepository()

	err = userRepository.Update(context.Background(), tx, &newUser)

	if err != nil {
		tx.Rollback()
		log.Println(err)
		log.Println("rollback")
	} else {
		tx.Commit()
		log.Println("commit")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NoError(t, err)

}

func TestDeleteUser_case_positive(t *testing.T) {

	id := 4

	/*init database*/
	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	//set mock
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM public."users"`).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	//membuat Database transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}

	repoManager := repository.NewRepositoryManager()
	userRepository := repoManager.UserRepository()

	err = userRepository.Delete(context.Background(), tx, id)

	if err == nil {
		tx.Commit()
		log.Println("commit")
	} else {
		tx.Rollback()
		log.Println("rollback, err : ", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	//assert.NoError(t, err)
}

func TestDeleteUser_case_negative(t *testing.T) {

	id := 4

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM public."users"`).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()

	ctx := context.Background()
	tx, err := db.Begin()
	if err != nil {
		t.Errorf("an error '%s' was not expected when begin database transactional", err)
	}

	repoManager := repository.NewRepositoryManager()
	userRepository := repoManager.UserRepository()

	err = userRepository.Delete(ctx, tx, id)

	if err != nil {
		tx.Rollback()
		log.Println("rollback database transaction, err : ", err)
	} else {
		tx.Commit()
		log.Println("commit database transaction")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
