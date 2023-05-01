package hello

import (
	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use:     "hello",
		Aliases: []string{"hey", "hi"},
		Short:   "Manage hellos from the CLI",
		Long:    `Manage hellos from the CLI`,
	}
)

func init() {
	RootCmd.AddCommand(sayHello)
}