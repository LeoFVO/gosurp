package inspect

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Domain struct {
	Name string
	Extension string
	Subdomains []string
}

func New(value string) (*Domain, error) {
	dns := strings.Split(value, ".")
	if !IsValidDomainName(value) {
		return nil, fmt.Errorf("invalid domain %s", value)
	}
	d := &Domain{
		Name: dns[0],
		Extension: dns[1],
		Subdomains: dns[2:],
	}
	log.Tracef("Domain struct for %s successfully crafted", d.Uri())
	return d, nil

}

func IsValidDomainName(value string) bool {
	dns := strings.Split(value, ".")
	
	// TODO: check if domain has at least 2 parts
	if len(dns) < 2 {
		log.Errorf("Domain %s isn't complete. (should at least contains name.extension)", value)
		return false
	}

	// TODO: check if name is valid
	if len(dns[0]) < 1 || strings.Contains(dns[0], " ") {
		log.Errorf("Domain name %s is invalid", value)
		return false
	}

	// TODO: check if extension is valid
	if len(dns[1]) < 1 || strings.Contains(dns[0], " ") {
		log.Errorf("Domain extension %s is invalid", value)
		return false
	}
		log.Tracef("Domain name %s seem's to be valid", dns[0])

	return true
}

func (d *Domain) Uri() string {
	return d.Name + "." + d.Extension
}