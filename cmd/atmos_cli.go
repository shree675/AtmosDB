package cmd

import (
	"github.com/spf13/cobra"

	"atmosdb/cli"
)

var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "Start an AtmosDB CLI session",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cli.Run(args[0])
	},
}

func init() {
	rootCmd.AddCommand(cliCmd)
}
