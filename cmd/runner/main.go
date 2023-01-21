package runner

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

func Run(binName string, execArgs []string, envVars []string, captureOutput bool, buffer bytes.Buffer) error {
	logrus.Debugf("Running '%s' %v", binName, execArgs)

	c := exec.Command(binName, execArgs...)

	env := os.Environ()
	env = append(env, envVars...)

	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	// optionally capture output into a buffer
	if captureOutput {
		c.Stdout = &buffer
		c.Stderr = &buffer
	}
	c.Stdin = os.Stdin
	c.Env = env
	if waitErr := c.Run(); waitErr != nil {
		return errors.Wrapf(waitErr, "error invoking '%s'", binName)
	}
	return nil
}
