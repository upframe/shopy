package mysql

import (
	"net/http"

	// Import driver
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

// insertQuery buils an insert query string based on the table and the name of
// the fields to be inserted. This query is supposed to be used
// as a prepared statement
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

func fieldsToColumns(m map[string]string, fields ...string) []string {
	columns := []string{}

	if fields[0] == UpdateAll {
		columns = getAllColumns(m)
	} else {
		for i := range fields {
			columns = append(columns, m[fields[i]])
		}
	}

	return columns
}

func getAllColumns(m map[string]string) []string {
	columns := []string{}

	for _, value := range m {
		columns = append(columns, value)
	}

	return columns
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
