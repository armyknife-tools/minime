package librato

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/henrikhodne/go-librato/librato"
)

func TestAccLibratoMetric_Basic(t *testing.T) {
	var metric librato.Metric
	name := fmt.Sprintf("tftest-metric-%s", acctest.RandString(10))
	typ := "counter"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLibratoMetricDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLibratoMetricConfig(name, typ),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLibratoMetricExists("librato_metric.foobar", &metric),
					testAccCheckLibratoMetricName(&metric, name),
					testAccCheckLibratoMetricType(&metric, []string{"gauge", "counter", "composite"}),
					resource.TestCheckResourceAttr(
						"librato_metric.foobar", "name", name),
				),
			},
		},
	})
}

func testAccCheckLibratoMetricDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*librato.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "librato_metric" {
			continue
		}

		_, _, err := client.Metrics.Get(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Metric still exists")
		}
	}

	return nil
}

func testAccCheckLibratoMetricConfig(name, typ string) string {
	return strings.TrimSpace(fmt.Sprintf(`
    resource "librato_metric" "foobar" {
        name = "%s"
        type = "%s"
        description = "A test composite metric"
        composite = "s(\"librato.cpu.percent.user\", {\"environment\" : \"prod\", \"service\": \"api\"})"
        attributes {
          display_stacked = true,
          created_by_ua = "go-librato/0.1"
        }
    }`, name, typ))
}

func testAccCheckLibratoMetricName(metric *librato.Metric, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if metric.Name == nil || *metric.Name != name {
			return fmt.Errorf("Bad name: %s", *metric.Name)
		}

		return nil
	}
}

func testAccCheckLibratoMetricType(metric *librato.Metric, validTypes []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		m := make(map[string]bool)
		for _, v := range validTypes {
			m[v] = true
		}

		if !m[*metric.Type] {
			return fmt.Errorf("Bad metric type: %s", *metric.Type)
		}

		return nil
	}
}

func testAccCheckLibratoMetricExists(n string, metric *librato.Metric) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Metric ID is set")
		}

		client := testAccProvider.Meta().(*librato.Client)

		foundMetric, _, err := client.Metrics.Get(rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundMetric.Name == nil || *foundMetric.Name != rs.Primary.ID {
			return fmt.Errorf("Metric not found")
		}

		*metric = *foundMetric

		return nil
	}
}
