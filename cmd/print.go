package cmd

import (
	"fmt"

	c "github.com/adikari/safebox/v2/config"
)

type Summary struct {
	Message string
	Config  c.Config
}

func PrintSummary(s Summary) {
	msg := ""

	if s.Message != "" {
		msg += fmt.Sprintf("%s", s.Message)
	}

	if s.Config.Service != "" && s.Config.Provider != "gpg" {
		msg += fmt.Sprintf(", service = %s", s.Config.Service)
	}

	if s.Config.Stage != "" {
		msg += fmt.Sprintf(", stage = %s", s.Config.Stage)
	}

	if s.Config.Region != "" && s.Config.Provider != "gpg" {
		msg += fmt.Sprintf(", region = %s", s.Config.Region)
	}

	if s.Config.Provider == "gpg" {
		msg += fmt.Sprintf(", file = %s", s.Config.Filepath)
	}

	fmt.Printf("%s\n", msg)
}
