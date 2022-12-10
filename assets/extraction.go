package assets

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

// ExtractAllFromMemory unpacks all stored libraries and binaries into single directory
func ExtractAllFromMemory(targetDir string) (bool, error) {
	var hasAtLeastOneError bool

	if err := os.RemoveAll(targetDir); err != nil {
		return false, errors.Wrapf(err, "cannot delete temporary directory at path: '%s'", targetDir)
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return false, errors.Wrapf(err, "cannot create directory '%s'", targetDir)
	}

	dir, readErr := Res.ReadDir(".build/data")
	if readErr != nil {
		return false, errors.Wrap(readErr, "cannot read directory from the go:embed")
	}

	for _, asset := range dir {
		assetName := asset.Name()
		if asset.IsDir() {
			continue
		}

		logrus.Debugf("Extracting '%s' into '%s'", assetName, targetDir)
		data, err := Res.ReadFile(".build/data/" + assetName)

		if err != nil {
			logrus.Error(errors.Wrap(err, "cannot unpack file from single-binary"))
			hasAtLeastOneError = true
		}

		subdir := ""
		baseName := path.Base(assetName)

		if baseName == "pg_dumpall" || baseName == "pg_restore" || baseName == "psql" || baseName == "pg_dump" {
			subdir = "bin/"
			_ = os.Mkdir(targetDir+"/"+subdir, 0755)
		}

		err = os.WriteFile(targetDir+"/"+subdir+baseName, data, 0755)
		if err != nil {
			logrus.Error(errors.Wrap(err, "cannot unpack file from single-binary"))
			hasAtLeastOneError = true
		}
	}

	return !hasAtLeastOneError, nil
}

func PatchBinaries(targetDir string) error {
	// todo: if patchelf is accessible, else return nil

	absDir, absErr := filepath.Abs(targetDir)
	if absErr != nil {
		return errors.Wrapf(absErr, "Cannot find absolute path for '%s'", targetDir)
	}

	interpreterPath := findInterpreterPath(absDir)
	logrus.Debugf("Interpreter path: '%s'", interpreterPath)

	if interpreterPath == "" {
		return errors.New("cannot find ld-linux or ld-musl")
	}

	binDir := absDir + "/bin"
	_ = os.MkdirAll(binDir, 0755)

	files, err := os.ReadDir(binDir)
	if err != nil {
		return errors.Wrapf(err, "Cannot list directory '%s' to patch binaries", targetDir)
	}

	for _, file := range files {
		logrus.Debugf("Processing '%s' with patchelf", binDir+"/"+file.Name())

		cmd := fmt.Sprintf("patchelf --set-interpreter %s %s", interpreterPath, binDir+"/"+file.Name())
		c := exec.Command("/bin/sh", "-c", cmd)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Stdin = os.Stdin
		c.Env = os.Environ()
		if waitErr := c.Run(); waitErr != nil {
			return errors.Wrapf(waitErr, "patchelf failed to set ld-linux path: %s", cmd)
		}
	}

	return nil
}

func findInterpreterPath(libPath string) string {
	possibleMatches := [][]string{
		// musl (Alpine Linux)
		glob(libPath + "/ld-musl-*"),
		// libc (all others)
		glob(libPath + "/ld-linux*"),
	}

	for _, matches := range possibleMatches {
		if len(matches) > 0 {
			return matches[0]
		}
	}

	return ""
}

func glob(pattern string) []string {
	matches, _ := filepath.Glob(pattern)
	return matches
}
