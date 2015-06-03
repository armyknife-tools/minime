package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

func TestAccAWSSNSTopic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSSNSTopicDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAWSSNSTopicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSNSTopicExists("aws_sns_topic.test_topic"),
				),
			},
		},
	})
}


func testAccCheckAWSSNSTopicDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).snsconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_sns_topic" {
			continue
		}

		// Check if the topic exists by fetching its attributes
		params := &sns.GetTopicAttributesInput{
			TopicARN: aws.String(rs.Primary.ID),
		}
		_, err := conn.GetTopicAttributes(params)
		if err == nil {
			return fmt.Errorf("Topic exists when it should be destroyed!")
		}

		// Verify the error is an API error, not something else
		_, ok := err.(awserr.Error)
		if !ok {
			return err
		}
	}

	return nil
}


func testAccCheckAWSSNSTopicExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SNS topic with that ARN exists")
		}

		conn := testAccProvider.Meta().(*AWSClient).snsconn

		params := &sns.GetTopicAttributesInput{
			TopicARN: aws.String(rs.Primary.ID),
		}
		_, err := conn.GetTopicAttributes(params)

		if err != nil {
			return err
		}

		return nil
	}
}

const testAccAWSSNSTopicConfig = `
resource "aws_sns_topic" "test_topic" {
    name = "terraform-test-topic"
}
`
