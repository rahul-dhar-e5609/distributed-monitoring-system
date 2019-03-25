package datamanager

import (
	"database/sql"
	// for the initialization of the library
	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres",
		"postres://postgres:password@localhost/distributed?sslmode=disable")
	if err != nil {
		panic(err.Error())
	}
}
