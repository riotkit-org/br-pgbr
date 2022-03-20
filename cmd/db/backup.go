package db

import (
	"fmt"
	"github.com/riotkit-org/br-pg-simple-backup/cmd/base"
	"github.com/riotkit-org/br-pg-simple-backup/cmd/wrapper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// NewBackupCommand creates the new command
func NewBackupCommand(libDir string) *cobra.Command {
	app := &BackupCommand{}
	var basicOpts base.BasicOptions

	command := &cobra.Command{
		Use:   "backup",
		Short: "Backup using pg_dump and pg_dumpall",
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
	command.Flags().StringVarP(&app.Database, "db-name", "d", "", "Database name. Leave empty to dump all")
	command.Flags().IntVarP(&app.CompressionLevel, "compression-level", "Z", 7, "Compression level 0-9")
	base.PopulateFlags(command, &basicOpts)

	return command
}

type BackupCommand struct {
	Hostname         string
	Port             int
	Username         string
	Password         string
	Database         string
	CompressionLevel int
	ExtraArgs        []string
}

// Run Executes the command and outputs a stream to the stdout
func (bc *BackupCommand) Run(libDir string) error {
	dumpArgs := []string{
		"--clean",
		"--host", bc.Hostname,
		fmt.Sprintf("--port=%v", bc.Port),
		"--username", bc.Username,
		"--format=c", // custom format for pg_restore
		fmt.Sprintf("--compress=%v", bc.CompressionLevel),
	}
	envVars := []string{
		"PGPASSWORD=" + bc.Password,
	}

	// difference between pg_dump and pg_dumpall
	if !bc.allDatabases() {
		// pg_dump
		dumpArgs = append(dumpArgs, "--create", "--blobs", bc.Database)
	} else {
		// pg_dumpall
		dumpArgs = append(dumpArgs, "--superuser="+bc.Username)
	}

	// passing extra arguments to the pg_dump/pg_dumpall
	if len(bc.ExtraArgs) > 0 {
		dumpArgs = append(dumpArgs, bc.ExtraArgs...)
	}

	return wrapper.RunWrappedPGCommand(libDir, "pg_dump", dumpArgs, envVars)
}

func (bc *BackupCommand) allDatabases() bool {
	return bc.Database == ""
}
