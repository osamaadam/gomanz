package manz

import (
	"fmt"
	"log"
	"strings"

	"github.com/pkg/errors"
)

var (
	ErrNoSoldiers = errors.New("no soldiers provided")
)

// Inserts many soldiers at once.
func (m *Manzoma) InsertMany(soldiers *[]Soldier) error {
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
		soldier.Selah = "اداره الاشاره"
		soldier.Note = "-"

		sqlInsert := fmt.Sprintf("(%d, %d, %d, N'%s', %d, %d, N'%s', N'%s', N'%s', N'%s', %d, N'%s')",
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

	res, err := m.DB.Exec(insertQuery)
	if err != nil {
		return errors.WithStack(err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.WithStack(err)
	}

	log.Printf("%d rows inserted\n", rowsAffected)

	return nil
}
