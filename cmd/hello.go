package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var helloCmd = &cobra.Command{
	Use:    "hello",
	Short:  "debugging command",
	Args:   cobra.ExactArgs(1),
	RunE:   hello,
	Hidden: true,
}

func hello(cmd *cobra.Command, args []string) error {
	fmt.Println("Hello,", args[0])

	return nil
}

func init() {
	rootCmd.AddCommand(helloCmd)
}
