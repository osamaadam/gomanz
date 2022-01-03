package manz

import (
	"database/sql"
	"net/url"

	_ "github.com/denisenkom/go-mssqldb"
)

func Connect(connection url.URL) (*sql.DB, error) {
	return sql.Open("sqlserver", connection.String())
}
