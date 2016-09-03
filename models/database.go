package models

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/jmoiron/sqlx"
)

const UpdateAll = "#update#"

var db *sqlx.DB

// InitDB establishes a connection with the database
func InitDB(database *sqlx.DB) {
	db = database
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
			query += fields[i] + ":" + fields[i] + ", "
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
