package client

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
)

// SMTPServer represents an SMTP server.
type SMTPServer struct {
    Host string
    Port string
    Username string
    Password string
}

// Mail represents an email message.
type Mail struct {
    From    mail.Address
    To      []mail.Address
    Subject string
    Body    string
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

    // Write the email headers.
    headers := fmt.Sprintf("From: %s\r\n", mail.From.String())
    headers += fmt.Sprintf("To: %s\r\n", mail.To[0].String())
    for _, to := range mail.To[1:] {
        headers += fmt.Sprintf("Cc: %s\r\n", to.String())
    }
    headers += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
    headers += "\r\n"

    if _, err := fmt.Fprintf(wc, headers+mail.Body); err != nil {
        return err
    }

    return nil
}
