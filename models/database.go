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

// updateQuery buils an update query string based on the table, the 'where' filter,
// and the name of the fields to be updated. This query is supposed to be used
// as a prepared statement
func updateQuery(table, where, what string, fields []string) string {
	// Starts writing the query
	query := "UPDATE " + table + " SET "

	// Add the fields that we are going to update to the queyr
	for i := range fields {
		if i == len(fields)-1 {
			query += fields[i] + "=:" + fields[i]
		} else {
			query += fields[i] + ":" + fields[i] + ", "
		}
	}

	// Finish the query adding the 'where' filter
	query += " WHERE " + where + "=" + what
	return query
}
