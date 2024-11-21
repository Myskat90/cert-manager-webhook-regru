package main

import (
	regruapi "github.com/daloman/regru-api-go/zonecontrol"
	"github.com/sirupsen/logrus"
)

type RegruClient struct {
	username string
	password string
	zone     string
}

// NewRegruClient создает новый экземпляр клиента для работы с API
func NewRegruClient(username string, password string, zone string) *RegruClient {
	return &RegruClient{
		username: username,
		password: password,
		zone:     zone,
	}
}

// CreateTXT создает TXT запись для домена
func (c *RegruClient) createTXT(domain string, value string) error {
	// Используем метод для добавления TXT записи
	logrus.Debugf("Add %v TXT resource record for %v domain with the following content %v", domain, c.zone, value)
	regruapi.AddTxtRr(c.username, c.password, c.zone, domain, value)
	return nil
}

// DeleteTXT удаляет TXT запись для домена
func (c *RegruClient) deleteTXT(domain string, value string) error {
	logrus.Debugf("Add %v TXT resource record for %v domain with the following content %v", domain, c.zone, value)
	regruapi.RmTxtRr(c.username, c.password, c.zone, domain, "TXT", value)
	return nil
}

// getRecords получает TXT записи для зоны
func (c *RegruClient) getRecords() error {
	regruapi.GetZones(c.username, c.password, c.zone)
	return nil
}
