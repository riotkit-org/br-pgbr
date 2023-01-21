package db

import (
	"bytes"
	"fmt"
	"github.com/riotkit-org/br-pg-simple-backup/cmd/base"
	"github.com/riotkit-org/br-pg-simple-backup/cmd/runner"
	"github.com/spf13/cobra"
)

// NewBackupCommand creates the new command
func NewBackupCommand(captureOutput bool, buffer bytes.Buffer) *cobra.Command {
	app := &BackupCommand{
		CaptureOutput: captureOutput,
		Buffer:        buffer,
	}
	var basicOpts base.BasicOptions

	command := &cobra.Command{
		Use:          "backup",
		SilenceUsage: true,
		Short:        "Backup using pg_dump and pg_dumpall",
		RunE: func(command *cobra.Command, args []string) error {
			app.ExtraArgs = command.Flags().Args()
			base.PreCommandRun(command, &basicOpts)
			return app.Run()
		},
	}

	command.Flags().StringVarP(&app.Hostname, "host", "H", "127.0.0.1", "Hostname or IP address (default: 127.0.0.1)")
	command.Flags().IntVarP(&app.Port, "port", "p", 5432, "Port (default: 5432)")
	command.Flags().StringVarP(&app.Username, "user", "U", "postgres", "Username (default: postgres)")
	command.Flags().StringVarP(&app.Password, "password", "P", "", "Password")
	command.Flags().StringVarP(&app.Database, "db-name", "d", "", "Database name. Leave empty to dump all")
	command.Flags().IntVarP(&app.CompressionLevel, "compression-level", "Z", 7, "Compression level 0-9")
	command.Flags().StringVarP(&app.InitialDbName, "connection-database", "D", "postgres", "Any, even empty database name to connect to initially")
	base.PopulateFlags(command, &basicOpts)

	return command
}

type BackupCommand struct {
	Hostname         string
	Port             int
	Username         string
	Password         string
	Database         string
	InitialDbName    string
	CompressionLevel int
	ExtraArgs        []string
	CaptureOutput    bool
	Buffer           bytes.Buffer
}

// Run Executes the command and outputs a stream to the stdout
func (bc *BackupCommand) Run() error {
	dumpArgs := []string{
		"--clean",
		"--host", bc.Hostname,
		fmt.Sprintf("--port=%v", bc.Port),
		"--username", bc.Username,
		fmt.Sprintf("--compress=%v", bc.CompressionLevel),
	}
	envVars := []string{
		"PGPASSWORD=" + bc.Password,
	}

	var binName string

	// difference between pg_dump and pg_dumpall
	if !bc.allDatabases() {
		// pg_dump
		dumpArgs = append(dumpArgs,
			"--create",
			"--blobs", bc.Database,
			"--format=c", // custom format for pg_restore
		)
		binName = "pg_dump"
	} else {
		// pg_dumpall
		dumpArgs = append(dumpArgs, "--superuser="+bc.Username, "--dbname="+bc.InitialDbName)
		binName = "pg_dumpall"
	}

	// passing extra arguments to the pg_dump/pg_dumpall
	if len(bc.ExtraArgs) > 0 {
		dumpArgs = append(dumpArgs, bc.ExtraArgs...)
	}

	return runner.Run(binName, dumpArgs, envVars, bc.CaptureOutput, bc.Buffer)
}

func (bc *BackupCommand) allDatabases() bool {
	return bc.Database == ""
}
