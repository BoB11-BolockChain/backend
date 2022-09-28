package utils

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func DB() *sql.DB {
	if db == nil {
		d, err := sql.Open("sqlite3", "maindb.db")
		db = d
		HandleError(err)
	}
	return db
}

func SaveUser(query string) sql.Result {
	res, err := DB().Exec(query)
	HandleError(err)
	return res
}

func CloseDB() {
	db.Close()
}
