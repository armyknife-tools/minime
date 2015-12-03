package cloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/xanzy/go-cloudstack/cloudstack"
)

func TestAccCloudStackFirewall_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackFirewallDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCloudStackFirewall_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackFirewallRulesExist("cloudstack_firewall.foo"),
					resource.TestCheckResourceAttr(
						"cloudstack_firewall.foo", "ipaddress", CLOUDSTACK_PUBLIC_IPADDRESS),
					resource.TestCheckResourceAttr(
						"cloudstack_firewall.foo", "rule.#", "2"),
					resource.TestCheckResourceAttr(
						"cloudstack_firewall.foo", "rule.60926170.cidr_list.3482919157", "10.0.0.0/24"),
					resource.TestCheckResourceAttr(
						"cloudstack_firewall.foo", "rule.60926170.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"cloudstack_firewall.foo", "rule.60926170.ports.32925333", "8080"),
					resource.TestCheckResourceAttr(
						"cloudstack_firewall.foo", "rule.716592205.source_cidr", "10.0.0.0/24"),
					resource.TestCheckResourceAttr(
						"cloudstack_firewall.foo", "rule.716592205.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"cloudstack_firewall.foo", "rule.716592205.ports.1209010669", "1000-2000"),
					resource.TestCheckResourceAttr(
						"cloudstack_firewall.foo", "rule.716592205.ports.1889509032", "80"),
				),
			},
		},
	})
}

func TestAccCloudStackFirewall_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackFirewallDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCloudStackFirewall_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackFirewallRulesExist("cloudstack_firewall.foo"),
					resource.TestCheckResourceAttr(
						"cloudstack_firewall.foo", "ipaddress", CLOUDSTACK_PUBLIC_IPADDRESS),
					resource.TestCheckResourceAttr(
						"cloudstack_firewall.foo", "rule.#", "2"),
					resource.TestCheckResourceAttr(
						"cloudstack_firewall.foo", "rule.60926170.cidr_list.3482919157", "10.0.0.0/24"),
					resource.TestCheckResourceAttr(
						"cloudstack_firewall.foo", "rule.60926170.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"cloudstack_firewall.foo", "rule.60926170.ports.32925333", "8080"),
					resource.TestCheckResourceAttr(
						"cloudstack_firewall.foo", "rule.716592205.source_cidr", "10.0.0.0/24"),
					resource.TestCheckResourceAttr(
						"cloudstack_firewall.foo", "rule.716592205.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"cloudstack_firewall.foo", "rule.716592205.ports.1209010669", "1000-2000"),
					resource.TestCheckResourceAttr(
						"cloudstack_firewall.foo", "rule.716592205.ports.1889509032", "80"),
				),
			},

			resource.TestStep{
				Config: testAccCloudStackFirewall_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackFirewallRulesExist("cloudstack_firewall.foo"),
					resource.TestCheckResourceAttr(
						"cloudstack_firewall.foo", "ipaddress", CLOUDSTACK_PUBLIC_IPADDRESS),
					resource.TestCheckResourceAttr(
						"cloudstack_firewall.foo", "rule.#", "3"),
					resource.TestCheckResourceAttr(
						"cloudstack_firewall.foo", "rule.2207610982.cidr_list.80081744", "10.0.1.0/24"),
					resource.TestCheckResourceAttr(
						"cloudstack_firewall.foo", "rule.2207610982.cidr_list.3482919157", "10.0.0.0/24"),
					resource.TestCheckResourceAttr(
						"cloudstack_firewall.foo", "rule.2207610982.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"cloudstack_firewall.foo", "rule.2207610982.ports.32925333", "8080"),
					resource.TestCheckResourceAttr(
						"cloudstack_firewall.foo", "rule.716592205.source_cidr", "10.0.0.0/24"),
					resource.TestCheckResourceAttr(
						"cloudstack_firewall.foo", "rule.716592205.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"cloudstack_firewall.foo", "rule.716592205.ports.1209010669", "1000-2000"),
					resource.TestCheckResourceAttr(
						"cloudstack_firewall.foo", "rule.716592205.ports.1889509032", "80"),
					resource.TestCheckResourceAttr(
						"cloudstack_firewall.foo", "rule.4449157.source_cidr", "172.16.100.0/24"),
					resource.TestCheckResourceAttr(
						"cloudstack_firewall.foo", "rule.4449157.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"cloudstack_firewall.foo", "rule.4449157.ports.1889509032", "80"),
					resource.TestCheckResourceAttr(
						"cloudstack_firewall.foo", "rule.4449157.ports.3638101695", "443"),
				),
			},
		},
	})
}

func testAccCheckCloudStackFirewallRulesExist(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No firewall ID is set")
		}

		for k, id := range rs.Primary.Attributes {
			if !strings.Contains(k, ".uuids.") || strings.HasSuffix(k, ".uuids.#") {
				continue
			}

			cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
			_, count, err := cs.Firewall.GetFirewallRuleByID(id)

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

func testAccCheckCloudStackFirewallDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_firewall" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No instance ID is set")
		}

		for k, id := range rs.Primary.Attributes {
			if !strings.Contains(k, ".uuids.") || strings.HasSuffix(k, ".uuids.#") {
				continue
			}

			_, _, err := cs.Firewall.GetFirewallRuleByID(id)
			if err == nil {
				return fmt.Errorf("Firewall rule %s still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

var testAccCloudStackFirewall_basic = fmt.Sprintf(`
resource "cloudstack_firewall" "foo" {
  ipaddress = "%s"

  rule {
    cidr_list = ["10.0.0.0/24"]
    protocol = "tcp"
    ports = ["8080"]
  }

  rule {
    source_cidr = "10.0.0.0/24"
    protocol = "tcp"
    ports = ["80", "1000-2000"]
  }
}`, CLOUDSTACK_PUBLIC_IPADDRESS)

var testAccCloudStackFirewall_update = fmt.Sprintf(`
resource "cloudstack_firewall" "foo" {
  ipaddress = "%s"

  rule {
    cidr_list = ["10.0.0.0/24", "10.0.1.0/24"]
    protocol = "tcp"
    ports = ["8080"]
  }

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
}`, CLOUDSTACK_PUBLIC_IPADDRESS)
