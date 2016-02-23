package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAWSDBClusterParameterGroup_basic(t *testing.T) {
	var v rds.DBClusterParameterGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDBClusterParameterGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAWSDBClusterParameterGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBClusterParameterGroupExists("aws_rds_cluster_parameter_group.bar", &v),
					testAccCheckAWSDBClusterParameterGroupAttributes(&v),
					resource.TestCheckResourceAttr(
						"aws_rds_cluster_parameter_group.bar", "name", "cluster-parameter-group-test-terraform"),
					resource.TestCheckResourceAttr(
						"aws_rds_cluster_parameter_group.bar", "family", "aurora5.6"),
					resource.TestCheckResourceAttr(
						"aws_rds_cluster_parameter_group.bar", "description", "Test cluster parameter group for terraform"),
					resource.TestCheckResourceAttr(
						"aws_rds_cluster_parameter_group.bar", "parameter.1708034931.name", "character_set_results"),
					resource.TestCheckResourceAttr(
						"aws_rds_cluster_parameter_group.bar", "parameter.1708034931.value", "utf8"),
					resource.TestCheckResourceAttr(
						"aws_rds_cluster_parameter_group.bar", "parameter.2421266705.name", "character_set_server"),
					resource.TestCheckResourceAttr(
						"aws_rds_cluster_parameter_group.bar", "parameter.2421266705.value", "utf8"),
					resource.TestCheckResourceAttr(
						"aws_rds_cluster_parameter_group.bar", "parameter.2478663599.name", "character_set_client"),
					resource.TestCheckResourceAttr(
						"aws_rds_cluster_parameter_group.bar", "parameter.2478663599.value", "utf8"),
					resource.TestCheckResourceAttr(
						"aws_rds_cluster_parameter_group.bar", "tags.#", "1"),
				),
			},
			resource.TestStep{
				Config: testAccAWSDBClusterParameterGroupAddParametersConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBClusterParameterGroupExists("aws_rds_cluster_parameter_group.bar", &v),
					testAccCheckAWSDBClusterParameterGroupAttributes(&v),
					resource.TestCheckResourceAttr(
						"aws_rds_cluster_parameter_group.bar", "name", "cluster-parameter-group-test-terraform"),
					resource.TestCheckResourceAttr(
						"aws_rds_cluster_parameter_group.bar", "family", "aurora5.6"),
					resource.TestCheckResourceAttr(
						"aws_rds_cluster_parameter_group.bar", "description", "Test cluster parameter group for terraform"),
					resource.TestCheckResourceAttr(
						"aws_rds_cluster_parameter_group.bar", "parameter.1706463059.name", "collation_connection"),
					resource.TestCheckResourceAttr(
						"aws_rds_cluster_parameter_group.bar", "parameter.1706463059.value", "utf8_unicode_ci"),
					resource.TestCheckResourceAttr(
						"aws_rds_cluster_parameter_group.bar", "parameter.1708034931.name", "character_set_results"),
					resource.TestCheckResourceAttr(
						"aws_rds_cluster_parameter_group.bar", "parameter.1708034931.value", "utf8"),
					resource.TestCheckResourceAttr(
						"aws_rds_cluster_parameter_group.bar", "parameter.2421266705.name", "character_set_server"),
					resource.TestCheckResourceAttr(
						"aws_rds_cluster_parameter_group.bar", "parameter.2421266705.value", "utf8"),
					resource.TestCheckResourceAttr(
						"aws_rds_cluster_parameter_group.bar", "parameter.2475805061.name", "collation_server"),
					resource.TestCheckResourceAttr(
						"aws_rds_cluster_parameter_group.bar", "parameter.2475805061.value", "utf8_unicode_ci"),
					resource.TestCheckResourceAttr(
						"aws_rds_cluster_parameter_group.bar", "parameter.2478663599.name", "character_set_client"),
					resource.TestCheckResourceAttr(
						"aws_rds_cluster_parameter_group.bar", "parameter.2478663599.value", "utf8"),
					resource.TestCheckResourceAttr(
						"aws_rds_cluster_parameter_group.bar", "tags.#", "2"),
				),
			},
		},
	})
}

func TestAccAWSDBClusterParameterGroupOnly(t *testing.T) {
	var v rds.DBClusterParameterGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDBClusterParameterGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAWSDBClusterParameterGroupOnlyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBClusterParameterGroupExists("aws_rds_cluster_parameter_group.bar", &v),
					testAccCheckAWSDBClusterParameterGroupAttributes(&v),
					resource.TestCheckResourceAttr(
						"aws_rds_cluster_parameter_group.bar", "name", "cluster-parameter-group-test-terraform"),
					resource.TestCheckResourceAttr(
						"aws_rds_cluster_parameter_group.bar", "family", "aurora5.6"),
					resource.TestCheckResourceAttr(
						"aws_rds_cluster_parameter_group.bar", "description", "Test cluster parameter group for terraform"),
				),
			},
		},
	})
}

func TestResourceAWSDBClusterParameterGroupName_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "tEsting123",
			ErrCount: 1,
		},
		{
			Value:    "testing123!",
			ErrCount: 1,
		},
		{
			Value:    "1testing123",
			ErrCount: 1,
		},
		{
			Value:    "testing--123",
			ErrCount: 1,
		},
		{
			Value:    "testing123-",
			ErrCount: 1,
		},
		{
			Value:    randomString(256),
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateDbParamGroupName(tc.Value, "aws_rds_cluster_parameter_group_name")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the DB Cluster Parameter Group Name to trigger a validation error")
		}
	}
}

func testAccCheckAWSDBClusterParameterGroupDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).rdsconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_rds_cluster_parameter_group" {
			continue
		}

		// Try to find the Group
		resp, err := conn.DescribeDBClusterParameterGroups(
			&rds.DescribeDBClusterParameterGroupsInput{
				DBClusterParameterGroupName: aws.String(rs.Primary.ID),
			})

		if err == nil {
			if len(resp.DBClusterParameterGroups) != 0 &&
				*resp.DBClusterParameterGroups[0].DBClusterParameterGroupName == rs.Primary.ID {
				return fmt.Errorf("DB Cluster Parameter Group still exists")
			}
		}

		// Verify the error
		newerr, ok := err.(awserr.Error)
		if !ok {
			return err
		}
		if newerr.Code() != "DBParameterGroupNotFound" {
			return err
		}
	}

	return nil
}

func testAccCheckAWSDBClusterParameterGroupAttributes(v *rds.DBClusterParameterGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if *v.DBClusterParameterGroupName != "cluster-parameter-group-test-terraform" {
			return fmt.Errorf("bad name: %#v", v.DBClusterParameterGroupName)
		}

		if *v.DBParameterGroupFamily != "aurora5.6" {
			return fmt.Errorf("bad family: %#v", v.DBParameterGroupFamily)
		}

		if *v.Description != "Test cluster parameter group for terraform" {
			return fmt.Errorf("bad description: %#v", v.Description)
		}

		return nil
	}
}

func testAccCheckAWSDBClusterParameterGroupExists(n string, v *rds.DBClusterParameterGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DB Cluster Parameter Group ID is set")
		}

		conn := testAccProvider.Meta().(*AWSClient).rdsconn

		opts := rds.DescribeDBClusterParameterGroupsInput{
			DBClusterParameterGroupName: aws.String(rs.Primary.ID),
		}

		resp, err := conn.DescribeDBClusterParameterGroups(&opts)

		if err != nil {
			return err
		}

		if len(resp.DBClusterParameterGroups) != 1 ||
			*resp.DBClusterParameterGroups[0].DBClusterParameterGroupName != rs.Primary.ID {
			return fmt.Errorf("DB Cluster Parameter Group not found")
		}

		*v = *resp.DBClusterParameterGroups[0]

		return nil
	}
}

const testAccAWSDBClusterParameterGroupConfig = `
resource "aws_rds_cluster_parameter_group" "bar" {
	name = "cluster-parameter-group-test-terraform"
	family = "aurora5.6"
	description = "Test cluster parameter group for terraform"
	parameter {
	  name = "character_set_server"
	  value = "utf8"
	}
	parameter {
	  name = "character_set_client"
	  value = "utf8"
	}
	parameter{
	  name = "character_set_results"
	  value = "utf8"
	}
	tags {
		foo = "bar"
	}
}
`

const testAccAWSDBClusterParameterGroupAddParametersConfig = `
resource "aws_rds_cluster_parameter_group" "bar" {
	name = "cluster-parameter-group-test-terraform"
	family = "aurora5.6"
	description = "Test cluster parameter group for terraform"
	parameter {
	  name = "character_set_server"
	  value = "utf8"
	}
	parameter {
	  name = "character_set_client"
	  value = "utf8"
	}
	parameter{
	  name = "character_set_results"
	  value = "utf8"
	}
	parameter {
	  name = "collation_server"
	  value = "utf8_unicode_ci"
	}
	parameter {
	  name = "collation_connection"
	  value = "utf8_unicode_ci"
	}
	tags {
		foo = "bar"
		baz = "foo"
	}
}
`

const testAccAWSDBClusterParameterGroupOnlyConfig = `
resource "aws_rds_cluster_parameter_group" "bar" {
	name = "cluster-parameter-group-test-terraform"
	family = "aurora5.6"
	description = "Test cluster parameter group for terraform"
}
`
