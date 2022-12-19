package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Query struct {
	Id     string
	Topic  string
	Option int
}

var commandChannel = make(chan Query)
var db *sql.DB

var Options map[string][]string

func Init(filename string) error {
	Options = map[string][]string{
		"Borrow a Device": {"Forgot device at home", "Forgot charger at home", "Lost Charger", "I don't have a school device"},
		"IT Help":         {"Software Install", "Wifi help", "Other"},
		"Broken Device":   {"Broken Screen", "Won't turn on/charge", "Broken elsewhere"},
	}
	var err error
	db, err = sql.Open("sqlite3", filename)
	_ = db
	if err != nil {
		log.Panic(err)
	}

	go pollCheckIns()

	return err
}

func CheckIn(q Query) {
	commandChannel <- q
}

func pollCheckIns() {
	defer db.Close()

	for {
		query := <-commandChannel
		switch query.Topic {
		case "Borrow a Device":
			queryBorrow(query.Id, query.Option)
		case "IT Help":
			queryHelp(query.Id, query.Option)
		case "Broken Device":
			queryBroken(query.Id, query.Option)
		}
	}
}

func queryBorrow(id string, reason int) {
	// map reason to column
	reasonCol := ""

	switch reason {
	case 0:
		reasonCol = "ForgotDevice"
	case 1:
		reasonCol = "ForgotCharger"
	case 2:
		reasonCol = "LostCharger"
	case 3:
		reasonCol = "NotIssuedDevice"
	}
	query := fmt.Sprintf("SELECT ID FROM BorrowDevice WHERE id='%s'", id)
	result, err := db.Query(query)
	exists := ""

	if err != nil {
		panic(err)
	}

	for result.Next() {
		result.Scan(&exists)
	}

	if exists == "" {
		query = fmt.Sprintf("INSERT INTO BorrowDevice VALUES ('%s', '0', '0', '0', '0')", id)
		_, err = db.Exec(query)
		if err != nil {
			panic(err)
		}
	}

	query = fmt.Sprintf("UPDATE BorrowDevice SET %s=%s + 1 WHERE ID='%s'", reasonCol, reasonCol, id)
	_, err = db.Exec(query)

	if err != nil {
		log.Fatal(err)
	}
}

func queryBroken(id string, reason int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered. Error:\n", r)
		}
	}()

	// map reason to column
	reasonCol := ""

	switch reason {
	case 0:
		reasonCol = "BrokenScreen"
	case 1:
		reasonCol = "WontTurnOnCharge"
	case 2:
		reasonCol = "BrokenOther"
	}
	query := fmt.Sprintf("SELECT ID FROM BrokenDevice WHERE id='%s'", id)
	result, err := db.Query(query)
	exists := ""

	for result.Next() {
		result.Scan(&exists)
	}

	if exists == "" {
		query = fmt.Sprintf("INSERT INTO BrokenDevice VALUES ('%s', '0', '0', '0')", id)
		_, err = db.Exec(query)
		if err != nil {
			panic(err)
		}
	}

	query = fmt.Sprintf("UPDATE BrokenDevice SET %s=%s + 1 WHERE ID='%s'", reasonCol, reasonCol, id)
	_, err = db.Exec(query)

	if err != nil {
		log.Fatal(err)
	}
}

func queryHelp(id string, reason int) {
	// map reason to column
	reasonCol := ""

	switch reason {
	case 0:
		reasonCol = "SoftwareInstall"
	case 1:
		reasonCol = "WifiHelp"
	case 2:
		reasonCol = "Other"
	}
	query := fmt.Sprintf("SELECT ID FROM ITHelp WHERE id='%s'", id)
	result, err := db.Query(query)
	exists := ""

	for result.Next() {
		result.Scan(&exists)
	}

	if exists == "" {
		query = fmt.Sprintf("INSERT INTO ITHelp VALUES ('%s', '0', '0', '0')", id)
		_, err = db.Exec(query)
		if err != nil {
			panic(err)
		}
	}

	query = fmt.Sprintf("UPDATE ITHelp SET %s=%s + 1 WHERE ID='%s'", reasonCol, reasonCol, id)
	_, err = db.Exec(query)

	if err != nil {
		log.Fatal(err)
	}
}
