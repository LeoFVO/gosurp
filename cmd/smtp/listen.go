package smtp

import (
	"github.com/LeoFVO/gosurp/pkg/smtp"
	"github.com/spf13/cobra"
)

var listen = &cobra.Command{
	Use:   "listen",
	Aliases: []string{"serve"},
	Short: "Start SMTP server and listen",
	Long:  `Start SMTP server and listen for incoming connections on the specified port`,
	Args:  cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		hostname, _ := cmd.Flags().GetString("hostname")
		port, _ := cmd.Flags().GetString("port")
		err := smtp.Server{Hostname: hostname, Port: port}.Listen()

		return err
	},
}

func init() {
	listen.PersistentFlags().StringP("hostname", "", "localhost", "Hostname to listen on")
	listen.PersistentFlags().StringP("port", "", "25", "Port to listen on")
}