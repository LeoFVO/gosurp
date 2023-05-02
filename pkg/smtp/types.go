package smtp

import "net/mail"

type Server struct {
    Hostname string `yaml:"host"`
    Port string `yaml:"port"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
}

// Mail represents an email message.
type Envelope struct {
    From    mail.Address `yaml:"from"`
    To      []mail.Address `yaml:"to"`
    Subject string `yaml:"subject"`
    Body    string `yaml:"body"`
}
