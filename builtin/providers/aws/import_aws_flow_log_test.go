package aws

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAWSFlowLog_importBasic(t *testing.T) {
	resourceName := "aws_flow_log.test_flow_log"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFlowLogDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccFlowLogConfig_basic,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
