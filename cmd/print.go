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

	if s.Config.Provider == "gpg" {
		fmt.Printf("%s\n", msg)
		return
	}

	if s.Config.Service != "" {
		msg += fmt.Sprintf(", service = %s", s.Config.Service)
	}

	if s.Config.Stage != "" {
		msg += fmt.Sprintf(", stage = %s", s.Config.Stage)
	}

	if s.Config.Region != "" {
		msg += fmt.Sprintf(", region = %s", s.Config.Region)
	}

	fmt.Printf("%s\n", msg)
}
