package util

import (
	"database/sql"
	"log"
)

func CommitOrRollback(err error, tx *sql.Tx) {
	if err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			log.Println("[ERROR] failed rollback : ", err2)
		} else {
			log.Println("[INFO] success rollback")
		}
	} else {
		if err2 := tx.Commit(); err2 != nil {
			log.Println("[ERROR] failed commit : ", err2)
		} else {
			log.Println("[INFO] success commit")
		}
	}
}
