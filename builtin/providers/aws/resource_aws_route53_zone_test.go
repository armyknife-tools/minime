package aws

import (
	"fmt"
	"sort"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/route53"
)

func TestCleanPrefix(t *testing.T) {
	cases := []struct {
		Input, Prefix, Output string
	}{
		{"/hostedzone/foo", "/hostedzone/", "foo"},
		{"/change/foo", "/change/", "foo"},
		{"/bar", "/test", "/bar"},
	}

	for _, tc := range cases {
		actual := cleanPrefix(tc.Input, tc.Prefix)
		if actual != tc.Output {
			t.Fatalf("input: %s\noutput: %s", tc.Input, actual)
		}
	}
}

func TestCleanZoneID(t *testing.T) {
	cases := []struct {
		Input, Output string
	}{
		{"/hostedzone/foo", "foo"},
		{"/change/foo", "/change/foo"},
		{"/bar", "/bar"},
	}

	for _, tc := range cases {
		actual := cleanZoneID(tc.Input)
		if actual != tc.Output {
			t.Fatalf("input: %s\noutput: %s", tc.Input, actual)
		}
	}
}

func TestCleanChangeID(t *testing.T) {
	cases := []struct {
		Input, Output string
	}{
		{"/hostedzone/foo", "/hostedzone/foo"},
		{"/change/foo", "foo"},
		{"/bar", "/bar"},
	}

	for _, tc := range cases {
		actual := cleanChangeID(tc.Input)
		if actual != tc.Output {
			t.Fatalf("input: %s\noutput: %s", tc.Input, actual)
		}
	}
}

func TestAccRoute53Zone(t *testing.T) {
	var zone route53.GetHostedZoneOutput
	var td route53.ResourceTagSet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRoute53ZoneDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRoute53ZoneConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRoute53ZoneExists("aws_route53_zone.main", &zone),
					testAccLoadTagsR53(&zone, &td),
					testAccCheckTagsR53(&td.Tags, "foo", "bar"),
				),
			},
		},
	})
}

func TestAccRoute53PrivateZone(t *testing.T) {
	var zone route53.GetHostedZoneOutput

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRoute53ZoneDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRoute53PrivateZoneConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRoute53ZoneExists("aws_route53_zone.main", &zone),
					testAccCheckRoute53ZoneAssociationExists("aws_vpc.main", &zone),
				),
			},
		},
	})
}

func testAccCheckRoute53ZoneDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).r53conn
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_route53_zone" {
			continue
		}

		_, err := conn.GetHostedZone(&route53.GetHostedZoneInput{ID: aws.String(rs.Primary.ID)})
		if err == nil {
			return fmt.Errorf("Hosted zone still exists")
		}
	}
	return nil
}

func testAccCheckRoute53ZoneExists(n string, zone *route53.GetHostedZoneOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No hosted zone ID is set")
		}

		conn := testAccProvider.Meta().(*AWSClient).r53conn
		resp, err := conn.GetHostedZone(&route53.GetHostedZoneInput{ID: aws.String(rs.Primary.ID)})
		if err != nil {
			return fmt.Errorf("Hosted zone err: %v", err)
		}

		if ! *resp.HostedZone.Config.PrivateZone {
			sorted_ns := make([]string, len(resp.DelegationSet.NameServers))
			for i, ns := range resp.DelegationSet.NameServers {
				sorted_ns[i] = *ns
			}
			sort.Strings(sorted_ns)
			for idx, ns := range sorted_ns {
				attribute := fmt.Sprintf("name_servers.%d", idx)
				dsns := rs.Primary.Attributes[attribute]
				if dsns != ns {
					return fmt.Errorf("Got: %v for %v, Expected: %v", dsns, attribute, ns)
				}
			}
		}

		*zone = *resp
		return nil
	}
}

func testAccCheckRoute53ZoneAssociationExists(n string, zone *route53.GetHostedZoneOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPC ID is set")
		}

		var associatedVPC *route53.VPC
		for _, vpc := range zone.VPCs {
			if *vpc.VPCID == rs.Primary.ID {
				associatedVPC = vpc
			}
		}
		if associatedVPC == nil {
			return fmt.Errorf("VPC: %v is not associated to Zone: %v")
		}
		return nil
	}
}

func testAccLoadTagsR53(zone *route53.GetHostedZoneOutput, td *route53.ResourceTagSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*AWSClient).r53conn

		zone := cleanZoneID(*zone.HostedZone.ID)
		req := &route53.ListTagsForResourceInput{
			ResourceID:   aws.String(zone),
			ResourceType: aws.String("hostedzone"),
		}

		resp, err := conn.ListTagsForResource(req)
		if err != nil {
			return err
		}

		if resp.ResourceTagSet != nil {
			*td = *resp.ResourceTagSet
		}

		return nil
	}
}

const testAccRoute53ZoneConfig = `
resource "aws_route53_zone" "main" {
	name = "hashicorp.com"

	tags {
		foo = "bar"
		Name = "tf-route53-tag-test"
	}
}
`

const testAccRoute53PrivateZoneConfig = `
resource "aws_vpc" "main" {
	cidr_block = "172.29.0.0/24"
	instance_tenancy = "default"
	enable_dns_support = true
	enable_dns_hostnames = true
}

resource "aws_route53_zone" "main" {
	name = "hashicorp.com"
	vpc_id = "${aws_vpc.main.id}"
}
`
