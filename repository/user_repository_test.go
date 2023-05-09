package repository

import (
	"context"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	conf "github.com/MCPutro/golang-docker/config"
	"github.com/MCPutro/golang-docker/model"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"
)

func setDatabase() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalln("Error loading .env file")
		return
	}

	conf.DB_Host = os.Getenv("POSTGRES_HOSTNAME")
	conf.DB_Pass = os.Getenv("POSTGRES_PASSWORD")
	conf.DB_Username = os.Getenv("POSTGRES_USER")
	conf.DB_Name = os.Getenv("POSTGRES_DB")
	conf.DB_Port = os.Getenv("POSTGRES_DB_PORT")
}

func TestSaveUser(t *testing.T) {

	ctx := context.Background()

	userRepository := NewUserRepositoryImpl()

	newUser := model.User{
		Username: "kuro",
		FullName: "kurokami",
		Password: "asd123",
	}

	//create database mock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock db: %v", err)
	}
	defer db.Close()

	mock.ExpectBegin()

	// expect insert query
	mock.ExpectExec(`INSERT INTO public."users"`).
		WithArgs(newUser.Username, newUser.FullName, newUser.Password).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	//membuat Database transaction
	tx, err := db.Begin()

	//call save func from UserRepository
	err1 := userRepository.Save(ctx, tx, &newUser)

	if err1 != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestFindAll(t *testing.T) {
	ctx := context.Background()

	userRepository := NewUserRepositoryImpl()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock db: %v", err)
	}
	defer db.Close()

	users := []model.User{
		{Id: 1, Username: "user1", FullName: "name1", CreationDate: time.Now().String()},
		{Id: 2, Username: "user2", FullName: "name2", CreationDate: time.Now().String()},
	}

	//set expect data
	rows := sqlmock.NewRows([]string{"user_id", "username", "fullname", "u.creation_date"})
	for _, user := range users {
		rows.AddRow(user.Id, user.Username, user.FullName, user.CreationDate)
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

	userRepository := NewUserRepositoryImpl()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock db: %v", err)
	}
	defer db.Close()

	//mock data resp
	users := []model.User{
		{Id: 1, Username: "user1", FullName: "name1"},
	}

	//set expect data
	rows := sqlmock.NewRows([]string{"user_id", "username", "fullname", "u.creation_date"})
	for _, user := range users {
		rows.AddRow(user.Id, user.Username, user.FullName, user.CreationDate)
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

	//assert.NoError(t, err)
	assert.Equal(t, findByID.Id, 1)
	//assert.Equal(t, findByID.Username, "user1")
	//assert.Equal(t, findByID.FullName, "name1")
}

func TestFindByUsername(t *testing.T) {
	ctx := context.Background()

	userRepository := NewUserRepositoryImpl()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock db: %v", err)
	}
	defer db.Close()

	//mock data resp
	users := []model.User{
		{Id: 1, Username: "user1", FullName: "name1"},
	}

	//set expect data
	rows := sqlmock.NewRows([]string{"user_id", "username", "fullname", "u.creation_date"})
	for _, user := range users {
		rows.AddRow(user.Id, user.Username, user.FullName, user.CreationDate)
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`select u.user_id, u.username, u.fullname, u.creation_date from public."users" u`).
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
	assert.Equal(t, findByUsername.FullName, "name1")
}

func TestUpdate(t *testing.T) {

	user := model.User{
		Id:       4,
		FullName: "empat",
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
		WithArgs(user.Id, user.FullName, user.Password).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	//membuat Database transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}

	userRepository := NewUserRepositoryImpl()
	err = userRepository.Update(context.Background(), tx, &user)

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

	userRepository := NewUserRepositoryImpl()

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

	userRepository := NewUserRepositoryImpl()
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
