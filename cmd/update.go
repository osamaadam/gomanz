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

var updateCmd = &cobra.Command{
	Use:     "update [access.mdb path]",
	Args:    cobra.NoArgs,
	Aliases: []string{"u", "insert", "update-manzoma"},
	Short:   "Update manzoma with the data in the access file provided",
	RunE:    updateRun,
	Example: strings.TrimSpace(`
		gomanz update --host http://remote-pc -u sa -p password1 --port 1433
	`),
}

type UpdateCmdOpts struct {
	Username      string
	Password      string
	Host          string
	Port          uint
	AccessSrcPath string
	Marhla        uint
}

var updateCmdOpts UpdateCmdOpts

func updateRun(cmd *cobra.Command, args []string) error {
	manzConnection := url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword(updateCmdOpts.Username, updateCmdOpts.Password),
		Host:   fmt.Sprintf("%s:%d", updateCmdOpts.Host, updateCmdOpts.Port),
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

	dataSrc := filepath.Clean(updateCmdOpts.AccessSrcPath)

	log.Printf("Creating a connection to %s\n", dataSrc)
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
	rows, err := a.DB.Query(sqlQuery, lastSeglNo, updateCmdOpts.Marhla)
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
	updateCmd.Flags().StringVarP(&updateCmdOpts.AccessSrcPath, "access", "a", "", "path to the .mdb file")
	updateCmd.Flags().StringVar(&updateCmdOpts.Host, "host", "localhost", "host of the mssql server")
	updateCmd.Flags().UintVar(&updateCmdOpts.Port, "port", 1433, "port of the running mssql database server")
	updateCmd.Flags().StringVarP(&updateCmdOpts.Username, "user", "u", "", "username for the mssql database")
	updateCmd.Flags().StringVarP(&updateCmdOpts.Password, "pass", "p", "", "password for the mssql database")
	updateCmd.Flags().UintVarP(&updateCmdOpts.Marhla, "marhla", "m", curMarhla, "the marhla to migrate")

	updateCmd.MarkFlagRequired("access")
	updateCmd.MarkFlagRequired("user")
	updateCmd.MarkFlagRequired("pass")

	rootCmd.AddCommand(updateCmd)
}
