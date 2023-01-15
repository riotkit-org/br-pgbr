package db

import (
	"bytes"
	"github.com/spf13/cobra"
)

// NewDbCommand creates the new command
func NewDbCommand(libDir string) *cobra.Command {
	command := &cobra.Command{
		Use:   "db",
		Short: "Operations on database level using dumps",
		RunE: func(command *cobra.Command, args []string) error {
			return command.Help()
		},
	}

	command.AddCommand(NewBackupCommand(libDir, false, bytes.Buffer{}))
	command.AddCommand(NewRestoreCommand(libDir, false, bytes.Buffer{}))

	return command
}
