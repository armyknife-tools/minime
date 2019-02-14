package test

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

// an empty config should be ok, because no deprecated/removed fields are set.
func TestResourceList_changed(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: strings.TrimSpace(`
resource "test_resource_list" "foo" {
	list_block {
		string = "a"
		int = 1
	}
}
				`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"test_resource_list.foo", "list_block.#", "1",
					),
					resource.TestCheckResourceAttr(
						"test_resource_list.foo", "list_block.0.string", "a",
					),
					resource.TestCheckResourceAttr(
						"test_resource_list.foo", "list_block.0.int", "1",
					),
				),
			},
			resource.TestStep{
				Config: strings.TrimSpace(`
resource "test_resource_list" "foo" {
	list_block {
		string = "a"
		int = 1
	}

	list_block {
		string = "b"
		int = 2
	}
}
				`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"test_resource_list.foo", "list_block.#", "2",
					),
					resource.TestCheckResourceAttr(
						"test_resource_list.foo", "list_block.0.string", "a",
					),
					resource.TestCheckResourceAttr(
						"test_resource_list.foo", "list_block.0.int", "1",
					),
					resource.TestCheckResourceAttr(
						"test_resource_list.foo", "list_block.1.string", "b",
					),
					resource.TestCheckResourceAttr(
						"test_resource_list.foo", "list_block.1.int", "2",
					),
				),
			},
			resource.TestStep{
				Config: strings.TrimSpace(`
resource "test_resource_list" "foo" {
	list_block {
		string = "a"
		int = 1
	}

	list_block {
		string = "c"
		int = 2
	}
}
				`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"test_resource_list.foo", "list_block.#", "2",
					),
					resource.TestCheckResourceAttr(
						"test_resource_list.foo", "list_block.0.string", "a",
					),
					resource.TestCheckResourceAttr(
						"test_resource_list.foo", "list_block.0.int", "1",
					),
					resource.TestCheckResourceAttr(
						"test_resource_list.foo", "list_block.1.string", "c",
					),
					resource.TestCheckResourceAttr(
						"test_resource_list.foo", "list_block.1.int", "2",
					),
				),
			},
		},
	})
}

func TestResourceList_sublist(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: strings.TrimSpace(`
resource "test_resource_list" "foo" {
	list_block {
		sublist_block {
			string = "a"
			int = 1
		}
	}
}
				`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"test_resource_list.foo", "list_block.0.sublist_block.#", "1",
					),
					resource.TestCheckResourceAttr(
						"test_resource_list.foo", "list_block.0.sublist_block.0.string", "a",
					),
					resource.TestCheckResourceAttr(
						"test_resource_list.foo", "list_block.0.sublist_block.0.int", "1",
					),
				),
			},
		},
	})
}

func TestResourceList_interpolationChanges(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: strings.TrimSpace(`
resource "test_resource_list" "foo" {
	list_block {
		string = "x"
	}
}
resource "test_resource_list" "bar" {
	list_block {
		string = test_resource_list.foo.id
	}
}
				`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"test_resource_list.foo", "list_block.0.string", "x",
					),
					resource.TestCheckResourceAttr(
						"test_resource_list.bar", "list_block.0.string", "testId",
					),
				),
			},
			resource.TestStep{
				Config: strings.TrimSpace(`
resource "test_resource_list" "baz" {
	list_block {
		string = "x"
		int = 1
	}
}
resource "test_resource_list" "bar" {
	list_block {
		string = test_resource_list.baz.id
		int = 3
	}
}
				`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"test_resource_list.baz", "list_block.0.string", "x",
					),
					resource.TestCheckResourceAttr(
						"test_resource_list.bar", "list_block.0.string", "testId",
					),
				),
			},
		},
	})
}

func TestResourceList_removedForcesNew(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: strings.TrimSpace(`
resource "test_resource_list" "foo" {
	list_block {
		force_new = "ok"
	}
}
				`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"test_resource_list.foo", "list_block.0.force_new", "ok",
					),
				),
			},
			resource.TestStep{
				Config: strings.TrimSpace(`
resource "test_resource_list" "foo" {
}
				`),
				Check: resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestResourceList_emptyStrings(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: strings.TrimSpace(`
resource "test_resource_list" "foo" {
  list_block {
    sublist = ["a", ""]
  }

  list_block {
    sublist = [""]
  }

  list_block {
    sublist = ["", "c", ""]
  }
}
				`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("test_resource_list.foo", "list_block.0.sublist.0", "a"),
					resource.TestCheckResourceAttr("test_resource_list.foo", "list_block.0.sublist.1", ""),
					resource.TestCheckResourceAttr("test_resource_list.foo", "list_block.1.sublist.0", ""),
					resource.TestCheckResourceAttr("test_resource_list.foo", "list_block.2.sublist.0", ""),
					resource.TestCheckResourceAttr("test_resource_list.foo", "list_block.2.sublist.1", "c"),
					resource.TestCheckResourceAttr("test_resource_list.foo", "list_block.2.sublist.2", ""),
				),
			},
			resource.TestStep{
				Config: strings.TrimSpace(`
resource "test_resource_list" "foo" {
  list_block {
    sublist = [""]
  }

  list_block {
    sublist = []
  }

  list_block {
    sublist = ["", "c"]
  }
}
			`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("test_resource_list.foo", "list_block.0.sublist.#", "1"),
					resource.TestCheckResourceAttr("test_resource_list.foo", "list_block.0.sublist.0", ""),
					resource.TestCheckResourceAttr("test_resource_list.foo", "list_block.1.sublist.#", "0"),
					resource.TestCheckResourceAttr("test_resource_list.foo", "list_block.2.sublist.1", "c"),
					resource.TestCheckResourceAttr("test_resource_list.foo", "list_block.2.sublist.#", "2"),
				),
			},
		},
	})
}
