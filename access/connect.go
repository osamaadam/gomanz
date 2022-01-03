package access

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-adodb"
)

func Connect(provider, dataSrc string) (*sql.DB, error) {
	connectionString := fmt.Sprintf("PROVIDER=%s;DATA SOURCE=%s", provider, dataSrc)
	return sql.Open("adodb", connectionString)
}
