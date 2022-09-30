package cmd

import (
	"log"

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
	_, err := getStore()
	log.Printf("%v", Config)
	log.Printf("%v", err)
	return nil
}
