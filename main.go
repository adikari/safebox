package main

import "github.com/adikari/safebox/v2/cmd"

var (
	Version = "dev"
)

func main() {
	cmd.Execute(Version)
}
