package db

import (
	"context"
	"fmt"
	pgx "github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/riotkit-org/br-pg-simple-backup/cmd/base"
	"github.com/riotkit-org/br-pg-simple-backup/cmd/wrapper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"
)

const TechDatabaseName = "br_empty_conn_db"

// NewRestoreCommand creates the new command
func NewRestoreCommand(libDir string) *cobra.Command {
	app := &RestoreCommand{}
	var basicOpts base.BasicOptions

	command := &cobra.Command{
		Use:          "restore",
		SilenceUsage: true,
		Short:        "Executes a backup restore procedure",
		RunE: func(command *cobra.Command, args []string) error {
			app.ExtraArgs = command.Flags().Args()
			base.PreCommandRun(command, &basicOpts)
			return app.Run(libDir)
		},
	}

	command.Flags().StringVarP(&app.Hostname, "host", "H", "127.0.0.1", "Hostname or IP address (default: 127.0.0.1)")
	command.Flags().IntVarP(&app.Port, "port", "p", 5432, "Port (default: 5432)")
	command.Flags().StringVarP(&app.Username, "user", "U", "postgres", "Username (default: postgres)")
	command.Flags().StringVarP(&app.Password, "password", "P", "", "Password")
	command.Flags().StringVarP(&app.DatabaseName, "db-name", "d", "", "Database name to restore. Leave empty if dump contains all databases. NOTICE: If dump contains more databases this switch does not allow to select them")
	command.Flags().StringVarP(&app.InitialDbName, "connection-database", "D", "postgres", "Any, even empty database name to connect to initially")
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
	DatabaseName  string
}

// Run Executes the command and outputs a stream to the stdout
func (bc *RestoreCommand) Run(libDir string) error {
	// 0) Prepare a database we will be connecting to, instead of connecting to target database
	if err := bc.createTechnicalDatabase(); err != nil {
		return err
	}
	bc.InitialDbName = TechDatabaseName

	client := bc.createClient()
	defer client.Close(context.Background())

	// 1) Kick off all clients
	if err := kickOffConnectedClients(client, bc.DatabaseName); err != nil {
		return errors.Wrap(err, "Cannot prepare to restore backup")
	}

	// 2) Set maintenance mode
	if err := setMaintenanceMode(client, true, bc.DatabaseName); err != nil {
		return errors.Wrap(err, "Cannot prepare to restore backup")
	}
	defer setMaintenanceMode(client, false, bc.DatabaseName)

	// 3) Restore structure & data
	logrus.Info("Restoring data...")
	envVars := []string{
		"PGPASSWORD=" + bc.Password,
	}
	if restoreErr := wrapper.RunWrappedPGCommand(libDir, "pg_restore", bc.buildRestoreArgs(), envVars); restoreErr != nil {
		return errors.Wrap(restoreErr, "Cannot restore backup")
	}

	logrus.Info("Database restored.")
	return nil
}

func (bc *RestoreCommand) buildRestoreArgs() []string {
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

	// passing extra arguments to the pg_dump/pg_dumpall
	if len(bc.ExtraArgs) > 0 {
		restoreArgs = append(restoreArgs, bc.ExtraArgs...)
	}

	return restoreArgs
}

// createTechnicalDatabase Creates an additional database to which we will be connecting, which will be excluded from the backup & restore process
func (bc *RestoreCommand) createTechnicalDatabase() error {
	client := bc.createClient()

	// check if database exists first
	var exists int
	row := client.QueryRow(context.TODO(), "SELECT 1 FROM pg_database WHERE datname = 'br_empty_conn_db'")
	if err := row.Scan(&exists); err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		return errors.Wrap(err, "Cannot parse result of database check query")
	}
	if exists == 1 {
		return nil
	}
	// create a database
	_, err := client.Exec(context.TODO(), "CREATE DATABASE br_empty_conn_db;")
	if err != nil {
		return errors.Wrap(err, "Cannot create database that would be used technically as connection database")
	}
	return nil
}

func (bc *RestoreCommand) createClient() *pgx.Conn {
	conn, _ := pgx.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%v/%s", bc.Username, bc.Password, bc.Hostname, bc.Port, bc.InitialDbName))
	return conn
}
