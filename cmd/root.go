package cmd

import (
	"github.com/riotkit-org/br-pg-simple-backup/cmd/wrapper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Main creates the new command
func Main(libDir string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pgbr",
		Short: "PostgreSQL backup & restore wrapper for Backup Repository. Works also as a standalone, single-binary backup make & restore utility",
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				logrus.Errorf(err.Error())
			}
		},
	}
	cmd.AddCommand(wrapper.NewCmdPostgresWrapper(libDir, "psql", "psql"))
	cmd.AddCommand(wrapper.NewCmdPostgresWrapper(libDir, "pg_dump", "pg_dump"))
	cmd.AddCommand(wrapper.NewCmdPostgresWrapper(libDir, "pg_dumpall", "pg_dumpall"))
	cmd.AddCommand(wrapper.NewCmdPostgresWrapper(libDir, "pg_restore", "pg_restore"))
	return cmd
}
