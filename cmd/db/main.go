package db

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// NewDbCommand creates the new command
func NewDbCommand(libDir string) *cobra.Command {
	command := &cobra.Command{
		Use:   "db",
		Short: "Operations on database level using dumps",
		Run: func(command *cobra.Command, args []string) {
			err := command.Help()
			if err != nil {
				logrus.Errorf(err.Error())
			}
		},
	}

	command.AddCommand(NewBackupCommand(libDir))
	command.AddCommand(NewRestoreCommand(libDir))

	return command
}
