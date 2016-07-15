package aws

import (
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/applicationautoscaling"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAWSAppautoScalingTarget_basic(t *testing.T) {
	var target applicationautoscaling.ScalableTarget
	var awsAccountId = os.Getenv("AWS_ACCOUNT_ID")

	randClusterName := fmt.Sprintf("cluster-%s", acctest.RandString(10))
	randResourceId := fmt.Sprintf("service/%s/%s", randClusterName, acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "aws_appautoscaling_target.bar",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAWSAppautoscalingTargetDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAWSAppautoscalingTargetConfig(randClusterName, randResourceId, awsAccountId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAppautoscalingTargetExists("aws_appautoscaling_target.bar", &target),
					testAccCheckAWSAppautoscalingTargetAttributes(&target, randResourceId),
					resource.TestCheckResourceAttr("aws_appautoscaling_target.bar", "service_namespace", "ecs"),
					resource.TestCheckResourceAttr("aws_appautoscaling_target.bar", "resource_id", fmt.Sprintf("service/%s/foobar", randClusterName)),
					resource.TestCheckResourceAttr("aws_appautoscaling_target.bar", "scalable_dimension", "ecs:service:DesiredCount"),
					resource.TestCheckResourceAttr("aws_appautoscaling_target.bar", "min_capacity", "1"),
					resource.TestCheckResourceAttr("aws_appautoscaling_target.bar", "max_capacity", "3"),
				),
			},

			resource.TestStep{
				Config: testAccAWSAppautoscalingTargetConfigUpdate(randClusterName, randResourceId, awsAccountId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAppautoscalingTargetExists("aws_appautoscaling_target.bar", &target),
					resource.TestCheckResourceAttr("aws_appautoscaling_target.bar", "min_capacity", "3"),
					resource.TestCheckResourceAttr("aws_appautoscaling_target.bar", "max_capacity", "6"),
				),
			},
		},
	})
}

func testAccCheckAWSAppautoscalingTargetDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).appautoscalingconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_appautoscaling_target" {
			continue
		}

		// Try to find the target
		describeTargets, err := conn.DescribeScalableTargets(
			&applicationautoscaling.DescribeScalableTargetsInput{
				ResourceIds: []*string{aws.String(rs.Primary.ID)},
			},
		)

		if err == nil {
			if len(describeTargets.ScalableTargets) != 0 &&
				*describeTargets.ScalableTargets[0].ResourceId == rs.Primary.ID {
				return fmt.Errorf("Application AutoScaling Target still exists")
			}
		}

		// Verify error
		e, ok := err.(awserr.Error)
		if !ok {
			return err
		}
		if e.Code() != "" {
			return e
		}
	}

	return nil
}

func testAccCheckAWSAppautoscalingTargetExists(n string, target *applicationautoscaling.ScalableTarget) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Application AutoScaling Target ID is set")
		}

		conn := testAccProvider.Meta().(*AWSClient).appautoscalingconn

		describeTargets, err := conn.DescribeScalableTargets(
			&applicationautoscaling.DescribeScalableTargetsInput{
				ResourceIds: []*string{aws.String(rs.Primary.ID)},
			},
		)

		if err != nil {
			return err
		}

		if len(describeTargets.ScalableTargets) != 1 ||
			*describeTargets.ScalableTargets[0].ResourceId != rs.Primary.ID {
			return fmt.Errorf("Application AutoScaling ResourceId not found")
		}

		*target = *describeTargets.ScalableTargets[0]

		return nil
	}
}

func testAccCheckAWSAppautoscalingTargetAttributes(target *applicationautoscaling.ScalableTarget, resourceId string) resource.TestCheckFunc {
	return nil
}

func testAccAWSAppautoscalingTargetConfig(
	randClusterName string,
	randResourceId string,
	awsAccountId string) string {
	return fmt.Sprintf(`
resource "aws_ecs_cluster" "foo" {
	name = "%s"
}
resource "aws_ecs_task_definition" "task" {
	family = "foobar"
	container_definitions = <<EOF
[
    {
        "name": "busybox",
        "image": "busybox:latest",
        "cpu": 10,
        "memory": 128,
        "essential": true
    }
]
EOF
}
resource "aws_ecs_service" "service" {
	name = "foobar"
	cluster = "${aws_ecs_cluster.foo.id}"
	task_definition = "${aws_ecs_task_definition.task.arn}"
	desired_count = 1

	deployment_maximum_percent = 200
	deployment_minimum_healthy_percent = 50
}
resource "aws_appautoscaling_target" "bar" {
	service_namespace = "ecs"
	resource_id = "service/${aws_ecs_cluster.foo.name}/${aws_ecs_service.service.name}"
	scalable_dimension = "ecs:service:DesiredCount"
	role_arn = "arn:aws:iam::%s:role/ecsAutoscaleRole"	
	min_capacity = 1
	max_capacity = 4
}
`, randClusterName, awsAccountId)
}

func testAccAWSAppautoscalingTargetConfigUpdate(
	randClusterName,
	randResourceId string,
	awsAccountId string) string {
	return fmt.Sprintf(`
resource "aws_ecs_cluster" "foo" {
	name = "%s"
}
resource "aws_ecs_task_definition" "task" {
	family = "foobar"
	container_definitions = <<EOF
[
    {
        "name": "busybox",
        "image": "busybox:latest",
        "cpu": 10,
        "memory": 128,
        "essential": true
    }
]
EOF
}
resource "aws_ecs_service" "service" {
	name = "foobar"
	cluster = "${aws_ecs_cluster.foo.id}"
	task_definition = "${aws_ecs_task_definition.task.arn}"
	desired_count = 1

	deployment_maximum_percent = 200
	deployment_minimum_healthy_percent = 50
}
resource "aws_appautoscaling_target" "bar" {
	service_namespace = "ecs"
	resource_id = "service/${aws_ecs_cluster.foo.name}/${aws_ecs_service.service.name}"
	scalable_dimension = "ecs:service:DesiredCount"
	role_arn = "arn:aws:iam::%s:role/ecsAutoscaleRole"	
	min_capacity = 2
	max_capacity = 8
}
`, randClusterName, awsAccountId)
}
