// +build !core

//
// This file is automatically generated by scripts/generate-plugins.go -- Do not edit!
//
package command

import (
	archiveprovider "github.com/hashicorp/terraform/builtin/providers/archive"
	atlasprovider "github.com/hashicorp/terraform/builtin/providers/atlas"
	awsprovider "github.com/hashicorp/terraform/builtin/providers/aws"
	azureprovider "github.com/hashicorp/terraform/builtin/providers/azure"
	azurermprovider "github.com/hashicorp/terraform/builtin/providers/azurerm"
	bitbucketprovider "github.com/hashicorp/terraform/builtin/providers/bitbucket"
	chefprovider "github.com/hashicorp/terraform/builtin/providers/chef"
	clcprovider "github.com/hashicorp/terraform/builtin/providers/clc"
	cloudflareprovider "github.com/hashicorp/terraform/builtin/providers/cloudflare"
	cloudstackprovider "github.com/hashicorp/terraform/builtin/providers/cloudstack"
	cobblerprovider "github.com/hashicorp/terraform/builtin/providers/cobbler"
	consulprovider "github.com/hashicorp/terraform/builtin/providers/consul"
	datadogprovider "github.com/hashicorp/terraform/builtin/providers/datadog"
	digitaloceanprovider "github.com/hashicorp/terraform/builtin/providers/digitalocean"
	dmeprovider "github.com/hashicorp/terraform/builtin/providers/dme"
	dnsimpleprovider "github.com/hashicorp/terraform/builtin/providers/dnsimple"
	dockerprovider "github.com/hashicorp/terraform/builtin/providers/docker"
	dynprovider "github.com/hashicorp/terraform/builtin/providers/dyn"
	fastlyprovider "github.com/hashicorp/terraform/builtin/providers/fastly"
	githubprovider "github.com/hashicorp/terraform/builtin/providers/github"
	googleprovider "github.com/hashicorp/terraform/builtin/providers/google"
	grafanaprovider "github.com/hashicorp/terraform/builtin/providers/grafana"
	herokuprovider "github.com/hashicorp/terraform/builtin/providers/heroku"
	influxdbprovider "github.com/hashicorp/terraform/builtin/providers/influxdb"
	libratoprovider "github.com/hashicorp/terraform/builtin/providers/librato"
	logentriesprovider "github.com/hashicorp/terraform/builtin/providers/logentries"
	mailgunprovider "github.com/hashicorp/terraform/builtin/providers/mailgun"
	mysqlprovider "github.com/hashicorp/terraform/builtin/providers/mysql"
	nomadprovider "github.com/hashicorp/terraform/builtin/providers/nomad"
	nullprovider "github.com/hashicorp/terraform/builtin/providers/null"
	openstackprovider "github.com/hashicorp/terraform/builtin/providers/openstack"
	packetprovider "github.com/hashicorp/terraform/builtin/providers/packet"
	pagerdutyprovider "github.com/hashicorp/terraform/builtin/providers/pagerduty"
	postgresqlprovider "github.com/hashicorp/terraform/builtin/providers/postgresql"
	powerdnsprovider "github.com/hashicorp/terraform/builtin/providers/powerdns"
	rabbitmqprovider "github.com/hashicorp/terraform/builtin/providers/rabbitmq"
	rancherprovider "github.com/hashicorp/terraform/builtin/providers/rancher"
	randomprovider "github.com/hashicorp/terraform/builtin/providers/random"
	rundeckprovider "github.com/hashicorp/terraform/builtin/providers/rundeck"
	scalewayprovider "github.com/hashicorp/terraform/builtin/providers/scaleway"
	softlayerprovider "github.com/hashicorp/terraform/builtin/providers/softlayer"
	statuscakeprovider "github.com/hashicorp/terraform/builtin/providers/statuscake"
	templateprovider "github.com/hashicorp/terraform/builtin/providers/template"
	terraformprovider "github.com/hashicorp/terraform/builtin/providers/terraform"
	testprovider "github.com/hashicorp/terraform/builtin/providers/test"
	tlsprovider "github.com/hashicorp/terraform/builtin/providers/tls"
	tritonprovider "github.com/hashicorp/terraform/builtin/providers/triton"
	ultradnsprovider "github.com/hashicorp/terraform/builtin/providers/ultradns"
	vaultprovider "github.com/hashicorp/terraform/builtin/providers/vault"
	vcdprovider "github.com/hashicorp/terraform/builtin/providers/vcd"
	vsphereprovider "github.com/hashicorp/terraform/builtin/providers/vsphere"
	chefresourceprovisioner "github.com/hashicorp/terraform/builtin/provisioners/chef"
	fileresourceprovisioner "github.com/hashicorp/terraform/builtin/provisioners/file"
	localexecresourceprovisioner "github.com/hashicorp/terraform/builtin/provisioners/local-exec"
	remoteexecresourceprovisioner "github.com/hashicorp/terraform/builtin/provisioners/remote-exec"

	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

var InternalProviders = map[string]plugin.ProviderFunc{
	"archive":      archiveprovider.Provider,
	"atlas":        atlasprovider.Provider,
	"aws":          awsprovider.Provider,
	"azure":        azureprovider.Provider,
	"azurerm":      azurermprovider.Provider,
	"bitbucket":    bitbucketprovider.Provider,
	"chef":         chefprovider.Provider,
	"clc":          clcprovider.Provider,
	"cloudflare":   cloudflareprovider.Provider,
	"cloudstack":   cloudstackprovider.Provider,
	"cobbler":      cobblerprovider.Provider,
	"consul":       consulprovider.Provider,
	"datadog":      datadogprovider.Provider,
	"digitalocean": digitaloceanprovider.Provider,
	"dme":          dmeprovider.Provider,
	"dnsimple":     dnsimpleprovider.Provider,
	"docker":       dockerprovider.Provider,
	"dyn":          dynprovider.Provider,
	"fastly":       fastlyprovider.Provider,
	"github":       githubprovider.Provider,
	"google":       googleprovider.Provider,
	"grafana":      grafanaprovider.Provider,
	"heroku":       herokuprovider.Provider,
	"influxdb":     influxdbprovider.Provider,
	"librato":      libratoprovider.Provider,
	"logentries":   logentriesprovider.Provider,
	"mailgun":      mailgunprovider.Provider,
	"mysql":        mysqlprovider.Provider,
	"nomad":        nomadprovider.Provider,
	"null":         nullprovider.Provider,
	"openstack":    openstackprovider.Provider,
	"packet":       packetprovider.Provider,
	"pagerduty":    pagerdutyprovider.Provider,
	"postgresql":   postgresqlprovider.Provider,
	"powerdns":     powerdnsprovider.Provider,
	"rabbitmq":     rabbitmqprovider.Provider,
	"rancher":      rancherprovider.Provider,
	"random":       randomprovider.Provider,
	"rundeck":      rundeckprovider.Provider,
	"scaleway":     scalewayprovider.Provider,
	"softlayer":    softlayerprovider.Provider,
	"statuscake":   statuscakeprovider.Provider,
	"template":     templateprovider.Provider,
	"terraform":    terraformprovider.Provider,
	"test":         testprovider.Provider,
	"tls":          tlsprovider.Provider,
	"triton":       tritonprovider.Provider,
	"ultradns":     ultradnsprovider.Provider,
	"vault":        vaultprovider.Provider,
	"vcd":          vcdprovider.Provider,
	"vsphere":      vsphereprovider.Provider,
}

var InternalProvisioners = map[string]plugin.ProvisionerFunc{
	"chef":        func() terraform.ResourceProvisioner { return new(chefresourceprovisioner.ResourceProvisioner) },
	"file":        func() terraform.ResourceProvisioner { return new(fileresourceprovisioner.ResourceProvisioner) },
	"local-exec":  func() terraform.ResourceProvisioner { return new(localexecresourceprovisioner.ResourceProvisioner) },
	"remote-exec": func() terraform.ResourceProvisioner { return new(remoteexecresourceprovisioner.ResourceProvisioner) },
}
