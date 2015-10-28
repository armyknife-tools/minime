package digitalocean

import (
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDigitalOceanDomain_Basic(t *testing.T) {
	var domain godo.Domain

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanDomainDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckDigitalOceanDomainConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDomainExists("digitalocean_domain.foobar", &domain),
					testAccCheckDigitalOceanDomainAttributes(&domain),
					resource.TestCheckResourceAttr(
						"digitalocean_domain.foobar", "name", "foobar-test-terraform.com"),
					resource.TestCheckResourceAttr(
						"digitalocean_domain.foobar", "ip_address", "192.168.0.10"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanDomainDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*godo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_domain" {
			continue
		}

		// Try to find the domain
		_, _, err := client.Domains.Get(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Domain still exists")
		}
	}

	return nil
}

func testAccCheckDigitalOceanDomainAttributes(domain *godo.Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if domain.Name != "foobar-test-terraform.com" {
			return fmt.Errorf("Bad name: %s", domain.Name)
		}

		return nil
	}
}

func testAccCheckDigitalOceanDomainExists(n string, domain *godo.Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*godo.Client)

		foundDomain, _, err := client.Domains.Get(rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundDomain.Name != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*domain = *foundDomain

		return nil
	}
}

const testAccCheckDigitalOceanDomainConfig_basic = `
resource "digitalocean_domain" "foobar" {
    name = "foobar-test-terraform.com"
    ip_address = "192.168.0.10"
}`
