package inspect

import (
	"log"
	"strings"

	"github.com/LeoFVO/gosurp/pkg/inspect"
	"github.com/spf13/cobra"
)

var inspectDMARC = &cobra.Command{
	Use:   "dmarc",
	Aliases: []string{},
	Short: "Inspect domain DMARC records",
	Long:  `Inspect domain DMARC records and configuration`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := strings.TrimSpace(strings.ToLower(args[0]))
		domain, err := inspect.New(name)
		if err != nil {
			log.Fatal(err)
			return err
		}

		err = domain.InspectDMARC()
		if err != nil {
			log.Fatal(err)
			return err
		}

		return nil
	},
}

func init() {
}