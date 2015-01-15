package cloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/xanzy/go-cloudstack/cloudstack"
)

func TestAccCloudStackEgressFirewall_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackEgressFirewallDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCloudStackEgressFirewall_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackEgressFirewallRulesExist("cloudstack_egress_firewall.foo"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "network", CLOUDSTACK_NETWORK_1),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.1702320581.source_cidr", "10.0.0.0/24"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.1702320581.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.1702320581.ports.#", "2"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.1702320581.ports.1209010669", "1000-2000"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.1702320581.ports.1889509032", "80"),
				),
			},
		},
	})
}

func TestAccCloudStackEgressFirewall_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackEgressFirewallDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCloudStackEgressFirewall_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackEgressFirewallRulesExist("cloudstack_egress_firewall.foo"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "ipaddress", CLOUDSTACK_NETWORK_1),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.#", "1"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.1702320581.source_cidr", "10.0.0.0/24"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.1702320581.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.1702320581.ports.#", "2"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.1702320581.ports.1209010669", "1000-2000"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.1702320581.ports.1889509032", "80"),
				),
			},

			resource.TestStep{
				Config: testAccCloudStackEgressFirewall_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackEgressFirewallRulesExist("cloudstack_egress_firewall.foo"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "ipaddress", CLOUDSTACK_NETWORK_1),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.#", "2"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.1702320581.source_cidr", "10.0.0.0/24"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.1702320581.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.1702320581.ports.#", "2"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.1702320581.ports.1209010669", "1000-2000"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.1702320581.ports.1889509032", "80"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.3779782959.source_cidr", "172.16.100.0/24"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.3779782959.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.3779782959.ports.#", "2"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.3779782959.ports.1889509032", "80"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.3779782959.ports.3638101695", "443"),
				),
			},
		},
	})
}

func testAccCheckCloudStackEgressFirewallRulesExist(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No firewall ID is set")
		}

		for k, uuid := range rs.Primary.Attributes {
			if !strings.Contains(k, "uuids") {
				continue
			}

			cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
			_, count, err := cs.Firewall.GetEgressFirewallRuleByID(uuid)

			if err != nil {
				return err
			}

			if count == 0 {
				return fmt.Errorf("Firewall rule for %s not found", k)
			}
		}

		return nil
	}
}

func testAccCheckCloudStackEgressFirewallDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_egress_firewall" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No instance ID is set")
		}

		for k, uuid := range rs.Primary.Attributes {
			if !strings.Contains(k, "uuids") {
				continue
			}

			p := cs.Firewall.NewDeleteEgressFirewallRuleParams(uuid)
			_, err := cs.Firewall.DeleteEgressFirewallRule(p)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

var testAccCloudStackEgressFirewall_basic = fmt.Sprintf(`
resource "cloudstack_egress_firewall" "foo" {
  ipaddress = "%s"

  rule {
    source_cidr = "10.0.0.0/24"
    protocol = "tcp"
    ports = ["80", "1000-2000"]
  }
}`, CLOUDSTACK_NETWORK_1)

var testAccCloudStackEgressFirewall_update = fmt.Sprintf(`
resource "cloudstack_egress_firewall" "foo" {
  ipaddress = "%s"

  rule {
    source_cidr = "10.0.0.0/24"
    protocol = "tcp"
    ports = ["80", "1000-2000"]
  }

  rule {
    source_cidr = "172.16.100.0/24"
    protocol = "tcp"
    ports = ["80", "443"]
  }
}`, CLOUDSTACK_NETWORK_1)
