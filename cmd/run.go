package cmd

import (
	"github.com/pkg/errors"
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
	store, err := getStore()

	if err != nil {
		return errors.Wrap(err, "failed to instantiate store")
	}

	err = store.PutMany(Config.Configs)

	if err != nil {
		return errors.Wrap(err, "failed to write param")
	}

	return nil
}
