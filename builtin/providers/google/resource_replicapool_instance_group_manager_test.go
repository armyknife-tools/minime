package google

import (
	"fmt"
	"log"
	"testing"

	"code.google.com/p/google-api-go-client/replicapool/v1beta2"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccInstanceGroupManager_basic(t *testing.T) {
	var manager replicapool.InstanceGroupManager

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceGroupManagerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccInstanceGroupManager_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceGroupManagerExists(
						"google_replicapool_instance_group_manager.foobar", &manager),
				),
			},
		},
	})
}

func TestAccInstanceGroupManager_update(t *testing.T) {
	var manager replicapool.InstanceGroupManager

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceGroupManagerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccInstanceGroupManager_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceGroupManagerExists(
						"google_replicapool_instance_group_manager.foobar", &manager),
				),
			},
			resource.TestStep{
				Config: testAccInstanceGroupManager_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceGroupManagerExists(
						"google_replicapool_instance_group_manager.foobar", &manager),
				),
			},
			resource.TestStep{
				Config: testAccInstanceGroupManager_update2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceGroupManagerExists(
						"google_replicapool_instance_group_manager.foobar", &manager),
					testAccCheckInstanceGroupManagerUpdated(
						"google_replicapool_instance_group_manager.foobar", 3,
						"google_compute_target_pool.foobaz", "terraform-test-foobaz"),
				),
			},
		},
	})
}

func testAccCheckInstanceGroupManagerDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_replicapool_instance_group_manager" {
			continue
		}
		_, err := config.clientReplicaPool.InstanceGroupManagers.Get(
			config.Project, rs.Primary.Attributes["zone"], rs.Primary.ID).Do()
		if err != nil {
			return fmt.Errorf("InstanceGroupManager still exists")
		}
	}

	return nil
}

func testAccCheckInstanceGroupManagerExists(n string, manager *replicapool.InstanceGroupManager) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.clientReplicaPool.InstanceGroupManagers.Get(
			config.Project, rs.Primary.Attributes["zone"], rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("InstanceGroupManager not found")
		}

		*manager = *found

		return nil
	}
}

func testAccCheckInstanceGroupManagerUpdated(n string, size int64, targetPool string, template string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		log.Printf("[DEBUG] XXXXXXXXXXXXXXXXXXXXXXXX Manager Test: %#v", rs)

		config := testAccProvider.Meta().(*Config)

		manager, err := config.clientReplicaPool.InstanceGroupManagers.Get(
			config.Project, rs.Primary.Attributes["zone"], rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		// check that total instance count is "size"
		log.Printf("[DEBUG] XXXXXXXXXXXXXXXXXXXXXXXX Manager Test: %#v", manager.TargetSize)
		if manager.CurrentSize != size {
			return fmt.Errorf("instance count incorrect")
		}

		// check that at least one instance exists in "targetpool"
		tp, ok := s.RootModule().Resources[targetPool]
		if !ok {
			return fmt.Errorf("Not found: %s", targetPool)
		}

		if tp.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		log.Printf("[DEBUG] XXXXXXXXXXXXXXXXXXXXXXXX Manager Test: %#v", tp)

		targetpool, err := config.clientCompute.TargetPools.Get(
			config.Project, config.Region, tp.Primary.ID).Do()
		if err != nil {
			return err
		}

		// check that total instance count is "size"
		log.Printf("[DEBUG] XXXXXXXXXXXXXXXXXXXXXXXX Manager Test: %#v", len(targetpool.Instances))
		if len(targetpool.Instances) == 0 {
			return fmt.Errorf("no instance in new targetpool")
		}

		// check that the instance template updated
		instanceTemplate, err := config.clientCompute.InstanceTemplates.Get(
			config.Project, template).Do()
		if err != nil {
			return fmt.Errorf("Error reading instance template: %s", err)
		}

		if instanceTemplate.Name != template {
			return fmt.Errorf("instance template not updated")
		}

		return nil
	}
}

const testAccInstanceGroupManager_basic = `
resource "google_compute_instance_template" "foobar" {
	name = "terraform-test-foobar"
	machine_type = "n1-standard-1"
	can_ip_forward = false
	tags = ["foo", "bar"]

	disk {
		source_image = "projects/debian-cloud/global/images/debian-7-wheezy-v20140814"
		auto_delete = true
		boot = true
	}

	network_interface {
		network = "default"
	}

	metadata {
		foo = "bar"
	}

	service_account {
		scopes = ["userinfo-email", "compute-ro", "storage-ro"]
	}
}

resource "google_compute_target_pool" "foobar" {
	description = "Resource created for Terraform acceptance testing"
	name = "terraform-test-foobar"
	session_affinity = "CLIENT_IP_PROTO"
}

resource "google_replicapool_instance_group_manager" "foobar" {
	description = "Terraform test instance group manager"
	name = "terraform-test"
	instance_template = "${google_compute_instance_template.foobar.self_link}"
	target_pools = ["${google_compute_target_pool.foobar.self_link}"]
	base_instance_name = "foobar"
	zone = "us-central1-a"
	size = 2
}`

const testAccInstanceGroupManager_update = `
resource "google_compute_instance_template" "foobar" {
	name = "terraform-test-foobar"
	machine_type = "n1-standard-1"
	can_ip_forward = false
	tags = ["foo", "bar"]

	disk {
		source_image = "projects/debian-cloud/global/images/debian-7-wheezy-v20140814"
		auto_delete = true
		boot = true
	}

	network_interface {
		network = "default"
	}

	metadata {
		foo = "bar"
	}

	service_account {
		scopes = ["userinfo-email", "compute-ro", "storage-ro"]
	}
}

resource "google_compute_instance_template" "foobaz" {
	name = "terraform-test-foobaz"
	machine_type = "n1-standard-1"
	can_ip_forward = false
	tags = ["foo", "bar"]

	disk {
		source_image = "projects/debian-cloud/global/images/debian-7-wheezy-v20140814"
		auto_delete = true
		boot = true
	}

	network_interface {
		network = "default"
	}

	metadata {
		foo = "bar"
	}

	service_account {
		scopes = ["userinfo-email", "compute-ro", "storage-ro"]
	}
}

resource "google_compute_target_pool" "foobar" {
	description = "Resource created for Terraform acceptance testing"
	name = "terraform-test-foobar"
	session_affinity = "CLIENT_IP_PROTO"
}

resource "google_compute_target_pool" "foobaz" {
	description = "Resource created for Terraform acceptance testing"
	name = "terraform-test-foobaz"
	session_affinity = "CLIENT_IP_PROTO"
}

resource "google_replicapool_instance_group_manager" "foobar" {
	description = "Terraform test instance group manager"
	name = "terraform-test"
	instance_template = "${google_compute_instance_template.foobar.self_link}"
	target_pools = ["${google_compute_target_pool.foobaz.self_link}"]
	base_instance_name = "foobar"
	zone = "us-central1-a"
	size = 2
}`

const testAccInstanceGroupManager_update2 = `
resource "google_compute_instance_template" "foobar" {
	name = "terraform-test-foobar"
	machine_type = "n1-standard-1"
	can_ip_forward = false
	tags = ["foo", "bar"]

	disk {
		source_image = "projects/debian-cloud/global/images/debian-7-wheezy-v20140814"
		auto_delete = true
		boot = true
	}

	network_interface {
		network = "default"
	}

	metadata {
		foo = "bar"
	}

	service_account {
		scopes = ["userinfo-email", "compute-ro", "storage-ro"]
	}
}

resource "google_compute_instance_template" "foobaz" {
	name = "terraform-test-foobaz"
	machine_type = "n1-standard-1"
	can_ip_forward = false
	tags = ["foo", "bar"]

	disk {
		source_image = "projects/debian-cloud/global/images/debian-7-wheezy-v20140814"
		auto_delete = true
		boot = true
	}

	network_interface {
		network = "default"
	}

	metadata {
		foo = "bar"
	}

	service_account {
		scopes = ["userinfo-email", "compute-ro", "storage-ro"]
	}
}

resource "google_compute_target_pool" "foobar" {
	description = "Resource created for Terraform acceptance testing"
	name = "terraform-test-foobar"
	session_affinity = "CLIENT_IP_PROTO"
}

resource "google_compute_target_pool" "foobaz" {
	description = "Resource created for Terraform acceptance testing"
	name = "terraform-test-foobaz"
	session_affinity = "CLIENT_IP_PROTO"
}

resource "google_replicapool_instance_group_manager" "foobar" {
	description = "Terraform test instance group manager"
	name = "terraform-test"
	instance_template = "${google_compute_instance_template.foobaz.self_link}"
	target_pools = ["${google_compute_target_pool.foobaz.self_link}"]
	base_instance_name = "foobar"
	zone = "us-central1-a"
	size = 3
}`
