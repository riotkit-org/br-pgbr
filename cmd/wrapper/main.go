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
		Use:   cmdName,
		Short: binName + " wrapper, pass " + binName + " parameters after '--'",
		Run: func(command *cobra.Command, args []string) {
			err := app.Run(binName, libDir, command.Flags().Args())

			if err != nil {
				logrus.Errorf(err.Error())
			}
		},
	}
	return command
}

type Wrapper struct{}

func (dw *Wrapper) Run(binName string, libDir string, execArgs []string) error {
	c := exec.Command(libDir+"/bin/"+binName, execArgs...)

	env := os.Environ()
	env = append(env, "LD_LIBRARY_PATH="+libDir)

	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	c.Env = env
	if waitErr := c.Run(); waitErr != nil {
		return errors.Wrap(waitErr, "error invoking "+binName)
	}

	return nil
}
