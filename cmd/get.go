package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	param string

	getCmd = &cobra.Command{
		Use:   "get",
		Short: "get parameter",
		RunE:  getE,
	}
)

func init() {
	getCmd.Flags().StringVarP(&param, "param", "p", "", "input format (json, yaml, dotenv)")
	getCmd.MarkFlagRequired("param")

	rootCmd.AddCommand(getCmd)
}

func getE(_ *cobra.Command, _ []string) error {
	log.Fatalf("not implemented")
	return nil
}
