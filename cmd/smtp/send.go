package smtp

import (
	"fmt"
	"net/mail"
	"os"

	smtp "github.com/LeoFVO/gosurp/pkg/smtp"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Server smtp.Server `yaml:"server"`
	Mail smtp.Envelope `yaml:"mail"`
	DKIM smtp.DKIM `yaml:"dkim"`
}

var send = &cobra.Command{
	Use:   "send",
	Short: "Send mail from local SMTP server",
	Long:  `Send mail from local SMTP server to the specified recipients`,
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
			
			cmd.Flags().Set("from", fmt.Sprintf("%s <%s>", config.Mail.From.Name, config.Mail.From.Address))
			for _, recipient := range config.Mail.To {
				cmd.Flags().Set("to", fmt.Sprintf("%s <%s>", recipient.Name, recipient.Address))
			}
			cmd.Flags().Set("subject", config.Mail.Subject)
			cmd.Flags().Set("body", config.Mail.Body)
			cmd.Flags().Set("hostname", config.Server.Hostname)
			cmd.Flags().Set("port", config.Server.Port)
			cmd.Flags().Set("username", config.Server.Username)
			cmd.Flags().Set("password", config.Server.Password)

			if config.DKIM.Selector != "" {
				cmd.Flags().Set("dkim-selector", config.DKIM.Selector)
			}
			if config.DKIM.PrivateKey != "" {
				cmd.Flags().Set("dkim-key", config.DKIM.PrivateKey)
			}
			if config.DKIM.Domain != "" {
				cmd.Flags().Set("dkim-domain", config.DKIM.Domain)
			}
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
		server := smtp.Server{
			Hostname: hostname,
			Port: port,
			Username: username,
			Password: password,
		}

		// Should we try to bypass DKIM 
		selector, _ := cmd.Flags().GetString("dkim-selector")
		privateKey, _ := cmd.Flags().GetString("dkim-key")
		domain, _ := cmd.Flags().GetString("dkim-domain")

		options := smtp.SendOptions{}
		if selector != "" && privateKey != "" {
			options.DKIM = smtp.DKIM{
				Selector: selector,
				PrivateKey: privateKey,
				Domain: domain,
			}
		}

		// Should we force TLS connection
		forceTLS, _ := cmd.Flags().GetBool("tls")
		options.WithTLS = forceTLS

		from, _ := mail.ParseAddress(fromEmail)
		to := []mail.Address{}
		for _, recipient := range toEmail {
			recipient, _ := mail.ParseAddress(recipient)
			to = append(to, *recipient)
		}

		mail := smtp.Envelope{From: *from, To: to, Subject: subject, Body: body}
		log.Println("Start sending email...")

		if err := server.Send(mail, options); err != nil {
			log.Errorf("Error while sending email: %v", err)
			return err
		} else {
			log.Println("Email sent successfully")
		}
		return nil
	},
}

func init() {
	// Mail flags
	send.PersistentFlags().StringP("from", "", "John Doe <john.doe@fsociety.com>", "Email address to send mail from")
	send.PersistentFlags().StringSlice("to", []string{"agent1337@evilcorp.com"}, "Email address to send mail to (separated by comma)")
	send.PersistentFlags().StringP("subject", "", "My phishing mail", "Subject to send mail with")
	send.PersistentFlags().StringP("body", "", "Here is my simple email test. Did you received it ?", "Body to send mail with")

	// SMTP server flags
	send.PersistentFlags().StringP("hostname", "", "localhost", "Hostname of the SMTP server to connect to")
	send.PersistentFlags().StringP("port", "", "25", "Port of the SMTP server to connect to")
	send.PersistentFlags().StringP("username", "", "", "Username to authenticate with")
	send.PersistentFlags().StringP("password", "", "", "Password to authenticate with")

	// Bypass DKIM flags
	send.PersistentFlags().StringP("dkim-selector", "", "", "DKIM selector to use")
	send.PersistentFlags().StringP("dkim-domain", "", "", "DKIM domain to use")
	send.PersistentFlags().StringP("dkim-key", "", "", "DKIM private key to use")

	// Force TLS flags
	send.PersistentFlags().BoolP("tls", "", false, "Force TLS connection")

	// Use can provide a config file to use instead of the flags
	send.PersistentFlags().StringP("from-file", "f", "", "Config file to use")
}