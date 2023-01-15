package cmd

import (
	"bytes"
	"github.com/riotkit-org/br-pg-simple-backup/cmd/db"
	"github.com/riotkit-org/br-pg-simple-backup/cmd/wrapper"
	"github.com/spf13/cobra"
)

// Main creates the new command
func Main(libDir string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pgbr",
		Short: "PostgreSQL backup & restore wrapper for Backup Repository. Works also as a standalone, single-binary backup make & restore utility",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(wrapper.NewCmdPostgresWrapper(libDir, "psql", "psql", false, bytes.Buffer{}))
	cmd.AddCommand(wrapper.NewCmdPostgresWrapper(libDir, "pg_dump", "pg_dump", false, bytes.Buffer{}))
	cmd.AddCommand(wrapper.NewCmdPostgresWrapper(libDir, "pg_dumpall", "pg_dumpall", false, bytes.Buffer{}))
	cmd.AddCommand(wrapper.NewCmdPostgresWrapper(libDir, "pg_restore", "pg_restore", false, bytes.Buffer{}))
	cmd.AddCommand(db.NewDbCommand(libDir))

	return cmd
}
