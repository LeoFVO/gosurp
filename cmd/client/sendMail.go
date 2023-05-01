package client

import (
	"fmt"
	"log"
	"net/mail"

	"github.com/LeoFVO/gosurp/pkg/client"
	"github.com/spf13/cobra"
)

var sendMail = &cobra.Command{
	Use:   "send",
	Short: "Send mail from the cli",
	Long:  `Send mail from the cli`,
	Args:  cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {

		fromEmail, _ := cmd.Flags().GetString("from")
		toEmail, _ := cmd.Flags().GetString("to")
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
		to := []mail.Address{{Name: "", Address: toEmail}}

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
	sendMail.PersistentFlags().StringP("to", "", "receiver@example.com", "Email address to send mail to")
	sendMail.PersistentFlags().StringP("subject", "", "Test mail", "Subject to send mail with")
	sendMail.PersistentFlags().StringP("body", "", "Here is my simple email test. Did you received it ?", "Body to send mail with")
	sendMail.PersistentFlags().StringP("hostname", "", "localhost", "Hostname of the SMTP server to connect to")
	sendMail.PersistentFlags().StringP("port", "", "25", "Port of the SMTP server to connect to")
	sendMail.PersistentFlags().StringP("username", "", "", "Username to authenticate with")
	sendMail.PersistentFlags().StringP("password", "", "", "Password to authenticate with")
}