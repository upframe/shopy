package mysql

import (

	// Import driver
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/upframe/fest"
)

// InitDB establishes a connection with the database
func InitDB(user, pass, host, port, name string) (*sqlx.DB, error) {
	var (
		err error
		db  *sqlx.DB
	)

	// Prepares the connection with the database
	db, err = sqlx.Open("mysql", user+":"+pass+"@tcp("+host+":"+port+")/"+name+"?parseTime=true")
	if err != nil {
		return db, err
	}

	// Pings the database to check if the connection is OK
	err = db.Ping()
	if err != nil {
		return db, err
	}

	return db, nil
}

// updateQuery buils an update query string based on the table, the 'where' filter,
// and the name of the fields to be updated. This query is supposed to be used
// as a prepared statement
func updateQuery(table, where string, fields []string) string {
	// Starts writing the query
	query := "UPDATE " + table + " SET "

	// Add the fields that we are going to update to the queyr
	for i := range fields {
		if strings.ToLower(fields[i]) == "id" {
			continue
		}

		query += fields[i] + "=:" + fields[i] + ", "
	}

	query = strings.TrimSuffix(query, ", ")

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

// fieldsToColumns converts a struct field name defined in the domain package
// (fest) to its name on the database using a map.
func fieldsToColumns(m map[string]string, fields ...string) []string {
	columns := []string{}

	if fields[0] == fest.UpdateAll {
		columns = getAllColumns(m)
	} else {
		for i := range fields {
			columns = append(columns, m[fields[i]])
		}
	}

	return columns
}

// getAllColumns gets all of the columns of the database from the map.
func getAllColumns(m map[string]string) []string {
	columns := []string{}

	for _, value := range m {
		columns = append(columns, value)
	}

	return columns
}

// getTableCount gets the number of records in a table
func getTableCount(db *sqlx.DB, table string) (int, error) {
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
