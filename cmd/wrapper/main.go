package wrapper

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"os/user"
)

// NewCmdPostgresWrapper creates the new command
func NewCmdPostgresWrapper(libDir string, binName string, cmdName string, captureOutput bool, buffer bytes.Buffer) *cobra.Command {
	app := &Wrapper{
		CaptureOutput: captureOutput,
		Buffer:        buffer,
	}

	command := &cobra.Command{
		Use:          cmdName,
		SilenceUsage: true,
		Short:        binName + " wrapper, pass " + binName + " parameters after '--'",
		RunE: func(command *cobra.Command, args []string) error {
			return app.Run(libDir, binName, command.Flags().Args())
		},
	}
	return command
}

type Wrapper struct {
	CaptureOutput bool
	Buffer        bytes.Buffer
}

func (dw *Wrapper) Run(libDir string, binName string, execArgs []string) error {
	return RunWrappedPGCommand(libDir, binName, execArgs, []string{}, dw.CaptureOutput, dw.Buffer)
}

func RunWrappedPGCommand(libDir string, binName string, execArgs []string, envVars []string, captureOutput bool, buffer bytes.Buffer) error {
	logrus.Debugf("Running '%s' %v", binName, execArgs)

	fullPath := libDir + "/bin/" + binName
	c := exec.Command(fullPath, execArgs...)

	// populate fake "/etc/passwd"
	u, _ := user.Current()
	etcPasswdPath := libDir + "/etc-passwd"
	etcPasswdContent := []byte(fmt.Sprintf("%s:x:%s:%s::/home/user:/bin/bash\n", u.Username, u.Uid, u.Gid))
	err := os.WriteFile(etcPasswdPath, etcPasswdContent, 0644)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("cannot write to '%s'", etcPasswdPath))
	}
	etcGroupPath := libDir + "/etc-group"
	etcGroupContent := []byte(fmt.Sprintf("%s:x:%s:\n", u.Username, u.Gid))
	err = os.WriteFile(etcGroupPath, etcGroupContent, 0644)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("cannot write to '%s'", etcGroupPath))
	}

	env := os.Environ()
	ldEnv := make([]string, 0)
	ldEnv = append(ldEnv, "LD_LIBRARY_PATH="+libDir)
	ldEnv = append(ldEnv, "LD_PRELOAD="+libDir+"/libnss_wrapper.so")
	ldEnv = append(ldEnv, "NSS_WRAPPER_PASSWD="+etcPasswdPath)
	ldEnv = append(ldEnv, "NSS_WRAPPER_GROUP="+etcGroupPath)

	env = append(env, ldEnv...)
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
		return errors.Wrapf(waitErr, "error invoking '%s' (%s), ldEnv: %v", binName, fullPath, ldEnv)
	}
	return nil
}
