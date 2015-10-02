package aws

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestDiffTagsKinesis(t *testing.T) {
	cases := []struct {
		Old, New       map[string]interface{}
		Create, Remove map[string]string
	}{
		// Basic add/remove
		{
			Old: map[string]interface{}{
				"foo": "bar",
			},
			New: map[string]interface{}{
				"bar": "baz",
			},
			Create: map[string]string{
				"bar": "baz",
			},
			Remove: map[string]string{
				"foo": "bar",
			},
		},

		// Modify
		{
			Old: map[string]interface{}{
				"foo": "bar",
			},
			New: map[string]interface{}{
				"foo": "baz",
			},
			Create: map[string]string{
				"foo": "baz",
			},
			Remove: map[string]string{
				"foo": "bar",
			},
		},
	}

	for i, tc := range cases {
		c, r := diffTagsKinesis(tagsFromMapKinesis(tc.Old), tagsFromMapKinesis(tc.New))
		cm := tagsToMapKinesis(c)
		rm := tagsToMapKinesis(r)
		if !reflect.DeepEqual(cm, tc.Create) {
			t.Fatalf("%d: bad create: %#v", i, cm)
		}
		if !reflect.DeepEqual(rm, tc.Remove) {
			t.Fatalf("%d: bad remove: %#v", i, rm)
		}
	}
}

// testAccCheckTags can be used to check the tags on a resource.
func testAccCheckKinesisTags(ts []*kinesis.Tag, key string, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		m := tagsToMapKinesis(ts)
		v, ok := m[key]
		if value != "" && !ok {
			return fmt.Errorf("Missing tag: %s", key)
		} else if value == "" && ok {
			return fmt.Errorf("Extra tag: %s", key)
		}
		if value == "" {
			return nil
		}

		if v != value {
			return fmt.Errorf("%s: bad value: %s", key, v)
		}

		return nil
	}
}
