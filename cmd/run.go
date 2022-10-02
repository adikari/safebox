package cmd

import (
	"fmt"
	"log"

	"github.com/adikari/safebox/v2/store"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	removeOrphans bool
	prompt        string

	deployCmd = &cobra.Command{
		Use:     "deploy",
		Short:   "Deploys all configurations specified in config file",
		RunE:    deploy,
		Example: `TODO: deploy command example`,
	}
)

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.Flags().BoolVarP(&removeOrphans, "remove-orphans", "r", false, "remove orphan configurations")
	deployCmd.Flags().StringVarP(&prompt, "prompt", "p", "missing", "prompt for configurations (missing or all)")
}

func deploy(cmd *cobra.Command, args []string) error {
	config, err := loadConfig()

	if prompt != "all" && prompt != "missing" {
		return errors.New("value for prompt must be \"all\" or \"missing\"")
	}

	if removeOrphans {
		log.Panic("remove orphans flag is not implemented")
	}

	if err != nil {
		return errors.Wrap(err, "failed to load config")
	}

	store, err := store.GetStore(config.Provider)

	if err != nil {
		return errors.Wrap(err, "failed to instantiate store")
	}

	for i, c := range config.Configs {
		if c.Value == "" {
			config.Configs[i].Value = promptConfig(c)
		}
	}

	err = store.PutMany(config.Configs)

	if err != nil {
		return errors.Wrap(err, "failed to write param")
	}

	return nil
}

func promptConfig(config store.ConfigInput) string {
	validate := func(input string) error {
		if len(input) < 1 {
			return fmt.Errorf("%s must not be empty", config.Name)
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    config.Name,
		Validate: validate,
	}

	result, _ := prompt.Run()

	return result
}
