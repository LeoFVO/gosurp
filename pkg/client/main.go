package client

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"strings"

	"github.com/LeoFVO/gosurp/pkg/dkim"
)

// SMTPServer represents an SMTP server.
type SMTPServer struct {
    Host string `yaml:"host"`
    Port string `yaml:"port"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
}

// Mail represents an email message.
type Mail struct {
    From    mail.Address `yaml:"from"`
    To      []mail.Address `yaml:"to"`
    Subject string `yaml:"subject"`
    Body    string `yaml:"body"`
}

// SendMail sends an email using the provided SMTP server and mail message.
func (s *SMTPServer) SendMail(mail Mail, withTLS bool) error {

    if s.Host == "" {
        return fmt.Errorf("SMTP server address required")
    }
    if s.Port == "" {
        return fmt.Errorf("SMTP server port required")
    }

    // Connect to the SMTP server.
    conn, err := net.Dial("tcp", s.Host+":"+s.Port)
    if err != nil {
        return err
    }
    defer conn.Close()

    if withTLS {
    // Convert the connection to a TLS connection.
        config := &tls.Config{
            ServerName: s.Host,
        }
        conn = tls.Client(conn, config)
    }
    

    // Set up an SMTP client.
    client, err := smtp.NewClient(conn, s.Host)
    if err != nil {
        return err
    }
    defer client.Quit()

    // Authenticate with the server.
    if s.Username != "" && s.Password != "" {
        if err := client.Auth(smtp.PlainAuth("", s.Username, s.Password, s.Host)); err != nil {
            return err
        }
    }

    // Set the sender and recipient of the email.
    if err := client.Mail(mail.From.Address); err != nil {
        return err
    }
    for _, to := range mail.To {
        if err := client.Rcpt(to.Address); err != nil {
            return err
        }
    }

    // Send the email body.
    wc, err := client.Data()
    if err != nil {
        return err
    }
    defer wc.Close()


    privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("error generating RSA key: %s", err) 
	}

	// Create signer
	signer := &dkim.Signer{
		PrivateKey: privateKey,
        /*
         * The domain name that the email is being sent from.
         * This is the domain that the DKIM public key will be published under in DNS.
         * The domain name should be the same as the domain name in the From header of the email.
         */
		Domain: strings.Split(mail.From.String(), "@")[1],
        /*
         * In DKIM, a selector is a string that identifies a specific public key in the DNS record for the domain. 
         * The DKIM signature for a message includes the selector value, allowing the recipient to look up the corresponding public key in DNS.
         * The "default" selector is just a convention that can be used if you don't need to have multiple selectors for the same domain. 
         * It's just a standard way to name the selector that holds the DKIM public key for the domain, and it makes it easier to configure DKIM for many email clients and services that expect the selector to be named "default".
         */
		Selector: "default",
	}


    // Write the email headers.
    headers := fmt.Sprintf("From: %s\r\n", mail.From.String())
    headers += fmt.Sprintf("To: %s\r\n", mail.To[0].String())
    for _, to := range mail.To[1:] {
        headers += fmt.Sprintf("Cc: %s\r\n", to.String())
    }
    headers += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
    header, err := signer.Sign(mail.Body)
	if err != nil {
		return fmt.Errorf("error signing message: %s", err)
	}
    headers += fmt.Sprintf("DKIM-Signature: %s\r\n", header)
    headers += "\r\n"

    if _, err := fmt.Fprintf(wc, headers+mail.Body); err != nil {
        return err
    }

    return nil
}
