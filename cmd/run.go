package cmd

import (
	"log"

	c "github.com/adikari/safebox/v2/config"
	"github.com/spf13/cobra"
)

// runCmd represents the exec command
var runCmd = &cobra.Command{
	Use:     "run",
	Short:   "Deploys all configurations specified in config file",
	RunE:    execute,
	Example: `TODO: run command example`,
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func execute(cmd *cobra.Command, args []string) error {
	params := c.LoadParam{
		Path:  config,
		Stage: stage,
	}

	config, err := c.Load(params)

	log.Printf("%v", config)
	log.Printf("%v", err)
	return nil
}
