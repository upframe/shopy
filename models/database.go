package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

// InitDB establishes a connection with the database
func InitDB(username, password, address, port, dbname string) error {
	var err error

	// Prepares the connection with the database
	db, err = sqlx.Open("mysql", username+":"+password+"@tcp("+address+":"+port+")/"+dbname+"?parseTime=true")
	if err != nil {
		return err
	}

	// Pings the database to check if the connection is OK
	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}
