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
		log.Debugf("SMTP server address is an IP address or localhost, using it as is")
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

func CheckIfDomainRegistered(domain string) error {
	_, err := net.LookupHost(domain)
	if err != nil {
		log.Debugf("Domain %s doesn't seem to be registered", domain)
		return err
	}
	log.Debugf("Domain %s is registered", domain)
	return nil
}

func IsSPFRecord(record string) bool {
	return strings.HasPrefix(record, "v=spf1")
}

func GetDMARCRecord(domain string) (string, error) {
	log.Tracef("Looking up DMARC record for %s", domain)

	txtRecords, err := GetTXTRecords("_dmarc." + domain)
	if err != nil {
		log.Debugf("No DMARC record found for %s", domain)
		return "", err
	}

	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			log.Debugf("DMARC record found for %s", domain)
			return record, nil
		}
	}

	log.Debugf("No DMARC record found for %s", domain)
	return "", nil
}

func GetDKIMRecords(domain string) ([]string, error) {
	log.Tracef("Looking up DKIM records for %s", domain)
	var dkimRecords []string

	// check if DKIM records are in TXT records
	txtRecords, err := GetTXTRecords("_domainkey." + domain)
	if err != nil {
		log.Debugf("No DKIM record found for %s using TXT records", domain)
	}
	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=DKIM1") {
			log.Debugf("DKIM record found for %s", domain)
			dkimRecords = append(dkimRecords, record)
		}
	}

	// TODO: check if default._domainkey is a valid prefix using CNAME
	dkimPrefixWords := []string{"default", "selector1", "selector2", "firebase1", "firebase2", "s1", "s2"}

	for _, prefix := range dkimPrefixWords {
		cnameRecords, err := net.LookupCNAME(prefix + "._domainkey." + domain)
		if err != nil {
			log.Debugf("No DKIM record found for %s using CNAME records", domain)
			continue
		}
		dkimRecords = append(dkimRecords, cnameRecords)
	}

	if len(dkimRecords) == 0 {
		log.Debugf("No DKIM record found for %s", domain)
		return nil, nil
	}

	return dkimRecords, nil
}