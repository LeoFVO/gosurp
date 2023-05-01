package client

import (
	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use:     "client",
		Aliases: []string{"c"},
		Short:   "Send mail from the CLI",
		Long:    `Send mail from the CLI`,
	}
)

func init() {
	RootCmd.AddCommand(sendMail)
}