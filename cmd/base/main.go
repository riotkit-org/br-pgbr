package base

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func PopulateFlags(command *cobra.Command, o *BasicOptions) {
	command.Flags().StringVarP(&o.LogLevel, "log-level", "", "info", "Logging level: error, warning, info, debug (default: info)")
}

func PreCommandRun(command *cobra.Command, o *BasicOptions) {
	level, _ := logrus.ParseLevel(o.LogLevel)
	logrus.SetLevel(level)
}

type BasicOptions struct {
	LogLevel string
}
