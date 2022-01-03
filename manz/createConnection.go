package manz

import (
	"database/sql"
	"net/url"

	_ "github.com/denisenkom/go-mssqldb"
)

// Creates an instance of the Manzoma.
func CreateConnection(connection url.URL) (*Manzoma, error) {
	db, err := sql.Open("sqlserver", connection.String())
	if err != nil {
		return nil, err
	}

	m := &Manzoma{
		DB: db,
	}

	return m, nil
}
