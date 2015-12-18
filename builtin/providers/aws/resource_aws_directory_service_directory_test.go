package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/directoryservice"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAWSDirectoryServiceDirectory_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDirectoryServiceDirectoryDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDirectoryServiceDirectoryConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceDirectoryExists("aws_directory_service_directory.bar"),
				),
			},
		},
	})
}

func TestAccAWSDirectoryServiceDirectory_microsoft(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDirectoryServiceDirectoryDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDirectoryServiceDirectoryConfig_microsoft,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceDirectoryExists("aws_directory_service_directory.bar"),
				),
			},
		},
	})
}

func TestAccAWSDirectoryServiceDirectory_withAliasAndSso(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDirectoryServiceDirectoryDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDirectoryServiceDirectoryConfig_withAlias,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceDirectoryExists("aws_directory_service_directory.bar_a"),
					testAccCheckServiceDirectoryAlias("aws_directory_service_directory.bar_a",
						fmt.Sprintf("tf-d-%d", randomInteger)),
					testAccCheckServiceDirectorySso("aws_directory_service_directory.bar_a", false),
				),
			},
			resource.TestStep{
				Config: testAccDirectoryServiceDirectoryConfig_withSso,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceDirectoryExists("aws_directory_service_directory.bar_a"),
					testAccCheckServiceDirectoryAlias("aws_directory_service_directory.bar_a",
						fmt.Sprintf("tf-d-%d", randomInteger)),
					testAccCheckServiceDirectorySso("aws_directory_service_directory.bar_a", true),
				),
			},
			resource.TestStep{
				Config: testAccDirectoryServiceDirectoryConfig_withSso_modified,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceDirectoryExists("aws_directory_service_directory.bar_a"),
					testAccCheckServiceDirectoryAlias("aws_directory_service_directory.bar_a",
						fmt.Sprintf("tf-d-%d", randomInteger)),
					testAccCheckServiceDirectorySso("aws_directory_service_directory.bar_a", false),
				),
			},
		},
	})
}

func testAccCheckDirectoryServiceDirectoryDestroy(s *terraform.State) error {
	if len(s.RootModule().Resources) > 0 {
		return fmt.Errorf("Expected all resources to be gone, but found: %#v",
			s.RootModule().Resources)
	}

	return nil
}

func testAccCheckServiceDirectoryExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		dsconn := testAccProvider.Meta().(*AWSClient).dsconn
		out, err := dsconn.DescribeDirectories(&directoryservice.DescribeDirectoriesInput{
			DirectoryIds: []*string{aws.String(rs.Primary.ID)},
		})

		if err != nil {
			return err
		}

		if len(out.DirectoryDescriptions) < 1 {
			return fmt.Errorf("No DS directory found")
		}

		if *out.DirectoryDescriptions[0].DirectoryId != rs.Primary.ID {
			return fmt.Errorf("DS directory ID mismatch - existing: %q, state: %q",
				*out.DirectoryDescriptions[0].DirectoryId, rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckServiceDirectoryAlias(name, alias string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		dsconn := testAccProvider.Meta().(*AWSClient).dsconn
		out, err := dsconn.DescribeDirectories(&directoryservice.DescribeDirectoriesInput{
			DirectoryIds: []*string{aws.String(rs.Primary.ID)},
		})

		if err != nil {
			return err
		}

		if *out.DirectoryDescriptions[0].Alias != alias {
			return fmt.Errorf("DS directory Alias mismatch - actual: %q, expected: %q",
				*out.DirectoryDescriptions[0].Alias, alias)
		}

		return nil
	}
}

func testAccCheckServiceDirectorySso(name string, ssoEnabled bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		dsconn := testAccProvider.Meta().(*AWSClient).dsconn
		out, err := dsconn.DescribeDirectories(&directoryservice.DescribeDirectoriesInput{
			DirectoryIds: []*string{aws.String(rs.Primary.ID)},
		})

		if err != nil {
			return err
		}

		if *out.DirectoryDescriptions[0].SsoEnabled != ssoEnabled {
			return fmt.Errorf("DS directory SSO mismatch - actual: %t, expected: %t",
				*out.DirectoryDescriptions[0].SsoEnabled, ssoEnabled)
		}

		return nil
	}
}

const testAccDirectoryServiceDirectoryConfig = `
resource "aws_directory_service_directory" "bar" {
  name = "corp.notexample.com"
  password = "SuperSecretPassw0rd"
  size = "Small"

  vpc_settings {
    vpc_id = "${aws_vpc.main.id}"
    subnet_ids = ["${aws_subnet.foo.id}", "${aws_subnet.bar.id}"]
  }
}

resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "foo" {
  vpc_id = "${aws_vpc.main.id}"
  availability_zone = "us-west-2a"
  cidr_block = "10.0.1.0/24"
}
resource "aws_subnet" "bar" {
  vpc_id = "${aws_vpc.main.id}"
  availability_zone = "us-west-2b"
  cidr_block = "10.0.2.0/24"
}
`

const testAccDirectoryServiceDirectoryConfig_microsoft = `
resource "aws_directory_service_directory" "bar" {
  name = "corp.notexample.com"
  password = "SuperSecretPassw0rd"
  type = "MicrosoftAD"

  vpc_settings {
    vpc_id = "${aws_vpc.main.id}"
    subnet_ids = ["${aws_subnet.foo.id}", "${aws_subnet.bar.id}"]
  }
}

resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "foo" {
  vpc_id = "${aws_vpc.main.id}"
  availability_zone = "us-west-2a"
  cidr_block = "10.0.1.0/24"
}
resource "aws_subnet" "bar" {
  vpc_id = "${aws_vpc.main.id}"
  availability_zone = "us-west-2b"
  cidr_block = "10.0.2.0/24"
}
`

var randomInteger = genRandInt()
var testAccDirectoryServiceDirectoryConfig_withAlias = fmt.Sprintf(`
resource "aws_directory_service_directory" "bar_a" {
  name = "corp.notexample.com"
  password = "SuperSecretPassw0rd"
  size = "Small"
  alias = "tf-d-%d"

  vpc_settings {
    vpc_id = "${aws_vpc.main.id}"
    subnet_ids = ["${aws_subnet.foo.id}", "${aws_subnet.bar.id}"]
  }
}

resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "foo" {
  vpc_id = "${aws_vpc.main.id}"
  availability_zone = "us-west-2a"
  cidr_block = "10.0.1.0/24"
}
resource "aws_subnet" "bar" {
  vpc_id = "${aws_vpc.main.id}"
  availability_zone = "us-west-2b"
  cidr_block = "10.0.2.0/24"
}
`, randomInteger)

var testAccDirectoryServiceDirectoryConfig_withSso = fmt.Sprintf(`
resource "aws_directory_service_directory" "bar_a" {
  name = "corp.notexample.com"
  password = "SuperSecretPassw0rd"
  size = "Small"
  alias = "tf-d-%d"
  enable_sso = true

  vpc_settings {
    vpc_id = "${aws_vpc.main.id}"
    subnet_ids = ["${aws_subnet.foo.id}", "${aws_subnet.bar.id}"]
  }
}

resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "foo" {
  vpc_id = "${aws_vpc.main.id}"
  availability_zone = "us-west-2a"
  cidr_block = "10.0.1.0/24"
}
resource "aws_subnet" "bar" {
  vpc_id = "${aws_vpc.main.id}"
  availability_zone = "us-west-2b"
  cidr_block = "10.0.2.0/24"
}
`, randomInteger)

var testAccDirectoryServiceDirectoryConfig_withSso_modified = fmt.Sprintf(`
resource "aws_directory_service_directory" "bar_a" {
  name = "corp.notexample.com"
  password = "SuperSecretPassw0rd"
  size = "Small"
  alias = "tf-d-%d"
  enable_sso = false

  vpc_settings {
    vpc_id = "${aws_vpc.main.id}"
    subnet_ids = ["${aws_subnet.foo.id}", "${aws_subnet.bar.id}"]
  }
}

resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "foo" {
  vpc_id = "${aws_vpc.main.id}"
  availability_zone = "us-west-2a"
  cidr_block = "10.0.1.0/24"
}
resource "aws_subnet" "bar" {
  vpc_id = "${aws_vpc.main.id}"
  availability_zone = "us-west-2b"
  cidr_block = "10.0.2.0/24"
}
`, randomInteger)
