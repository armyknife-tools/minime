package opc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/compute"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// TODO (@jake): Properly create a vNIC Set once instances are finished
func TestAccOPCRoute_Basic(t *testing.T) {
	rInt := acctest.RandInt()
	resName := "opc_compute_route.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccOPCCheckRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOPCRouteConfig_Basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccOPCCheckRouteExists,
					resource.TestCheckResourceAttr(resName, "admin_distance", "1"),
					resource.TestCheckResourceAttr(resName, "ip_address_prefix", "10.0.12.0/24"),
					resource.TestCheckResourceAttr(resName, "name", fmt.Sprintf("testing-route-%d", rInt)),
				),
			},
			{
				Config: testAccOPCRouteConfig_BasicUpdate(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccOPCCheckRouteExists,
					resource.TestCheckResourceAttr(resName, "admin_distance", "2"),
					resource.TestCheckResourceAttr(resName, "ip_address_prefix", "10.0.14.0/24"),
				),
			},
		},
	})
}

// TODO (@jake): Properly create a vNIC Set once instances are finished
func testAccOPCRouteConfig_Basic(rInt int) string {
	return fmt.Sprintf(`
resource "opc_compute_vnic_set" "test" {
  name = "route-test-%d"
  description = "route-testing-%d"
  virtual_nics = ["jake-manual_eth1"]
}

resource "opc_compute_route" "test" {
  name = "testing-route-%d"
  description = "testing-desc-%d"
  admin_distance = 1
  ip_address_prefix = "10.0.12.0/24"
  next_hop_vnic_set = "${opc_compute_vnic_set.test.name}"
}`, rInt, rInt, rInt, rInt)
}

func testAccOPCRouteConfig_BasicUpdate(rInt int) string {
	return fmt.Sprintf(`
resource "opc_compute_vnic_set" "test" {
  name = "route-test-%d"
  description = "route-testing-%d"
  virtual_nics = ["jake-manual_eth1"]
}

resource "opc_compute_route" "test" {
  name = "testing-route-%d"
  description = "testing-desc-%d"
  admin_distance = 2
  ip_address_prefix = "10.0.14.0/24"
  next_hop_vnic_set = "${opc_compute_vnic_set.test.name}"
}`, rInt, rInt, rInt, rInt)
}

func testAccOPCCheckRouteExists(s *terraform.State) error {
	client := testAccProvider.Meta().(*compute.Client).Routes()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opc_compute_route" {
			continue
		}

		input := compute.GetRouteInput{
			Name: rs.Primary.Attributes["name"],
		}
		if _, err := client.GetRoute(&input); err != nil {
			return fmt.Errorf("Error retrieving state of Rule %s: %s", input.Name, err)
		}
	}

	return nil
}

func testAccOPCCheckRouteDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*compute.Client).Routes()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opc_compute_route" {
			continue
		}

		input := compute.GetRouteInput{
			Name: rs.Primary.Attributes["name"],
		}
		if info, err := client.GetRoute(&input); err == nil {
			return fmt.Errorf("Rule %s still exists: %#v", input.Name, info)
		}
	}

	return nil
}
