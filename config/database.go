package config

import (
	// Calls the mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// initDB establishes a connection with the database
func initDB() error {
	var err error

	// Prepares the connection with the database
	db, err = sqlx.Open("mysql", dbUser+":"+dbPass+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"?parseTime=true")
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
