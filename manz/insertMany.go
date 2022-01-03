package manz

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type Soldier struct {
	Marhla      uint
	Segl_no     uint
	Military_no uint
	S_name      string
	Moahel_code uint
	Etgah_code  uint
	Moahel      string
	Etgah       string
	Selah       string
	Gov         int
	Note        string
	Tawzeaa     string
}

var (
	ErrNoSoldiers = errors.New("no soldiers provided")
)

func InsertMany(manzDB *sql.DB, soldiers *[]Soldier) error {
	if len(*soldiers) <= 0 {
		return errors.WithStack(ErrNoSoldiers)
	}
	insertQuery := `
		insert into [Manzoma].[dbo].[Soldiers]
		([marhla],
		[segl_no],
		[military_no],
		[s_name],
		[moahel_code],
		[etgah_code],
		[moahel],
		[etgah],
		[tawzeaa],
		[selah],
		[gov],
		[note])
		values
	`

	var soldiersSql []string

	for _, soldier := range *soldiers {
		soldier.Tawzeaa = "بدون توزيع"
		soldier.Selah = "ادارة الاشارة"
		soldier.Note = "-"

		sqlInsert := fmt.Sprintf("(%d, %d, %d, N'%s', %d, %d, N'%s', N'%s', N'%s', N'%s', %d, N'%s')\n",
			soldier.Marhla,
			soldier.Segl_no,
			soldier.Military_no,
			soldier.S_name,
			soldier.Moahel_code,
			soldier.Etgah_code,
			soldier.Moahel,
			soldier.Etgah,
			soldier.Tawzeaa,
			soldier.Selah,
			soldier.Gov,
			soldier.Note)

		soldiersSql = append(soldiersSql, sqlInsert)
	}

	insertQuery += strings.Join(soldiersSql, ", ")

	res, err := manzDB.Exec(insertQuery)
	if err != nil {
		return errors.WithStack(err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.WithStack(err)
	}

	fmt.Printf("%d rows inserted\n", rowsAffected)

	return nil
}
