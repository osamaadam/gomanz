package manz

import "database/sql"

type Manzoma struct {
	DB *sql.DB
}

// Soldier according to the Manzoma's database.
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
