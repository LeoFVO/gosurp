package smtp

import (
	"sync"

	"github.com/LeoFVO/gosurp/pkg/smtp"
	log "github.com/sirupsen/logrus"
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
		ports, _ := cmd.Flags().GetStringSlice("port")
		
		var wg sync.WaitGroup // Use a WaitGroup to block main() exit until all goroutines have completed.

		for _, port := range ports {
			log.Infof("Starting server on %s:%s", hostname, port)
			wg.Add(1)
			go func(port string) {
				err := smtp.Server{Hostname: hostname, Port: port}.Listen()
				if err != nil {
					log.Errorf("Error starting server on %s:%s: %s", hostname, port, err)
					wg.Done()
				}
			}(port)
		}
		wg.Wait()

		return nil
	},
}

func init() {
	listen.PersistentFlags().StringP("hostname", "", "localhost", "Hostname to listen on")
	listen.PersistentFlags().StringSliceP("port", "", []string{"25"}, "Ports to listen on, comma separated")
}