package sqruct

import "database/sql"

// DB represents subset of database/sql.DB or database/sql.Tx.
type DB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}
