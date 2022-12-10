package assets

import (
	"log"
	"os"
)

func UnpackOrExit() string {
	tempDir := "/tmp/br-pgbr"
	if os.Getenv("BR_TEMP_DIR") != "" {
		tempDir = os.Getenv("BR_TEMP_DIR")
	}

	// prepare binaries
	if unpacked, err := ExtractAllFromMemory(tempDir); err != nil || !unpacked {
		if err == nil && !unpacked {
			log.Fatalf("Cannot unpack binaries and libraries to '%s'", tempDir)
		}
		log.Fatal(err)
	}
	if err := PatchBinaries(tempDir); err != nil {
		log.Fatal(err)
	}
	return tempDir
}
