package cmd

import (
	"github.com/adikari/safebox/v2/store"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// runCmd represents the exec command
var runCmd = &cobra.Command{
	Use:     "run",
	Short:   "Deploys all configurations specified in config file",
	RunE:    run,
	Example: `TODO: run command example`,
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func run(cmd *cobra.Command, args []string) error {
	config, err := loadConfig()

	if err != nil {
		return errors.Wrap(err, "failed to load config")
	}

	store, err := store.GetStore(config.Provider)

	if err != nil {
		return errors.Wrap(err, "failed to instantiate store")
	}

	err = store.PutMany(config.Configs)

	if err != nil {
		return errors.Wrap(err, "failed to write param")
	}

	return nil
}
