package main

import (
	"github.com/riotkit-org/br-pg-simple-backup/assets"
	"github.com/riotkit-org/br-pg-simple-backup/cmd"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	tempDir := "/tmp/br-pgbr"
	if os.Getenv("BR_TEMP_DIR") != "" {
		tempDir = os.Getenv("BR_TEMP_DIR")
	}

	// prepare binaries
	if unpacked, err := assets.UnpackAll(tempDir); err != nil || !unpacked {
		if err == nil && !unpacked {
			log.Fatalf("Cannot unpack binaries and libraries to '%s'", tempDir)
		}

		log.Fatal(err)
	}
	if err := assets.PatchBinaries(tempDir); err != nil {
		log.Fatal(err)
	}

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
