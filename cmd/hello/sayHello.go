package hello

import (
	"fmt"

	"github.com/spf13/cobra"
)

var sayHello = &cobra.Command{
	Use:   "to",
	Short: "Say hello to someone",
	Long:  `Say hello to someone`,
	Args:  cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")

		fmt.Printf("Hello %s!\n", name)
		return nil
	},
}

func init() {
	sayHello.PersistentFlags().StringP("name", "n", "world", "The name to say hello to")
}