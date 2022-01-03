package access

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-adodb"
)

func CreateConnection(provider, dataSrc string) (*Access, error) {
	connectionString := fmt.Sprintf("PROVIDER=%s;DATA SOURCE=%s", provider, dataSrc)
	db, err := sql.Open("adodb", connectionString)
	if err != nil {
		return nil, err
	}

	a := &Access{
		DB: db,
	}

	return a, nil
}
