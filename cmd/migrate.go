package cmd

import (
	"fmt"
	"log"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/osamaadam/gomanz/access"
	"github.com/osamaadam/gomanz/manz"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:     "migrate [access.mdb path]",
	Args:    cobra.NoArgs,
	Aliases: []string{"m", "mig", "update-manzoma"},
	Short:   "Migrate soldiers from Access to Manzoma",
	RunE:    migrateRun,
	Example: strings.TrimSpace(`
		gomanz migrate --host http://remote-pc -u sa -p password1 --port 1433
	`),
}

type MigrateOpts struct {
	Username      string
	Password      string
	Host          string
	Port          uint
	AccessSrcPath string
	Marhla        uint
}

var migrateOpts MigrateOpts

func migrateRun(cmd *cobra.Command, args []string) error {
	manzConnection := url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword(migrateOpts.Username, migrateOpts.Password),
		Host:   fmt.Sprintf("%s:%d", migrateOpts.Host, migrateOpts.Port),
	}

	log.Println("creating a connection to Manzoma")
	m, err := manz.CreateConnection(manzConnection)
	if err != nil {
		return err
	}
	defer m.DB.Close()

	log.Println("finding the last registered soldier in Manzoma")
	lastSeglNo, err := m.FindLastSoldier(20221)
	if err != nil {
		return errors.WithStack(err)
	}

	dataSrc := filepath.Clean(migrateOpts.AccessSrcPath)

	log.Printf("creating a connection to %s\n", dataSrc)
	a, err := access.CreateConnection("Microsoft.ACE.OLEDB.12.0", dataSrc)
	if err != nil {
		return errors.WithStack(err)
	}
	defer a.DB.Close()

	sqlQuery := `
		select mrhla, segl_no, military_no, soldier_name, s.moahel_code, s.etgah as etgah_code, m.moahel_name as moahel,
			e.etgah as etgah, governorate_fk as gov
		from 
		(((
			src_soldiers s
			left join moahel_type m on s.moahel_code = m.moahel_code)
			left join etgah e on s.etgah = e.` + "`etgah c`" + `)
		)
		where segl_no > ?
			and mrhla = ?
	`

	log.Println("querying the access database for unregistered soldiers")
	rows, err := a.DB.Query(sqlQuery, lastSeglNo, migrateOpts.Marhla)
	if err != nil {
		return errors.WithStack(err)
	}
	defer rows.Close()

	var soldiers []manz.Soldier

	for rows.Next() {
		var soldier manz.Soldier

		err := rows.Scan(&soldier.Marhla, &soldier.Segl_no, &soldier.Military_no, &soldier.S_name, &soldier.Moahel_code, &soldier.Etgah_code,
			&soldier.Moahel, &soldier.Etgah, &soldier.Gov)
		if err != nil {
			return errors.WithStack(err)
		}
		soldiers = append(soldiers, soldier)
	}

	if len(soldiers) <= 0 {
		fmt.Println("manzoma database already up to date")
		return nil
	}

	log.Printf("found %d new soldiers starting at segl_no %d\n", len(soldiers), soldiers[0].Segl_no)
	log.Println("migrating new soldiers from the access database to Manzoma")
	if err := m.InsertMany(&soldiers); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func guessMarhla() (uint, error) {
	curTime := time.Now()
	curYear := curTime.Year()
	curMonth := curTime.Month()

	curMarhla := fmt.Sprintf("%d", curYear)

	switch curMonth {
	case time.January, time.February, time.March:
		curMarhla += "1"
	case time.April, time.May, time.June:
		curMarhla += "2"
	case time.July, time.August, time.September:
		curMarhla += "3"
	case time.October, time.November, time.December:
		curMarhla += "4"
	}

	intMarhla, err := strconv.Atoi(curMarhla)
	if err != nil {
		return 0, err
	}

	return uint(intMarhla), nil
}

func init() {
	curMarhla, _ := guessMarhla()
	migrateCmd.Flags().StringVarP(&migrateOpts.AccessSrcPath, "access", "a", "", "path to the .mdb file")
	migrateCmd.Flags().StringVar(&migrateOpts.Host, "host", "localhost", "host of the mssql server")
	migrateCmd.Flags().UintVar(&migrateOpts.Port, "port", 1433, "port of the running mssql database server")
	migrateCmd.Flags().StringVarP(&migrateOpts.Username, "user", "u", "", "username for the mssql database")
	migrateCmd.Flags().StringVarP(&migrateOpts.Password, "pass", "p", "", "password for the mssql database")
	migrateCmd.Flags().UintVarP(&migrateOpts.Marhla, "marhla", "m", curMarhla, "the marhla to migrate")

	migrateCmd.MarkFlagRequired("access")
	migrateCmd.MarkFlagRequired("user")
	migrateCmd.MarkFlagRequired("pass")

	rootCmd.AddCommand(migrateCmd)
}
