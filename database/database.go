package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var commandChannel = make(chan string)
var db *sql.DB

func Init(filename string) error {
	var id, forgotD, forgotC, lostC = 0, 0, 0, 0
	db, err := open(filename)
	if err != nil{
		log.Panic(err)
	}
	rows, err := db.Query("SELECT * FROM BorrowDevice WHERE ID='000036122'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id, &forgotD, &forgotC, &lostC)
		if err != nil {
			log.Fatal(err)
		}
		print(id, forgotD, forgotC, lostC)
	}

	return err
}

func NewVisit() {

}

func open(filename string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filename)
	return db, err
}
