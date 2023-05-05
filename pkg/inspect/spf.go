package inspect

import (
	"fmt"

	"github.com/LeoFVO/gosurp/pkg/utils"
	log "github.com/sirupsen/logrus"
)

func (d *Domain) InspectSPF() error {
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

	// get TXT records
	txtRecords, err := utils.GetTXTRecords(d.Uri())
	if err != nil {
		log.Errorf("Unable to get TXT records for domain %s", d.Uri())
		return err
	}
	log.Infof("Domain %s contains the following SPF records:\n", d.Uri())
	for _, txt := range txtRecords {
		// check if record is SPF
		if utils.IsSPFRecord(txt) {
			log.Infof("Found SPF record: %s", txt)
		}
	}

	// TODO: Analyse SPF records

	return nil
}
