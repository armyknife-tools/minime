package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/volumeattach"
)

func TestAccComputeV2VolumeAttach_basic(t *testing.T) {
	var va volumeattach.VolumeAttachment

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2VolumeAttachDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeV2VolumeAttach_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2VolumeAttachExists(t, "openstack_compute_volume_attach_v2.va_1", &va),
				),
			},
		},
	})
}

func testAccCheckComputeV2VolumeAttachDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	computeClient, err := config.computeV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_compute_volume_attach_v2" {
			continue
		}

		instanceId, volumeId, err := parseComputeVolumeAttachmentId(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = volumeattach.Get(computeClient, instanceId, volumeId).Extract()
		if err == nil {
			return fmt.Errorf("Volume attachment still exists")
		}
	}

	return nil
}

func testAccCheckComputeV2VolumeAttachExists(t *testing.T, n string, va *volumeattach.VolumeAttachment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		computeClient, err := config.computeV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("(testAccCheckComputeV2VolumeAttachExists) Error creating OpenStack compute client: %s", err)
		}

		instanceId, volumeId, err := parseComputeVolumeAttachmentId(rs.Primary.ID)
		if err != nil {
			return err
		}

		found, err := volumeattach.Get(computeClient, instanceId, volumeId).Extract()
		if err != nil {
			return err
		}

		if found.ServerID != instanceId || found.VolumeID != volumeId {
			return fmt.Errorf("VolumeAttach not found")
		}

		*va = *found

		return nil
	}
}

var testAccComputeV2VolumeAttach_basic = `
	resource "openstack_blockstorage_volume_v2" "volume_1" {
		name = "volume_1"
		size = 1
	}

	resource "openstack_compute_instance_v2" "instance_1" {
		name = "instance_1"
		security_groups = ["default"]
	}

	resource "openstack_compute_volume_attach_v2" "va_1" {
		instance_id = "${openstack_compute_instance_v2.instance_1.id}"
		volume_id = "${openstack_blockstorage_volume_v2.volume_1.id}"
	}
`
