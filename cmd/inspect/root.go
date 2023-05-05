package inspect

import (
	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use:   "inspect",
		Aliases: []string{"get"},
		Short:   "Inspect DNS records",
		Long:    `Inspect DNS records and configuration publicly available`,
	}
)

func init() {
	RootCmd.AddCommand(inspectDomain)
	RootCmd.AddCommand(inspectSPF)
	RootCmd.AddCommand(inspectDMARC)
	RootCmd.AddCommand(inspectDKIM)
}