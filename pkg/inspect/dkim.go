package inspect

import (
	"fmt"

	"github.com/LeoFVO/gosurp/pkg/utils"
	log "github.com/sirupsen/logrus"
)

func (d *Domain) InspectDKIM() error {
	log.Infof("Inspecting domain %s", d.Uri())
	if err := utils.CheckIfDomainRegistered(d.Uri()); err != nil {
		log.Errorf("Domain %s does not exist", d.Uri())
		return fmt.Errorf("Domain %s does not exist", d.Uri())
	}	


	// get DKIM records
	dkimRecords, err := utils.GetDKIMRecords(d.Uri())
	if err != nil {
		log.Errorf("Unable to get DKIM records for domain %s", d.Uri())
		return err
	}
	log.Infof("Domain %s contains the following DKIM records:\n", d.Uri())
	for _, dkim := range dkimRecords {
		log.Infof("- %s", dkim)
	}

	return nil
}