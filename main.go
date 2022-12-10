package main

import (
	"github.com/riotkit-org/br-pg-simple-backup/assets"
	"github.com/riotkit-org/br-pg-simple-backup/cmd"
	"os"
)

func main() {
	tempDir := assets.UnpackOrExit()
	command := cmd.Main(tempDir)
	args := os.Args

	if args != nil {
		args = args[1:]
		command.SetArgs(args)
	}

	err := command.Execute()
	if err != nil {
		os.Exit(1)
	}
}
