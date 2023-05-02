package smtp

import (
	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use:     "smtp",
		Aliases: []string{},
		Short:   "Send mail from the CLI using SMTP",
		Long:    `Send mail from the CLI using SMTP`,
	}
)

func init() {
	RootCmd.AddCommand(send)
	RootCmd.AddCommand(listen)
}