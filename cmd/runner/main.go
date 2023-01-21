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
	var c *exec.Cmd

	if os.Getenv("PGBR_USE_CONTAINER") != "" && os.Getenv("POSTGRES_VERSION") != "" {
		containerImage := "bitnami/postgresql"
		// allow to override the image
		if os.Getenv("PGBR_CONTAINER_IMAGE") != "" {
			containerImage = os.Getenv("PGBR_CONTAINER_IMAGE")
		}
		containerImage += ":" + os.Getenv("POSTGRES_VERSION")

		args := []string{"run"}

		// apply env variables as "-e" docker run commandline switches
		for _, envVar := range envVars {
			args = append(args, "-e", envVar)
		}

		args = append(args, "-i", "--entrypoint", binName, "--rm", containerImage)
		args = append(args, execArgs...)

		c = exec.Command("docker", args...)
	} else {
		c = exec.Command(binName, execArgs...)
		env := os.Environ()
		env = append(env, envVars...)
		c.Env = env
	}

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
