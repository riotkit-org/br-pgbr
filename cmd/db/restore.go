package db

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/riotkit-org/br-pg-simple-backup/cmd/base"
	"github.com/riotkit-org/br-pg-simple-backup/cmd/wrapper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// NewRestoreCommand creates the new command
func NewRestoreCommand(libDir string) *cobra.Command {
	app := &RestoreCommand{}
	var basicOpts base.BasicOptions

	command := &cobra.Command{
		Use:   "restore",
		Short: "Executes a backup restore procedure",
		Run: func(command *cobra.Command, args []string) {
			app.ExtraArgs = command.Flags().Args()
			base.PreCommandRun(command, &basicOpts)
			err := app.Run(libDir)

			if err != nil {
				logrus.Errorf(err.Error())
			}
		},
	}

	command.Flags().StringVarP(&app.Hostname, "host", "H", "127.0.0.1", "Hostname or IP address (default: 127.0.0.1)")
	command.Flags().IntVarP(&app.Port, "port", "p", 5432, "Port (default: 5432)")
	command.Flags().StringVarP(&app.Username, "user", "U", "postgres", "Username (default: postgres)")
	command.Flags().StringVarP(&app.Password, "password", "P", "", "Password")
	command.Flags().StringVarP(&app.InitialDbName, "connection-database", "d", "postgres", "Database to connect to")
	base.PopulateFlags(command, &basicOpts)

	return command
}

type RestoreCommand struct {
	InitialDbName string
	Hostname      string
	Port          int
	Username      string
	Password      string
	ExtraArgs     []string
}

// Run Executes the command and outputs a stream to the stdout
func (bc *RestoreCommand) Run(libDir string) error {
	restoreArgs := []string{
		"--clean",
		"--create",
		"--exit-on-error",

		"--host", bc.Hostname,
		"--port", fmt.Sprintf("%v", bc.Port),
		"--username", bc.Username,
		"--format=c",
		"--dbname=" + bc.InitialDbName,
	}
	envVars := []string{
		"PGPASSWORD=" + bc.Password,
	}

	// passing extra arguments to the pg_dump/pg_dumpall
	if len(bc.ExtraArgs) > 0 {
		restoreArgs = append(restoreArgs, bc.ExtraArgs...)
	}

	if restoreErr := wrapper.RunWrappedPGCommand(libDir, "pg_restore", restoreArgs, envVars); restoreErr != nil {
		return errors.Wrap(restoreErr, "Cannot restore backup")
	}

	return nil
}
