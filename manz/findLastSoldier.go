package manz

import (
	"database/sql"

	"github.com/pkg/errors"
)

// Finds the `segl_no` of the last soldier registered on Manzoma.
//
// This doesn't find all missing `segl_no`s in the database, only returns the last one.
func (m *Manzoma) FindLastSoldier(marhla int) (int, error) {
	sqlQuery := `
		select top (1)
			segl_no
		from [manzoma].[dbo].[soldiers]
		where marhla = @marhla
		order by segl_no desc
	`

	rows, err := m.DB.Query(sqlQuery, sql.Named("marhla", marhla))
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
