package main

import (
	"context"
	"net"

	"github.com/Sirupsen/logrus"
	"github.com/cloudflare/cloudflare-go"
)

type CFDNSUpdater struct {
	log    *logrus.Entry
	cf     *cloudflare.API
	zoneId string
}

func NewCFDNSUpdater(zoneId string, apiKey string, apiEmail string, log *logrus.Entry) (*CFDNSUpdater, error) {
	api, err := cloudflare.New(apiKey, apiEmail)
	if err != nil {
		return nil, err
	}

	return &CFDNSUpdater{
		cf:     api,
		zoneId: zoneId,
		log:    log,
	}, nil
}

func NewCFDNSUpdaterToken(zoneId string, apiKey string, log *logrus.Entry) (*CFDNSUpdater, error) {
	api, err := cloudflare.NewWithAPIToken(apiKey)
	if err != nil {
		return nil, err
	}

	return &CFDNSUpdater{
		cf:     api,
		zoneId: zoneId,
		log:    log,
	}, nil
}

func (c *CFDNSUpdater) UpdateRecordA(host string, ip net.IP) error {
	ctx := context.Background()
	// Fetch record IDs for the records we need to update.
	records, err := c.cf.DNSRecords(ctx, c.zoneId, cloudflare.DNSRecord{Name: host, Type: "A"})
	if err != nil {
		return err
	}

	for _, r := range records {
		c.log.Infof("Updating record with ID %s to %s", r.ID, ip.String())
		err := c.cf.UpdateDNSRecord(ctx, c.zoneId, r.ID, cloudflare.DNSRecord{Content: ip.String()})
		if err != nil {
			return err
		}
	}
	return nil
}
