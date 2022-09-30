package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	cleanCmd = &cobra.Command{
		Use:     "clean",
		Short:   "Cleans all orphan configurations",
		RunE:    clean,
		Example: `TODO: clean command example`,
	}
)

func init() {
	rootCmd.AddCommand(cleanCmd)
}

func clean(cmd *cobra.Command, args []string) error {
	log.Panic("missing implementation")
	return nil
}
