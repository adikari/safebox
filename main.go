package main

import (
	"github.com/adikari/safebox/v2/cmd"
)

var (
	version = "dev"
)

func main() {
	cmd.Execute(version)
}
