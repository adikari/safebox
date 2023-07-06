package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	importFormat string
	inputFile    string

	importCmd = &cobra.Command{
		Use:     "import",
		Short:   "Imports all configuration from a file",
		RunE:    importE,
		Example: `TODO: import command example`,
	}
)

func init() {
	importCmd.Flags().StringVarP(&importFormat, "format", "f", "json", "input format (json, yaml, dotenv)")
	importCmd.Flags().StringVarP(&inputFile, "input-file", "i", "", "input file")
	importCmd.MarkFlagRequired("input-file")
	importCmd.MarkFlagFilename("input-file")

	rootCmd.AddCommand(importCmd)
}

func importE(_ *cobra.Command, _ []string) error {
	log.Fatalf("not implemented")
	return nil
}
