package azurerm

import (
	"fmt"
	"reflect"
	"strings"

	"sync"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform/helper/mutexkv"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	riviera "github.com/jen20/riviera/azure"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	var p *schema.Provider
	p = &schema.Provider{
		Schema: map[string]*schema.Schema{
			"subscription_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_SUBSCRIPTION_ID", ""),
			},

			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_ID", ""),
			},

			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_SECRET", ""),
			},

			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_TENANT_ID", ""),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"azurerm_client_config": dataSourceArmClientConfig(),
		},

		ResourcesMap: map[string]*schema.Resource{
			// These resources use the Azure ARM SDK
			"azurerm_availability_set": resourceArmAvailabilitySet(),
			"azurerm_cdn_endpoint":     resourceArmCdnEndpoint(),
			"azurerm_cdn_profile":      resourceArmCdnProfile(),

			"azurerm_eventhub":                resourceArmEventHub(),
			"azurerm_eventhub_consumer_group": resourceArmEventHubConsumerGroup(),
			"azurerm_eventhub_namespace":      resourceArmEventHubNamespace(),

			"azurerm_lb":                      resourceArmLoadBalancer(),
			"azurerm_lb_backend_address_pool": resourceArmLoadBalancerBackendAddressPool(),
			"azurerm_lb_nat_rule":             resourceArmLoadBalancerNatRule(),
			"azurerm_lb_nat_pool":             resourceArmLoadBalancerNatPool(),
			"azurerm_lb_probe":                resourceArmLoadBalancerProbe(),
			"azurerm_lb_rule":                 resourceArmLoadBalancerRule(),

			"azurerm_key_vault":                 resourceArmKeyVault(),
			"azurerm_local_network_gateway":     resourceArmLocalNetworkGateway(),
			"azurerm_network_interface":         resourceArmNetworkInterface(),
			"azurerm_network_security_group":    resourceArmNetworkSecurityGroup(),
			"azurerm_network_security_rule":     resourceArmNetworkSecurityRule(),
			"azurerm_public_ip":                 resourceArmPublicIp(),
			"azurerm_redis_cache":               resourceArmRedisCache(),
			"azurerm_route":                     resourceArmRoute(),
			"azurerm_route_table":               resourceArmRouteTable(),
			"azurerm_servicebus_namespace":      resourceArmServiceBusNamespace(),
			"azurerm_servicebus_subscription":   resourceArmServiceBusSubscription(),
			"azurerm_servicebus_topic":          resourceArmServiceBusTopic(),
			"azurerm_storage_account":           resourceArmStorageAccount(),
			"azurerm_storage_blob":              resourceArmStorageBlob(),
			"azurerm_storage_container":         resourceArmStorageContainer(),
			"azurerm_storage_share":             resourceArmStorageShare(),
			"azurerm_storage_queue":             resourceArmStorageQueue(),
			"azurerm_storage_table":             resourceArmStorageTable(),
			"azurerm_subnet":                    resourceArmSubnet(),
			"azurerm_template_deployment":       resourceArmTemplateDeployment(),
			"azurerm_traffic_manager_endpoint":  resourceArmTrafficManagerEndpoint(),
			"azurerm_traffic_manager_profile":   resourceArmTrafficManagerProfile(),
			"azurerm_virtual_machine_extension": resourceArmVirtualMachineExtensions(),
			"azurerm_virtual_machine":           resourceArmVirtualMachine(),
			"azurerm_virtual_machine_scale_set": resourceArmVirtualMachineScaleSet(),
			"azurerm_virtual_network":           resourceArmVirtualNetwork(),
			"azurerm_virtual_network_peering":   resourceArmVirtualNetworkPeering(),

			// These resources use the Riviera SDK
			"azurerm_dns_a_record":      resourceArmDnsARecord(),
			"azurerm_dns_aaaa_record":   resourceArmDnsAAAARecord(),
			"azurerm_dns_cname_record":  resourceArmDnsCNameRecord(),
			"azurerm_dns_mx_record":     resourceArmDnsMxRecord(),
			"azurerm_dns_ns_record":     resourceArmDnsNsRecord(),
			"azurerm_dns_srv_record":    resourceArmDnsSrvRecord(),
			"azurerm_dns_txt_record":    resourceArmDnsTxtRecord(),
			"azurerm_dns_zone":          resourceArmDnsZone(),
			"azurerm_resource_group":    resourceArmResourceGroup(),
			"azurerm_search_service":    resourceArmSearchService(),
			"azurerm_sql_database":      resourceArmSqlDatabase(),
			"azurerm_sql_firewall_rule": resourceArmSqlFirewallRule(),
			"azurerm_sql_server":        resourceArmSqlServer(),
		},
	}

	p.ConfigureFunc = providerConfigure(p)

	return p
}

// Config is the configuration structure used to instantiate a
// new Azure management client.
type Config struct {
	ManagementURL string

	SubscriptionID string
	ClientID       string
	ClientSecret   string
	TenantID       string

	validateCredentialsOnce sync.Once
}

func (c *Config) validate() error {
	var err *multierror.Error

	if c.SubscriptionID == "" {
		err = multierror.Append(err, fmt.Errorf("Subscription ID must be configured for the AzureRM provider"))
	}
	if c.ClientID == "" {
		err = multierror.Append(err, fmt.Errorf("Client ID must be configured for the AzureRM provider"))
	}
	if c.ClientSecret == "" {
		err = multierror.Append(err, fmt.Errorf("Client Secret must be configured for the AzureRM provider"))
	}
	if c.TenantID == "" {
		err = multierror.Append(err, fmt.Errorf("Tenant ID must be configured for the AzureRM provider"))
	}

	return err.ErrorOrNil()
}

func providerConfigure(p *schema.Provider) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		config := &Config{
			SubscriptionID: d.Get("subscription_id").(string),
			ClientID:       d.Get("client_id").(string),
			ClientSecret:   d.Get("client_secret").(string),
			TenantID:       d.Get("tenant_id").(string),
		}

		if err := config.validate(); err != nil {
			return nil, err
		}

		client, err := config.getArmClient()
		if err != nil {
			return nil, err
		}

		client.StopContext = p.StopContext()

		err = registerAzureResourceProvidersWithSubscription(client.rivieraClient)
		if err != nil {
			return nil, err
		}

		return client, nil
	}
}

func registerProviderWithSubscription(providerName string, client *riviera.Client) error {
	request := client.NewRequest()
	request.Command = riviera.RegisterResourceProvider{
		Namespace: providerName,
	}

	response, err := request.Execute()
	if err != nil {
		return fmt.Errorf("Cannot request provider registration for Azure Resource Manager: %s.", err)
	}

	if !response.IsSuccessful() {
		return fmt.Errorf("Credentials for accessing the Azure Resource Manager API are likely " +
			"to be incorrect, or\n  the service principal does not have permission to use " +
			"the Azure Service Management\n  API.")
	}

	return nil
}

var providerRegistrationOnce sync.Once

// registerAzureResourceProvidersWithSubscription uses the providers client to register
// all Azure resource providers which the Terraform provider may require (regardless of
// whether they are actually used by the configuration or not). It was confirmed by Microsoft
// that this is the approach their own internal tools also take.
func registerAzureResourceProvidersWithSubscription(client *riviera.Client) error {
	var err error
	providerRegistrationOnce.Do(func() {
		// We register Microsoft.Compute during client initialization
		providers := []string{
			"Microsoft.Cache",
			"Microsoft.Network",
			"Microsoft.Cdn",
			"Microsoft.Storage",
			"Microsoft.Sql",
			"Microsoft.Search",
			"Microsoft.Resources",
			"Microsoft.ServiceBus",
			"Microsoft.KeyVault",
			"Microsoft.EventHub",
		}

		var wg sync.WaitGroup
		wg.Add(len(providers))
		for _, providerName := range providers {
			go func(p string) {
				defer wg.Done()
				if innerErr := registerProviderWithSubscription(p, client); err != nil {
					err = innerErr
				}
			}(providerName)
		}
		wg.Wait()
	})

	return err
}

// armMutexKV is the instance of MutexKV for ARM resources
var armMutexKV = mutexkv.NewMutexKV()

func azureStateRefreshFunc(resourceURI string, client *ArmClient, command riviera.APICall) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		req := client.rivieraClient.NewRequestForURI(resourceURI)
		req.Command = command

		res, err := req.Execute()
		if err != nil {
			return nil, "", fmt.Errorf("Error executing %T command in azureStateRefreshFunc", req.Command)
		}

		var value reflect.Value
		if reflect.ValueOf(res.Parsed).Kind() == reflect.Ptr {
			value = reflect.ValueOf(res.Parsed).Elem()
		} else {
			value = reflect.ValueOf(res.Parsed)
		}

		for i := 0; i < value.NumField(); i++ { // iterates through every struct type field
			tag := value.Type().Field(i).Tag // returns the tag string
			tagValue := tag.Get("mapstructure")
			if tagValue == "provisioningState" {
				return res.Parsed, value.Field(i).Elem().String(), nil
			}
		}

		panic(fmt.Errorf("azureStateRefreshFunc called on structure %T with no mapstructure:provisioningState tag. This is a bug", res.Parsed))
	}
}

// Resource group names can be capitalised, but we store them in lowercase.
// Use a custom diff function to avoid creation of new resources.
func resourceAzurermResourceGroupNameDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	return strings.ToLower(old) == strings.ToLower(new)
}
