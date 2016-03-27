package dns

import "github.com/jen20/riviera/azure"

type GetMXRecordSetResponse struct {
	ID        string             `mapstructure:"id"`
	Name      string             `mapstructure:"name"`
	Location  string             `mapstructure:"location"`
	Tags      map[string]*string `mapstructure:"tags"`
	TTL       *int               `mapstructure:"TTL"`
	MXRecords []MXRecord         `mapstructure:"MXRecords"`
}

type GetMXRecordSet struct {
	Name              string `json:"-"`
	ResourceGroupName string `json:"-"`
	ZoneName          string `json:"-"`
}

func (command GetMXRecordSet) APIInfo() azure.APIInfo {
	return azure.APIInfo{
		APIVersion:  apiVersion,
		Method:      "GET",
		URLPathFunc: dnsRecordSetDefaultURLPathFunc(command.ResourceGroupName, command.ZoneName, "MX", command.Name),
		ResponseTypeFunc: func() interface{} {
			return &GetMXRecordSetResponse{}
		},
	}
}
