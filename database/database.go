package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var commandChannel = make(chan string)
var db *sql.DB

func Init(filename string) error {
	db, err := open(filename)
	print(db)
	return err
}

func NewVisit() {

}

func open(filename string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filename)
	return db, err
}
