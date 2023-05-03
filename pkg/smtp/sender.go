package smtp

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"

	"github.com/LeoFVO/gosurp/pkg/dkim"
	"github.com/LeoFVO/gosurp/pkg/utils"
	log "github.com/sirupsen/logrus"
)

type SendOptions struct {
    WithTLS bool
    DKIM DKIM
}

// SendMail sends an email using the provided SMTP server and mail message.
func (s *Server) Send(mail Envelope, options SendOptions) error {

    if s.Hostname == "" {
        return fmt.Errorf("SMTP server address required")
    }
    if s.Port == "" {
        return fmt.Errorf("SMTP server port required")
    }

    s.Hostname = utils.GetSMTPServerAddress(s.Hostname)

    // Connect to the SMTP server.
    conn, err := net.Dial("tcp", s.Hostname+":"+s.Port)
    if err != nil {
        return err
    }
    defer conn.Close()

    if options.WithTLS {
    // Convert the connection to a TLS connection.
        config := &tls.Config{
            ServerName: s.Hostname,
        }
        conn = tls.Client(conn, config)
    }
    

    // Set up an SMTP client.
    client, err := smtp.NewClient(conn, s.Hostname)
    if err != nil {
        return err
    }
    defer client.Quit()

    // Authenticate with the server.
    if s.Username != "" && s.Password != "" {
        if err := client.Auth(smtp.PlainAuth("", s.Username, s.Password, s.Hostname)); err != nil {
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

    signer := &dkim.Signer{
        /*
         * The domain name that the email is being sent from.
         * This is the domain that the DKIM public key will be published under in DNS.
         * The domain name should be the same as the domain name in the From header of the email.
         */
		Domain: strings.TrimSuffix(strings.Split(mail.From.String(), "@")[1], ">"),
        /*
         * In DKIM, a selector is a string that identifies a specific public key in the DNS record for the domain. 
         * The DKIM signature for a message includes the selector value, allowing the recipient to look up the corresponding public key in DNS.
         * The "default" selector is just a convention that can be used if you don't need to have multiple selectors for the same domain. 
         * It's just a standard way to name the selector that holds the DKIM public key for the domain, and it makes it easier to configure DKIM for many email clients and services that expect the selector to be named "default".
         */
		Selector: "default",
        /*
         * The headers that should be included in the DKIM signature.
         * The From, To, Subject, and Date headers are always included in the signature, so you don't need to include them here.
         * The headers should be in the order that they appear in the email.    
         */
        Headers: []string{
            "From",
            "To",
            "Subject",
        },
	}

    // if bypass DKIM is enabled, we need to sign the email using the provided private key and change selector

    if options.DKIM != (DKIM{}) {
        log.Tracef("Custom DKIM config provided, using selector %s and private key %s", options.DKIM.Selector, options.DKIM.PrivateKey)
        signer.Selector = options.DKIM.Selector

        // Load the private key from the provided path.
        privateKey := utils.LoadRSAPrivateKey(options.DKIM.PrivateKey)
        signer.PrivateKey = privateKey

        if options.DKIM.Domain != "" {
            log.Tracef("Custom DKIM config provided, using domain %s", options.DKIM.Domain)
            signer.Domain = options.DKIM.Domain
        }
    } else {
        log.Tracef("Custom DKIM config disabled, generating new RSA key pair")
        // Generate a new RSA key pair to use for the DKIM signature.
        privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
        if err != nil {
            return fmt.Errorf("error generating RSA key: %s", err) 
        }
        signer.PrivateKey = privateKey
    }


    // Write the email headers.
    log.Tracef("Starting to write email headers")
    headers := fmt.Sprintf("From: %s\r\n", mail.From.String())
    headers += fmt.Sprintf("To: %s\r\n", mail.To[0].String())
    for _, to := range mail.To[1:] {
        headers += fmt.Sprintf("Cc: %s\r\n", to.String())
    }
    headers += fmt.Sprintf("Subject: %s\r\n", mail.Subject)

    // Sign the email body and add the DKIM-Signature header.
    log.Tracef("Starting to sign email body")
    header, err := signer.Sign(headers+mail.Body)
	if err != nil {
		return fmt.Errorf("error signing message: %s", err)
	}
    log.Tracef("Starting to add DKIM-Signature header")
    headers += fmt.Sprintf("DKIM-Signature: %s\r\n", header)
    headers += "\r\n"

    // Write the email headers and body to the SMTP client.
    log.Tracef("Starting to write email headers and body to the SMTP client")
    if _, err := fmt.Fprintf(wc, headers+mail.Body); err != nil {
        return err
    }

    return nil
}
