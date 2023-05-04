package smtp

import "net/mail"

type Server struct {
    Hostname string `yaml:"host"`
    Port string `yaml:"port"`
    Ports []string `yaml:"ports"`
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

type DKIM struct {
    Selector string `yaml:"selector"`
    PrivateKey string `yaml:"private_key"`
    Domain string `yaml:"domain"`
}