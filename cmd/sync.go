package cmd

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var syncCommand = &cobra.Command{
	Use:     "sync [sync_url]",
	Aliases: []string{"s"},
	Short:   "Sync timezone provided a time url",
	Args:    cobra.ExactArgs(1),
	RunE:    runSync,
	Example: strings.TrimSpace(`
		gomanz sync http://remote-api/time
	`),
}

func runSync(cmd *cobra.Command, args []string) error {
	syncUrl := strings.TrimSpace(args[0])
	if len(syncUrl) == 0 {
		errors.WithStack(errors.New("sync uri argument not provided."))
	}

	if err := syncTime(syncUrl); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func syncTime(syncUrl string) error {
	resp, err := http.Get(syncUrl)
	if err != nil {
		return errors.WithStack(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.WithStack(err)
	}

	if _, err = exec.LookPath("powershell"); err != nil {
		return errors.WithStack(errors.New("date binary doesn't exist on this machine"))
	}

	command := exec.Command("powershell", []string{"Set-Date", string(body)}...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		return errors.WithStack(err)
	}

	curTime, err := time.Parse(time.RFC3339, string(body))
	if err != nil {
		return errors.WithStack(err)
	}
	windowsTimeFormat := "2 Jan 2006 15:04:05"
	curTimeWindows := curTime.Format(windowsTimeFormat)

	log.Println("date set successfully to:", curTimeWindows)

	return nil
}

func init() {
	rootCmd.AddCommand(syncCommand)
}
