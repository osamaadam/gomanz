package manz

import (
	"database/sql"

	"github.com/pkg/errors"
)

func FindLastSoldier(manzDB sql.DB, marhla int) (int, error) {
	sqlQuery := `
		select top (1)
			segl_no
		from [manzoma].[dbo].[soldiers]
		where marhla = @marhla
		order by segl_no desc
	`

	rows, err := manzDB.Query(sqlQuery, sql.Named("marhla", marhla))
	if err != nil {
		return 0, errors.WithStack(err)
	}

	for rows.Next() {
		var segl_no int

		rows.Scan(&segl_no)

		return segl_no, nil
	}

	return 0, nil
}
