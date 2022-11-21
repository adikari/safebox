package cmd

import (
	"fmt"
	"log"

	"github.com/adikari/safebox/v2/cloudformation"
	conf "github.com/adikari/safebox/v2/config"
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

	if err != nil {
		return errors.Wrap(err, "failed to load config")
	}

	variables := map[string]string{
		"stage":   stage,
		"service": config.Service,
	}

	if len(config.Stacks) > 0 {
		cf := cloudformation.NewCloudformation()
		outputs, err := cf.GetOutput(config.Stacks[0])

		if err != nil {
			return errors.Wrap(err, "failed to load outputs")
		}

		for key, value := range outputs {
			variables[key] = value
		}
	}

	st, err := store.GetStore(config.Provider)

	if err != nil {
		return errors.Wrap(err, "failed to instantiate store")
	}

	all, err := st.GetMany(config.All)

	if err != nil {
		return errors.Wrap(err, "failed to read existing params")
	}

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
	for i, c := range config.Configs {
		co := config.Configs[i]
		v, err := conf.Interpolate(c.Value, variables)

		if err != nil {
			return errors.Wrap(err, "failed to interpolate template variables")
		}

		co.Value = v
		found := false
		for _, a := range all {
			if co.Name == *a.Name {
				found = true

				if co.Value != *a.Value {
					configsToDeploy = append(configsToDeploy, co)
				}
				break
			}
		}

		if !found {
			configsToDeploy = append(configsToDeploy, co)
		}
	}

	err = st.PutMany(configsToDeploy)

	if err != nil {
		return errors.Wrap(err, "failed to write params")
	}

	if removeOrphans {
		orphans, err := doRemoveOrphans(st, config.Prefix, config.All)
		if err != nil {
			log.Print("failed to remove orphans")
		}

		fmt.Printf("%d orphans removed.\n", len(orphans))
	}

	fmt.Printf("%d new configs deployed. service = %s, stage = %s\n", len(configsToDeploy), config.Service, stage)

	return nil
}

func doRemoveOrphans(store store.Store, prefix string, all []store.ConfigInput) ([]string, error) {
	var orphans []string
	params, err := store.GetByPath(prefix)

	if err != nil {
		return nil, err
	}

	for _, param := range params {
		exists := false

		for _, config := range all {
			if config.Name == *param.Name {
				exists = true
				break
			}
		}

		if !exists {
			orphans = append(orphans, *param.Name)
		}
	}

	return orphans, nil
}

func promptConfig(config store.ConfigInput) store.ConfigInput {
	validate := func(input string) error {
		if len(input) < 1 {
			return fmt.Errorf("%s must not be empty", config.Name)
		}
		return nil
	}

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
