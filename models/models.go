package models

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// UpdateAll is used as a placeholder to update all of the fields
const UpdateAll = "#update#"

var db *sqlx.DB

// InitDB establishes a connection with the database
func InitDB(user, pass, host, port, name string) error {
	var err error

	// Prepares the connection with the database
	db, err = sqlx.Open("mysql", user+":"+pass+"@tcp("+host+":"+port+")/"+name+"?parseTime=true")
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
func updateQuery(table, where string, fields []string) string {
	// Starts writing the query
	query := "UPDATE " + table + " SET "

	// Add the fields that we are going to update to the queyr
	for i := range fields {
		if i == len(fields)-1 {
			query += fields[i] + "=:" + fields[i]
		} else {
			query += fields[i] + "=:" + fields[i] + ", "
		}
	}

	// Finish the query adding the 'where' filter
	query += " WHERE " + where + "=:" + where
	return query
}

func insertQuery(table string, fields []string) string {
	query := "INSERT INTO " + table
	values := " VALUES ("
	vars := "("

	for i := range fields {
		if i == len(fields)-1 {
			vars += fields[i]
			values += ":" + fields[i]
		} else {
			vars += fields[i] + ", "
			values += ":" + fields[i] + ", "
		}
	}

	values += ")"
	vars += ")"
	query += " " + vars + values

	return query
}

// UniqueHash returns a SHA256 hash based on the string and on the
// current time
func UniqueHash(phrase string) string {
	data := phrase + time.Now().Format(time.ANSIC)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// GetTableRecords gets the number of records in a table
func GetTableRecords(table string) (int, error) {
	rows, err := db.Query("SELECT COUNT(*) FROM " + table)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	count := 0

	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	return count, nil
}
