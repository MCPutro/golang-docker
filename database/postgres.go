package database

import (
	"database/sql"
	"fmt"
	"github.com/MCPutro/golang-docker/config"
	_ "github.com/lib/pq"

	"log"
	"time"
)

// InitDatabase create database connection
func InitDatabase() (*sql.DB, error) {
	var db *sql.DB
	var err error

	for i := 1; i <= 5; i++ {
		dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s sslmode=%s TimeZone=Asia/Jakarta",
			config.DB_Host, config.DB_Username, config.DB_Pass, config.DB_Port, config.DB_Name, "disable")

		db, err = sql.Open("postgres", dsn)

		if err = db.Ping(); err != nil {
			log.Printf("error create db connection [rety %d times], error message : %s will retry in 5 second.", i, err)
		} else {
			db.SetMaxIdleConns(5)
			db.SetMaxOpenConns(100)
			db.SetConnMaxLifetime(60 * time.Minute)
			db.SetConnMaxIdleTime(10 * time.Minute)

			//out from looping
			break
		}

		time.Sleep(5 * time.Second)

	}

	//test ping
	if err = db.Ping(); err != nil {
		log.Fatalf("error ping db connection, err: %s", err)
		return nil, err
	}

	return db, nil
}
