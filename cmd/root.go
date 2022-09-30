package cmd

import (
	"fmt"
	"os"
	"strings"

	c "github.com/adikari/safebox/v2/config"
	"github.com/adikari/safebox/v2/store"
	"github.com/spf13/cobra"
)

var (
	stage        string
	pathToConfig string
	Config       *c.Config
)

var rootCmd = &cobra.Command{
	Use:              "safebox",
	Short:            "SafeBox is a secret manager CLI program",
	Long:             `A Fast and Flexible secret manager built with love by adikari in Go.`,
	PersistentPreRun: prerun,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&stage, "stage", "s", "", "stage to deploy to")
	rootCmd.PersistentFlags().StringVarP(&pathToConfig, "config", "c", "safebox.yml", "path to safebox configuration file")
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

func prerun(cmd *cobra.Command, args []string) {
	c, _ := loadConfig()
	Config = c
}

func loadConfig() (*c.Config, error) {
	params := c.LoadParam{
		Path:  pathToConfig,
		Stage: stage,
	}

	return c.Load(params)
}

func getStore() (store.Store, error) {
	switch Config.Provider {
	case "ssm":
		return store.NewSSMStore()
	default:
		return nil, fmt.Errorf("invalid provider `%s`", Config.Provider)
	}
}
