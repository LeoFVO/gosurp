package client

import (
	"fmt"
	"log"
	"net/mail"
	"os"

	"github.com/LeoFVO/gosurp/pkg/client"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server client.SMTPServer `yaml:"server"`
	Mail client.Mail `yaml:"mail"`
}

var sendMail = &cobra.Command{
	Use:   "send",
	Short: "Send mail from the cli",
	Long:  `Send mail from the cli`,
	Args:  cobra.MinimumNArgs(0),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		fromFile, _ := cmd.Flags().GetString("from-file")
		if fromFile != "" {
			var config Config
			file, err := os.Open(fromFile)
			if err != nil {
					return err
			}
			defer file.Close()

			decoder := yaml.NewDecoder(file)
			err = decoder.Decode(&config)
			if err != nil {
					return err
			}
			cmd.Flags().Set("from", config.Mail.From.Address)
			for _, recipient := range config.Mail.To {
				cmd.Flags().Set("to", recipient.Address)
			}
			cmd.Flags().Set("subject", config.Mail.Subject)
			cmd.Flags().Set("body", config.Mail.Body)
			cmd.Flags().Set("hostname", config.Server.Host)
			cmd.Flags().Set("port", config.Server.Port)
			cmd.Flags().Set("username", config.Server.Username)
			cmd.Flags().Set("password", config.Server.Password)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		fromEmail, _ := cmd.Flags().GetString("from")
		toEmail, _ := cmd.Flags().GetStringSlice("to")
		subject, _ := cmd.Flags().GetString("subject")
		body, _ := cmd.Flags().GetString("body")

		// Create a new SMTP server instance
		hostname, _ := cmd.Flags().GetString("hostname")
		port, _ := cmd.Flags().GetString("port")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		server := client.SMTPServer{
			Host: hostname,
			Port: port,
			Username: username,
			Password: password,
		}

		from := mail.Address{Name: "", Address: fromEmail}
		to := []mail.Address{}
		for _, recipient := range toEmail {
			to = append(to, mail.Address{Name: "", Address: recipient})
		}

		mail := client.Mail{From: from, To: to, Subject: subject, Body: body}
		log.Println("Sending email...")
		if err := server.SendMail(mail, false); err != nil {
			fmt.Println("Error sending email:", err)
		} else {
			fmt.Println("Email sent successfully!")
		}
		return nil
	},
}

func init() {
	sendMail.PersistentFlags().StringP("from", "", "sender@example.com", "Email address to send mail from")
	sendMail.PersistentFlags().StringSlice("to", []string{"receiver@victim.com"}, "Email address to send mail to (separated by comma)")
	sendMail.PersistentFlags().StringP("subject", "", "Test mail", "Subject to send mail with")
	sendMail.PersistentFlags().StringP("body", "", "Here is my simple email test. Did you received it ?", "Body to send mail with")
	sendMail.PersistentFlags().StringP("hostname", "", "localhost", "Hostname of the SMTP server to connect to")
	sendMail.PersistentFlags().StringP("port", "", "25", "Port of the SMTP server to connect to")
	sendMail.PersistentFlags().StringP("username", "", "", "Username to authenticate with")
	sendMail.PersistentFlags().StringP("password", "", "", "Password to authenticate with")
	sendMail.PersistentFlags().StringP("from-file", "f", "", "Config file to use")
}