package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "gomanz",
	Short: "manzoma multipurpose utility",
	Args: cobra.NoArgs,
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
