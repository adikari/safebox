package cmd

import (
	"os"
	"strings"

	c "github.com/adikari/safebox/v2/config"
	"github.com/spf13/cobra"
)

var (
	stage        string
	pathToConfig string
	TimeFormat   = "2006-01-02 15:04:05"
)

var rootCmd = &cobra.Command{
	Use:          "safebox",
	Short:        "SafeBox is a secret manager CLI program",
	Long:         `A Fast and Flexible secret manager built with love by adikari in Go.`,
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, _ []string) {
		cmd.Usage()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&stage, "stage", "s", "", "stage to deploy to")

	rootCmd.PersistentFlags().StringVarP(&pathToConfig, "config", "c", "", "path to safebox configuration file")
	rootCmd.MarkFlagFilename("config")
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

func loadConfig() (*c.Config, error) {
	return c.Load(c.LoadConfigInput{
		Path:  pathToConfig,
		Stage: stage,
	})
}
