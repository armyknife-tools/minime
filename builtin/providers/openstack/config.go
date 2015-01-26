package openstack

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
)

type Config struct {
	Username         string
	UserID           string
	Password         string
	APIKey           string
	IdentityEndpoint string
	TenantID         string
	TenantName       string
	DomainID         string
	DomainName       string

	osClient *gophercloud.ProviderClient
}

func (c *Config) loadAndValidate() error {
	ao := gophercloud.AuthOptions{
		Username:         c.Username,
		UserID:           c.UserID,
		Password:         c.Password,
		APIKey:           c.APIKey,
		IdentityEndpoint: c.IdentityEndpoint,
		TenantID:         c.TenantID,
		TenantName:       c.TenantName,
		DomainID:         c.DomainID,
		DomainName:       c.DomainName,
	}

	client, err := openstack.AuthenticatedClient(ao)
	if err != nil {
		return err
	}

	c.osClient = client

	return nil
}
