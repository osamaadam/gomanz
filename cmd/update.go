package cmd

import (
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/osamaadam/gomanz/access"
	"github.com/osamaadam/gomanz/manz"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update [access.mdb path]",
	Args:  cobra.NoArgs,
	Short: "update manzoma with the data in the access file provided",
	RunE:  updateRun,
}

type UpdateCmdOpts struct {
	Username      string
	Password      string
	Host          string
	Port          uint
	AccessSrcPath string
}

var updateCmdOpts UpdateCmdOpts

func updateRun(cmd *cobra.Command, args []string) error {
	manzConnection := url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword(updateCmdOpts.Username, updateCmdOpts.Password),
		Host:   fmt.Sprintf("%s:%d", updateCmdOpts.Host, updateCmdOpts.Port),
	}

	manzDB, err := manz.Connect(manzConnection)
	if err != nil {
		return errors.WithStack(err)
	}
	defer manzDB.Close()

	lastSeglNo, err := manz.FindLastSoldier(*manzDB, 20221)
	if err != nil {
		return errors.WithStack(err)
	}

	dataSrc := filepath.Clean(updateCmdOpts.AccessSrcPath)

	accessDB, err := access.Connect("Microsoft.ACE.OLEDB.12.0", dataSrc)
	if err != nil {
		return errors.WithStack(err)
	}
	defer accessDB.Close()

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
	`

	rows, err := accessDB.Query(sqlQuery, lastSeglNo)
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

	if err := manz.InsertMany(manzDB, &soldiers); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func init() {
	updateCmd.Flags().StringVarP(&updateCmdOpts.AccessSrcPath, "access", "a", "", "path to the .mdb file")
	updateCmd.Flags().StringVar(&updateCmdOpts.Host, "host", "localhost", "host of the mssql server")
	updateCmd.Flags().UintVar(&updateCmdOpts.Port, "port", 1433, "port of the running mssql database server")
	updateCmd.Flags().StringVarP(&updateCmdOpts.Username, "user", "u", "", "username for the mssql database")
	updateCmd.Flags().StringVarP(&updateCmdOpts.Password, "pass", "p", "", "password for the mssql database")

	updateCmd.MarkFlagRequired("access")
	updateCmd.MarkFlagRequired("user")
	updateCmd.MarkFlagRequired("pass")

	rootCmd.AddCommand(updateCmd)
}
