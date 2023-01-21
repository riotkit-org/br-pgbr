package cmd

import (
	"github.com/riotkit-org/br-pg-simple-backup/cmd/db"
	"github.com/spf13/cobra"
)

// Main creates the new command
func Main() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pgbr",
		Short: "PostgreSQL backup & restore runner for Backup Repository (and not only)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(db.NewDbCommand())

	return cmd
}
