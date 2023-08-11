package cmd

import (
	"fmt"

	"github.com/adikari/safebox/v2/store"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	getParam string

	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Gets parameter",
		RunE:  getE,
	}
)

func init() {
	getCmd.Flags().StringVarP(&getParam, "param", "p", "", "parameter to get")
	getCmd.MarkFlagRequired("param")

	rootCmd.AddCommand(getCmd)
}

func getE(_ *cobra.Command, _ []string) error {
	config, err := loadConfig()

	if err != nil {
		return errors.Wrap(err, "failed to load config")
	}

	st, err := store.GetStore(store.StoreConfig{
		Provider: config.Provider,
		Region:   config.Region,
		FilePath: config.Filepath,
	})

	if err != nil {
		return errors.Wrap(err, "failed to instantiate store")
	}

	found, err := st.Get(store.ConfigInput{Name: fmt.Sprintf("%s%s", config.Prefix, getParam)})

	if err != nil {
		return errors.Wrap(err, "failed to get param")
	}

	if found != nil {
		fmt.Printf("%s\n", *found.Value)
	}
	return nil
}
