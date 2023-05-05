package inspect

import (
	"fmt"

	"github.com/LeoFVO/gosurp/pkg/utils"
	log "github.com/sirupsen/logrus"
)

func (d *Domain) InspectDMARC() error {
	log.Infof("Inspecting domain %s", d.Uri())
	if err := utils.CheckIfDomainRegistered(d.Uri()); err != nil {
		log.Errorf("Domain %s does not exist", d.Uri())
		return fmt.Errorf("Domain %s does not exist", d.Uri())
	}

	// get DMARC records
	dmarcRecord, err := utils.GetDMARCRecord(d.Uri())
	if err != nil {
		log.Errorf("Unable to get DMARC record for domain %s", d.Uri())
		return err
	}

	log.Infof("Found DMARC record: %s", dmarcRecord)

	// TODO: Analyse DMARC record

	return nil
}