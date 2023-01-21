package db

import (
	"bytes"
	"github.com/spf13/cobra"
)

// NewDbCommand creates the new command
func NewDbCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "db",
		Short: "Operations on database level using dumps",
		RunE: func(command *cobra.Command, args []string) error {
			return command.Help()
		},
	}

	backupCmd, _ := NewBackupCommand(false, &bytes.Buffer{})
	command.AddCommand(backupCmd)

	restoreCmd, _ := NewRestoreCommand(false, &bytes.Buffer{})
	command.AddCommand(restoreCmd)

	return command
}
