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
	deployCmd.Flags().StringVarP(&prompt, "prompt", "p", "", "prompt for configurations (missing or all)")
}

func deploy(cmd *cobra.Command, args []string) error {
	config, err := loadConfig()

	if prompt != "" && prompt != "all" && prompt != "missing" {
		return errors.New("value for prompt must be \"all\" or \"missing\"")
	}

	if removeOrphans {
		log.Panic("remove orphans flag is not implemented")
	}

	if err != nil {
		return errors.Wrap(err, "failed to load config")
	}

	st, err := store.GetStore(config.Provider)

	if err != nil {
		return errors.Wrap(err, "failed to instantiate store")
	}

	all, _ := st.GetMany(config.All)

	missing := getMissing(config.Secrets, all)

	if len(missing) > 0 && prompt == "" {
		return errors.New("config values missing. run deploy with \"--prompt\" flag")
	}

	configsToDeploy := []store.ConfigInput{}

	// prompt for missing secrets
	if prompt == "missing" {
		for _, c := range missing {
			if c.Value == "" {
				configsToDeploy = append(configsToDeploy, promptConfig(c))
			}
		}
	}

	// prompt for all secrets and provide existing value as default
	if prompt == "all" {
		for _, c := range config.Secrets {
			var existingValue string
			for _, a := range all {
				if c.Name == *a.Name {
					existingValue = *a.Value
					c.Value = *a.Value
				}
			}

			userInput := promptConfig(c)

			if userInput.Value != existingValue {
				configsToDeploy = append(configsToDeploy, userInput)
			}
		}
	}

	// filter configs with changed values
	for _, c := range config.Configs {
		for _, a := range all {
			if c.Name == *a.Name && c.Value != *a.Value {
				configsToDeploy = append(configsToDeploy, c)
			}
		}
	}

	err = st.PutMany(configsToDeploy)

	if err != nil {
		return errors.Wrap(err, "failed to write param")
	}

	fmt.Printf("%d new configs deployed", len(configsToDeploy))

	return nil
}

func promptConfig(config store.ConfigInput) store.ConfigInput {
	validate := func(input string) error {
		if len(input) < 1 {
			return fmt.Errorf("%s must not be empty", config.Name)
		}
		return nil
	}

	log.Printf("value %s", config.Value)

	prompt := promptui.Prompt{
		Label:    config.Key(),
		Validate: validate,
		Default:  config.Value,
	}

	result, _ := prompt.Run()

	config.Value = result

	return config
}

func getMissing(a []store.ConfigInput, b []store.Config) []store.ConfigInput {
	mb := make(map[string]struct{}, len(b))

	for _, x := range b {
		mb[*x.Name] = struct{}{}
	}

	var diff []store.ConfigInput
	for _, x := range a {
		if _, found := mb[x.Name]; !found {
			diff = append(diff, x)
		}
	}

	return diff
}
