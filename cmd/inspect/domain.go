package inspect

import (
	"log"
	"strings"

	"github.com/LeoFVO/gosurp/pkg/inspect"
	"github.com/spf13/cobra"
)

var inspectDomain = &cobra.Command{
	Use:   "domain <yourdomain.com>",
	Aliases: []string{"dns"},
	Short: "Inspect all domain public records",
	Long:  `Inspectall domain public records and configuration`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := strings.TrimSpace(strings.ToLower(args[0]))
		domain, err := inspect.New(name)
		if err != nil {
			log.Fatal(err)
			return err
		}

		err = domain.Inspect()
		if err != nil {
			log.Fatal(err)
			return err
		}

		return nil
	},
}

func init() {
}