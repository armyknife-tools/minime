package google

import (
	"fmt"
	"testing"

	"code.google.com/p/google-api-go-client/autoscaler/v1beta2"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAutoscaler_basic(t *testing.T) {
	var ascaler autoscaler.Autoscaler

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAutoscalerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAutoscaler_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoscalerExists(
						"google_autoscaler.foobar", &ascaler),
				),
			},
		},
	})
}

func TestAccAutoscaler_update(t *testing.T) {
	var ascaler autoscaler.Autoscaler

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAutoscalerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAutoscaler_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoscalerExists(
						"google_autoscaler.foobar", &ascaler),
				),
			},
			resource.TestStep{
				Config: testAccAutoscaler_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoscalerExists(
						"google_autoscaler.foobar", &ascaler),
					testAccCheckAutoscalerUpdated(
						"google_autoscaler.foobar", 10),
				),
			},
		},
	})
}

func testAccCheckAutoscalerDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_autoscaler" {
			continue
		}

		_, err := config.clientAutoscaler.Autoscalers.Get(
			config.Project, rs.Primary.Attributes["zone"], rs.Primary.ID).Do()
		if err == nil {
			return fmt.Errorf("Autoscaler still exists")
		}
	}

	return nil
}

func testAccCheckAutoscalerExists(n string, ascaler *autoscaler.Autoscaler) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.clientAutoscaler.Autoscalers.Get(
			config.Project, rs.Primary.Attributes["zone"], rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("Autoscaler not found")
		}

		*ascaler = *found

		return nil
	}
}

func testAccCheckAutoscalerUpdated(n string, max int64) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		ascaler, err := config.clientAutoscaler.Autoscalers.Get(
			config.Project, rs.Primary.Attributes["zone"], rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		if ascaler.AutoscalingPolicy.MaxNumReplicas != max {
			return fmt.Errorf("maximum replicas incorrect")
		}

		return nil
	}
}

const testAccAutoscaler_basic = `
resource "google_compute_instance_template" "foobar" {
	name = "terraform-test-template-foobar"
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
	name = "terraform-test-tpool-foobar"
	session_affinity = "CLIENT_IP_PROTO"
}

resource "google_replicapool_instance_group_manager" "foobar" {
	description = "Terraform test instance group manager"
	name = "terraform-test-groupmanager"
	instance_template = "${google_compute_instance_template.foobar.self_link}"
	target_pools = ["${google_compute_target_pool.foobar.self_link}"]
	base_instance_name = "foobar"
	zone = "us-central1-a"
	size = 0
}

resource "google_autoscaler" "foobar" {
	description = "Resource created for Terraform acceptance testing"
	name = "terraform-test-ascaler"
	zone = "us-central1-a"
	target = "${google_replicapool_instance_group_manager.foobar.self_link}"
	autoscaling_policy = {
		max_replicas = 5
		min_replicas = 0
		cooldown_period = 60
		cpu_utilization = {
			target = 0.5
		}
	}

}`

const testAccAutoscaler_update = `
resource "google_compute_instance_template" "foobar" {
	name = "terraform-test-template-foobar"
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
	name = "terraform-test-tpool-foobar"
	session_affinity = "CLIENT_IP_PROTO"
}

resource "google_replicapool_instance_group_manager" "foobar" {
	description = "Terraform test instance group manager"
	name = "terraform-test-groupmanager"
	instance_template = "${google_compute_instance_template.foobar.self_link}"
	target_pools = ["${google_compute_target_pool.foobar.self_link}"]
	base_instance_name = "foobar"
	zone = "us-central1-a"
	size = 0
}

resource "google_autoscaler" "foobar" {
	description = "Resource created for Terraform acceptance testing"
	name = "terraform-test-ascaler"
	zone = "us-central1-a"
	target = "${google_replicapool_instance_group_manager.foobar.self_link}"
	autoscaling_policy = {
		max_replicas = 10
		min_replicas = 0
		cooldown_period = 60
		cpu_utilization = {
			target = 0.5
		}
	}

}`
