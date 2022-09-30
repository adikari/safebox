package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	getCmd = &cobra.Command{
		Use:     "get",
		Short:   "Gets given configuration from store",
		RunE:    get,
		Example: `TODO: get command example`,
	}
)

func init() {
	rootCmd.AddCommand(getCmd)
}

func get(cmd *cobra.Command, args []string) error {
	log.Panic("missing implementation")
	return nil
}
