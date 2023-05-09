package smtp

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"

	"github.com/LeoFVO/gosurp/pkg/dkim"
	"github.com/LeoFVO/gosurp/pkg/utils"
	log "github.com/sirupsen/logrus"
)

type SendOptions struct {
    WithTLS bool
    DKIM DKIM
    Timeout time.Duration
    Port string
}

// SendMail sends an email using the provided SMTP server and mail message.
func (s Server) Send(mail Envelope, options SendOptions) error {
    
    var usurpedServer Server

    // Define the hostname of the smtp server we are usurping.
    if s.Hostname != "" {
        usurpedServer = s
        log.Debugf("Using custom usurped server hostname: %s", usurpedServer.Hostname)
    } else {
        usurpedServer.Hostname = strings.Split(mail.From.Address, "@")[1]
        log.Debugf("No custom usurped server hostname provided, using the domain from the \"From\" header : %s", usurpedServer.Hostname)
    }
    
    // Get the target SMTP server address.
    targetServer := utils.GetSMTPServerAddress(strings.Split(mail.To[0].Address, "@")[1])
    log.Debugf("Trying to connect to SMTP server %s:%s", targetServer, options.Port)

    // Connect to the SMTP server.
    session := net.Dialer{Timeout: options.Timeout}
    // Connect to the target SMTP server
    conn, err := session.Dial("tcp", fmt.Sprintf("%s:%s", targetServer, options.Port))
    if err != nil {
        return err
    }
    defer conn.Close()
    log.Tracef("SMTP server connection established to %s:%s", targetServer, options.Port)

    // Set up an SMTP client.
    client, err := smtp.NewClient(conn, targetServer)
    if err != nil {
        log.Errorf("Error creating SMTP client: %s", err)
        return err
    }
    err = client.Hello(usurpedServer.Hostname)
    if err != nil {
        log.Errorf("Error sending HELO: %s", err)
        return err
    }

    log.Infof("SMTP client successfully created and connected from %s to %s", conn.LocalAddr().String(), targetServer)
    defer client.Quit()
 
    if options.WithTLS {
       log.Debugf("SMTP server requires TLS, converting connection to TLS")
        // Convert the connection to a TLS connection.
       config := &tls.Config{
           ServerName: usurpedServer.Hostname,
       }

       //ask for STARTTLS
       if ok, _ := client.Extension("STARTTLS"); ok {
           if err := client.StartTLS(config); err != nil {
               return err
           }
       } else {
           return fmt.Errorf("SMTP server does not support STARTTLS")
       }
       
       log.Debugf("Connection converted to TLS")
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

    // Building the mail envelope
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
		Domain: strings.Split(mail.From.Address, "@")[1],
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
