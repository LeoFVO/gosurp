package inspect

import (
	"fmt"

	"github.com/LeoFVO/gosurp/pkg/utils"
	log "github.com/sirupsen/logrus"
)

func (d *Domain) Inspect() error {
	log.Infof("Inspecting domain %s", d.Uri())
	if err := utils.CheckIfDomainRegistered(d.Uri()); err != nil {
		log.Errorf("Domain %s does not exist", d.Uri())
		return fmt.Errorf("Domain %s does not exist", d.Uri())
	}	

	// Get IP records
	ipRecords, err := utils.GetIPRecords(d.Uri())
	if err != nil {
		log.Errorf("Unable to get IP records for domain %s", d.Uri())
		return err
	}
	log.Infof("Domain %s contains the following IP records:\n", d.Uri())
	for _, ip := range ipRecords {
		log.Infof("- %s", ip.String())
	}

	// get MX records
	mxRecords, err := utils.GetMXRecords(d.Uri())
	if err != nil {
		log.Errorf("Unable to get MX records for domain %s", d.Uri())
		return err
	}
	log.Infof("Domain %s contains the following MX records:\n", d.Uri())
	for _, mx := range mxRecords {
		log.Infof("- %s (preference: %d)", mx.Host, mx.Pref)
	}

	// get TXT records
	txtRecords, err := utils.GetTXTRecords(d.Uri())
	if err != nil {
		log.Errorf("Unable to get TXT records for domain %s", d.Uri())
		return err
	}
	log.Infof("Domain %s contains the following TXT records:\n", d.Uri())
	for _, txt := range txtRecords {
		log.Infof("- %s", txt)
	}

	// get NS records
	nsRecords, err := utils.GetNSRecords(d.Uri())
	if err != nil {
		log.Errorf("Unable to get NS records for domain %s", d.Uri())
		return err
	}
	log.Infof("Domain %s contains the following NS records:\n", d.Uri())
	for _, ns := range nsRecords {
		log.Infof("- %s", ns.Host)
	}
	
	return nil
}