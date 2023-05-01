package server

import (
	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use:     "server",
		Aliases: []string{"s"},
		Short:   "Start SMTP server and listen",
		Long:    `Start SMTP server and listen`,
	}
)

func init() {
	RootCmd.AddCommand(startServer)
}