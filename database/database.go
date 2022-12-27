package database

import (
	"database/sql"

	"github.com/backend/utils"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func DB() *sql.DB {
	if db == nil {
		d, err := sql.Open("sqlite3", "maindb.db")
		db = d
		utils.HandleError(err)
	}
	return db
}
