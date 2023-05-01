package server

import (
	"github.com/LeoFVO/gosurp/pkg/server"
	"github.com/spf13/cobra"
)

var startServer = &cobra.Command{
	Use:   "start",
	Short: "Start SMTP server and listen",
	Long:  `Start SMTP server and listen for incoming connections on the specified port`,
	Args:  cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		hostname, _ := cmd.Flags().GetString("hostname")
		port, _ := cmd.Flags().GetString("port")
		err := server.Start(server.Server{Hostname: hostname, Port: port})

		return err
	},
}

func init() {
	startServer.PersistentFlags().StringP("hostname", "", "localhost", "Hostname to listen on")
	startServer.PersistentFlags().StringP("port", "", "25", "Port to listen on")
}