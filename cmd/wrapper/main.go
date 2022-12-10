package wrapper

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

// NewCmdPostgresWrapper creates the new command
func NewCmdPostgresWrapper(libDir string, binName string, cmdName string) *cobra.Command {
	app := &Wrapper{}

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

type Wrapper struct{}

func (dw *Wrapper) Run(libDir string, binName string, execArgs []string) error {
	return RunWrappedPGCommand(libDir, binName, execArgs, []string{})
}

func RunWrappedPGCommand(libDir string, binName string, execArgs []string, envVars []string) error {
	logrus.Debugf("Running '%s' %v", binName, execArgs)

	fullPath := libDir + "/bin/" + binName
	c := exec.Command(fullPath, execArgs...)

	env := os.Environ()
	env = append(env, "LD_LIBRARY_PATH="+libDir)
	env = append(env, envVars...)

	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	c.Env = env
	if waitErr := c.Run(); waitErr != nil {
		return errors.Wrapf(waitErr, "error invoking '%s' (%s), LD_LIBRARY_PATH=%s", binName, fullPath, libDir)
	}
	return nil
}
