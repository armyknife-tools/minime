package azure

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureHostedServiceBasic(t *testing.T) {
	name := "azure_hosted_service.foo"

	hostedServiceName := fmt.Sprintf("terraform-testing-service%d", acctest.RandInt())
	config := fmt.Sprintf(testAccAzureHostedServiceBasic, hostedServiceName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAzureHostedServiceDestroyed,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAzureHostedServiceExists(name),
					resource.TestCheckResourceAttr(name, "name", hostedServiceName),
					resource.TestCheckResourceAttr(name, "location", "North Europe"),
					resource.TestCheckResourceAttr(name, "ephemeral_contents", "false"),
					resource.TestCheckResourceAttr(name, "description", "very discriptive"),
					resource.TestCheckResourceAttr(name, "label", "very identifiable"),
				),
			},
		},
	})
}

func TestAccAzureHostedServiceUpdate(t *testing.T) {
	name := "azure_hosted_service.foo"

	hostedServiceName := fmt.Sprintf("terraform-testing-service%d", acctest.RandInt())

	basicConfig := fmt.Sprintf(testAccAzureHostedServiceBasic, hostedServiceName)
	updateConfig := fmt.Sprintf(testAccAzureHostedServiceUpdate, hostedServiceName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAzureHostedServiceDestroyed,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: basicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAzureHostedServiceExists(name),
					resource.TestCheckResourceAttr(name, "name", hostedServiceName),
					resource.TestCheckResourceAttr(name, "location", "North Europe"),
					resource.TestCheckResourceAttr(name, "ephemeral_contents", "false"),
					resource.TestCheckResourceAttr(name, "description", "very discriptive"),
					resource.TestCheckResourceAttr(name, "label", "very identifiable"),
				),
			},

			resource.TestStep{
				Config: updateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAzureHostedServiceExists(name),
					resource.TestCheckResourceAttr(name, "name", hostedServiceName),
					resource.TestCheckResourceAttr(name, "location", "North Europe"),
					resource.TestCheckResourceAttr(name, "ephemeral_contents", "true"),
					resource.TestCheckResourceAttr(name, "description", "very discriptive"),
					resource.TestCheckResourceAttr(name, "label", "very identifiable"),
				),
			},
		},
	})
}

func testAccCheckAzureHostedServiceExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Hosted Service resource not found.")
		}

		if resource.Primary.ID == "" {
			return fmt.Errorf("Resource's ID is not set.")
		}

		hostedServiceClient := testAccProvider.Meta().(*Client).hostedServiceClient
		_, err := hostedServiceClient.GetHostedService(resource.Primary.ID)
		return err
	}
}

func testAccCheckAzureHostedServiceDestroyed(s *terraform.State) error {
	hostedServiceClient := testAccProvider.Meta().(*Client).hostedServiceClient

	for _, resource := range s.RootModule().Resources {
		if resource.Type != "azure_hosted_service" {
			continue
		}

		if resource.Primary.ID == "" {
			return fmt.Errorf("No Azure Hosted Service Resource found.")
		}

		_, err := hostedServiceClient.GetHostedService(resource.Primary.ID)

		return testAccResourceDestroyedErrorFilter("Hosted Service", err)
	}

	return nil
}

const testAccAzureHostedServiceBasic = `
resource "azure_hosted_service" "foo" {
	name = "%s"
	location = "North Europe"
    ephemeral_contents = false
	description = "very discriptive"
    label = "very identifiable"
}
`
const testAccAzureHostedServiceUpdate = `
resource "azure_hosted_service" "foo" {
	name = "%s"
	location = "North Europe"
    ephemeral_contents = true
	description = "very discriptive"
    label = "very identifiable"
}
`
