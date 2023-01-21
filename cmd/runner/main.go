package runner

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

func Run(binName string, execArgs []string, envVars []string, captureOutput bool, inputStdin *bytes.Buffer) ([]byte, error) {
	logrus.Debugf("Running '%s' %v", binName, execArgs)

	c := exec.Command(binName, execArgs...)

	env := os.Environ()
	env = append(env, envVars...)
	c.Env = env

	// allow to optionally capture output into a buffer
	if captureOutput {
		c.Stdin = inputStdin
		out, waitErr := c.CombinedOutput()
		if waitErr != nil {
			return out, errors.Wrapf(waitErr, "error invoking '%s'", binName)
		}
		return out, nil
	}

	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin

	if waitErr := c.Run(); waitErr != nil {
		return []byte{}, errors.Wrapf(waitErr, "error invoking '%s'", binName)
	}
	return []byte{}, nil
}
