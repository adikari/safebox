package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	stage  string
	config string
)

var rootCmd = &cobra.Command{
	Use:   "safebox",
	Short: "SafeBox is a secret manager CLI program",
	Long:  `A Fast and Flexible secret manager built with love by adikari in Go.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&stage, "stage", "s", "", "stage to deploy to")
	rootCmd.PersistentFlags().StringVarP(&config, "config", "c", "safebox.yml", "path to safebox configuration file")
	rootCmd.MarkPersistentFlagRequired("stage")
}

func Execute(version string) {
	rootCmd.Version = version

	if cmd, err := rootCmd.ExecuteC(); err != nil {
		if strings.Contains(err.Error(), "arg(s)") || strings.Contains(err.Error(), "usage") {
			cmd.Usage()
		}
		os.Exit(1)
	}
}
