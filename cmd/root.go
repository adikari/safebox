package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "SafeBox",
	Short: "SafeBox is a secret manager CLI program",
	Long:  `A Fast and Flexible secret manager built with love by adikari in Go.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	if cmd, err := RootCmd.ExecuteC(); err != nil {
		if strings.Contains(err.Error(), "arg(s)") || strings.Contains(err.Error(), "usage") {
			cmd.Usage()
		}
		os.Exit(1)
	}
}
