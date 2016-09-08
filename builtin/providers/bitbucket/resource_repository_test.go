package bitbucket

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccBitbucketRepository_basic(t *testing.T) {
	var repo Repository

	testUser := os.Getenv("BITBUCKET_USERNAME")
	testAccBitbucketRepositoryConfig := fmt.Sprintf(`
		resource "bitbucket_repository" "test_repo" {
			owner = "%s"
			name = "%s"
		}
	`, testUser, testRepo)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketRepositoryDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccBitbucketRepositoryConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketRepositoryExists("bitbucket_repository.test_repo", &repo),
				),
			},
		},
	})
}

func testAccCheckBitbucketRepositoryDestroy(s *terraform.State) error {
	return nil
}

func testAccCheckBitbucketRepositoryExists(n string, repository *Repository) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No repository ID is set")
		}
		return nil
	}
}
