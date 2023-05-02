package utils

import (
	"net"
	"strings"

	log "github.com/sirupsen/logrus"
)

func IsIp(value string) bool {
	if value == "localhost" {
		return true
	}
	return net.ParseIP(value) != nil
}

func GetSMTPServerAddress(domain string) string {
	var server string
	if!IsIp(domain) {
		log.Debugf("SMTP server address is not an IP address, looking up MX records")
		mxRecords, err := GetMXRecords(domain)
		if len(mxRecords) == 0 || err != nil {
				log.Warnf("no MX records found for %s", domain)
		}
		server = strings.TrimSuffix(mxRecords[0].Host, ".")
		log.Debugf("Using %s as SMTP server address", server)
	} else {
		log.Debugf("SMTP server address is an IP address, using it as is")
		server = domain
	}
	return server
}

func GetMXRecords(domain string) ([]*net.MX, error) {
	log.Tracef("Looking up MX records for %s", domain)
	return net.LookupMX(domain)
}

func GetIPRecords(domain string) ([]net.IP, error) {
	log.Tracef("Looking up IP records for %s", domain)
	return net.LookupIP(domain)
}

func GetTXTRecords(domain string) ([]string, error) {
	log.Tracef("Looking up TXT records for %s", domain)
	return net.LookupTXT(domain)
}

func GetNSRecords(domain string) ([]*net.NS, error) {
	log.Tracef("Looking up NS records for %s", domain)
	return net.LookupNS(domain)
}

func GetPTRRecords(domain string) ([]string, error) {
	log.Tracef("Looking up PTR records for %s", domain)
	return net.LookupAddr(domain)
}