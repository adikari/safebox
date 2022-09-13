package cmd

import (
	c "github.com/adikari/safebox/v2/config"
	"github.com/spf13/cobra"
)

// runCmd represents the exec command
var runCmd = &cobra.Command{
	Use:     "run",
	Short:   "Deploys all configurations specified in config file",
	RunE:    execute,
	Example: `TODO`,
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func execute(cmd *cobra.Command, args []string) error {
	c.LoadConfig()
	return nil
}
