package schema

import (
	"math"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/terraform"
)

func TestResourceDataGet(t *testing.T) {
	cases := []struct {
		Schema map[string]*Schema
		State  *terraform.InstanceState
		Diff   *terraform.InstanceDiff
		Key    string
		Value  interface{}
	}{
		// #0
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: nil,

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						Old:         "foo",
						New:         "bar",
						NewComputed: true,
					},
				},
			},

			Key:   "availability_zone",
			Value: "",
		},

		// #1
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: nil,

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						Old:         "",
						New:         "foo",
						RequiresNew: true,
					},
				},
			},

			Key: "availability_zone",

			Value: "foo",
		},

		// #2
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: nil,

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						Old:      "",
						New:      "foo!",
						NewExtra: "foo",
					},
				},
			},

			Key:   "availability_zone",
			Value: "foo",
		},

		// #3
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"availability_zone": "bar",
				},
			},

			Diff: nil,

			Key: "availability_zone",

			Value: "bar",
		},

		// #4
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"availability_zone": "foo",
				},
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						Old:         "foo",
						New:         "bar",
						NewComputed: true,
					},
				},
			},

			Key:   "availability_zone",
			Value: "",
		},

		// #5
		{
			Schema: map[string]*Schema{
				"port": &Schema{
					Type:     TypeInt,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"port": "80",
				},
			},

			Diff: nil,

			Key: "port",

			Value: 80,
		},

		// #6
		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeList,
					Required: true,
					Elem:     &Schema{Type: TypeInt},
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"ports.#": "3",
					"ports.0": "1",
					"ports.1": "2",
					"ports.2": "5",
				},
			},

			Key: "ports.1",

			Value: 2,
		},

		// #7
		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeList,
					Required: true,
					Elem:     &Schema{Type: TypeInt},
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"ports.#": "3",
					"ports.0": "1",
					"ports.1": "2",
					"ports.2": "5",
				},
			},

			Key: "ports.#",

			Value: 3,
		},

		// #8
		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeList,
					Required: true,
					Elem:     &Schema{Type: TypeInt},
				},
			},

			State: nil,

			Key: "ports.#",

			Value: 0,
		},

		// #9
		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeList,
					Required: true,
					Elem:     &Schema{Type: TypeInt},
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"ports.#": "3",
					"ports.0": "1",
					"ports.1": "2",
					"ports.2": "5",
				},
			},

			Key: "ports",

			Value: []interface{}{1, 2, 5},
		},

		// #10
		{
			Schema: map[string]*Schema{
				"ingress": &Schema{
					Type:     TypeList,
					Required: true,
					Elem: &Resource{
						Schema: map[string]*Schema{
							"from": &Schema{
								Type:     TypeInt,
								Required: true,
							},
						},
					},
				},
			},

			State: nil,

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"ingress.#": &terraform.ResourceAttrDiff{
						Old: "",
						New: "1",
					},
					"ingress.0.from": &terraform.ResourceAttrDiff{
						Old: "",
						New: "8080",
					},
				},
			},

			Key: "ingress.0",

			Value: map[string]interface{}{
				"from": 8080,
			},
		},

		// #11
		{
			Schema: map[string]*Schema{
				"ingress": &Schema{
					Type:     TypeList,
					Required: true,
					Elem: &Resource{
						Schema: map[string]*Schema{
							"from": &Schema{
								Type:     TypeInt,
								Required: true,
							},
						},
					},
				},
			},

			State: nil,

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"ingress.#": &terraform.ResourceAttrDiff{
						Old: "",
						New: "1",
					},
					"ingress.0.from": &terraform.ResourceAttrDiff{
						Old: "",
						New: "8080",
					},
				},
			},

			Key: "ingress",

			Value: []interface{}{
				map[string]interface{}{
					"from": 8080,
				},
			},
		},

		// #12 Computed get
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Computed: true,
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"availability_zone": "foo",
				},
			},

			Key: "availability_zone",

			Value: "foo",
		},

		// #13 Full object
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: nil,

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						Old:         "",
						New:         "foo",
						RequiresNew: true,
					},
				},
			},

			Key: "",

			Value: map[string]interface{}{
				"availability_zone": "foo",
			},
		},

		// #14 List of maps
		{
			Schema: map[string]*Schema{
				"config_vars": &Schema{
					Type:     TypeList,
					Optional: true,
					Computed: true,
					Elem: &Schema{
						Type: TypeMap,
					},
				},
			},

			State: nil,

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"config_vars.#": &terraform.ResourceAttrDiff{
						Old: "0",
						New: "2",
					},
					"config_vars.0.foo": &terraform.ResourceAttrDiff{
						Old: "",
						New: "bar",
					},
					"config_vars.1.bar": &terraform.ResourceAttrDiff{
						Old: "",
						New: "baz",
					},
				},
			},

			Key: "config_vars",

			Value: []interface{}{
				map[string]interface{}{
					"foo": "bar",
				},
				map[string]interface{}{
					"bar": "baz",
				},
			},
		},

		// #15 List of maps in state
		{
			Schema: map[string]*Schema{
				"config_vars": &Schema{
					Type:     TypeList,
					Optional: true,
					Computed: true,
					Elem: &Schema{
						Type: TypeMap,
					},
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"config_vars.#":     "2",
					"config_vars.0.foo": "baz",
					"config_vars.1.bar": "bar",
				},
			},

			Diff: nil,

			Key: "config_vars",

			Value: []interface{}{
				map[string]interface{}{
					"foo": "baz",
				},
				map[string]interface{}{
					"bar": "bar",
				},
			},
		},

		// #16 List of maps with removal in diff
		{
			Schema: map[string]*Schema{
				"config_vars": &Schema{
					Type:     TypeList,
					Optional: true,
					Computed: true,
					Elem: &Schema{
						Type: TypeMap,
					},
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"config_vars.#":     "1",
					"config_vars.0.FOO": "bar",
				},
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"config_vars.#": &terraform.ResourceAttrDiff{
						Old: "1",
						New: "0",
					},
					"config_vars.0.FOO": &terraform.ResourceAttrDiff{
						Old:        "bar",
						NewRemoved: true,
					},
				},
			},

			Key: "config_vars",

			Value: []interface{}{},
		},

		// #17 Sets
		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeSet,
					Optional: true,
					Computed: true,
					Elem:     &Schema{Type: TypeInt},
					Set: func(a interface{}) int {
						return a.(int)
					},
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"ports.#":  "1",
					"ports.80": "80",
				},
			},

			Diff: nil,

			Key: "ports",

			Value: []interface{}{80},
		},

		// #18
		{
			Schema: map[string]*Schema{
				"data": &Schema{
					Type:     TypeSet,
					Optional: true,
					Elem: &Resource{
						Schema: map[string]*Schema{
							"index": &Schema{
								Type:     TypeInt,
								Required: true,
							},

							"value": &Schema{
								Type:     TypeString,
								Required: true,
							},
						},
					},
					Set: func(a interface{}) int {
						m := a.(map[string]interface{})
						return m["index"].(int)
					},
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"data.#":        "1",
					"data.10.index": "10",
					"data.10.value": "50",
				},
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"data.10.value": &terraform.ResourceAttrDiff{
						Old: "50",
						New: "80",
					},
				},
			},

			Key: "data",

			Value: []interface{}{
				map[string]interface{}{
					"index": 10,
					"value": "80",
				},
			},
		},

		// #19 Empty Set
		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeSet,
					Optional: true,
					Computed: true,
					Elem:     &Schema{Type: TypeInt},
					Set: func(a interface{}) int {
						return a.(int)
					},
				},
			},

			State: nil,

			Diff: nil,

			Key: "ports",

			Value: []interface{}{},
		},

		// #20 Float zero
		{
			Schema: map[string]*Schema{
				"ratio": &Schema{
					Type:     TypeFloat,
					Optional: true,
					Computed: true,
				},
			},

			State: nil,

			Diff: nil,

			Key: "ratio",

			Value: 0.0,
		},

		// #21 Float given
		{
			Schema: map[string]*Schema{
				"ratio": &Schema{
					Type:     TypeFloat,
					Optional: true,
					Computed: true,
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"ratio": "0.5",
				},
			},

			Diff: nil,

			Key: "ratio",

			Value: 0.5,
		},

		// #22 Float diff
		{
			Schema: map[string]*Schema{
				"ratio": &Schema{
					Type:     TypeFloat,
					Optional: true,
					Computed: true,
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"ratio": "-0.5",
				},
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"ratio": &terraform.ResourceAttrDiff{
						Old: "-0.5",
						New: "33.0",
					},
				},
			},

			Key: "ratio",

			Value: 33.0,
		},
	}

	for i, tc := range cases {
		d, err := schemaMap(tc.Schema).Data(tc.State, tc.Diff)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		v := d.Get(tc.Key)
		if s, ok := v.(*Set); ok {
			v = s.List()
		}

		if !reflect.DeepEqual(v, tc.Value) {
			t.Fatalf("Bad: %d\n\n%#v\n\nExpected: %#v", i, v, tc.Value)
		}
	}
}

func TestResourceDataGetChange(t *testing.T) {
	cases := []struct {
		Schema   map[string]*Schema
		State    *terraform.InstanceState
		Diff     *terraform.InstanceDiff
		Key      string
		OldValue interface{}
		NewValue interface{}
	}{
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: nil,

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						Old:         "",
						New:         "foo",
						RequiresNew: true,
					},
				},
			},

			Key: "availability_zone",

			OldValue: "",
			NewValue: "foo",
		},

		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"availability_zone": "foo",
				},
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						Old:         "",
						New:         "foo",
						RequiresNew: true,
					},
				},
			},

			Key: "availability_zone",

			OldValue: "foo",
			NewValue: "foo",
		},
	}

	for i, tc := range cases {
		d, err := schemaMap(tc.Schema).Data(tc.State, tc.Diff)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		o, n := d.GetChange(tc.Key)
		if !reflect.DeepEqual(o, tc.OldValue) {
			t.Fatalf("Old Bad: %d\n\n%#v", i, o)
		}
		if !reflect.DeepEqual(n, tc.NewValue) {
			t.Fatalf("New Bad: %d\n\n%#v", i, n)
		}
	}
}

func TestResourceDataGetOk(t *testing.T) {
	cases := []struct {
		Schema map[string]*Schema
		State  *terraform.InstanceState
		Diff   *terraform.InstanceDiff
		Key    string
		Value  interface{}
		Ok     bool
	}{
		/*
		 * Primitives
		 */
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: nil,

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						Old: "",
						New: "",
					},
				},
			},

			Key:   "availability_zone",
			Value: "",
			Ok:    false,
		},

		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: nil,

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						Old:         "",
						New:         "",
						NewComputed: true,
					},
				},
			},

			Key:   "availability_zone",
			Value: "",
			Ok:    false,
		},

		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: nil,

			Diff: nil,

			Key:   "availability_zone",
			Value: "",
			Ok:    false,
		},

		/*
		 * Lists
		 */

		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeList,
					Optional: true,
					Elem:     &Schema{Type: TypeInt},
				},
			},

			State: nil,

			Diff: nil,

			Key:   "ports",
			Value: []interface{}{},
			Ok:    false,
		},

		/*
		 * Map
		 */

		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeMap,
					Optional: true,
				},
			},

			State: nil,

			Diff: nil,

			Key:   "ports",
			Value: map[string]interface{}{},
			Ok:    false,
		},

		/*
		 * Set
		 */

		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeSet,
					Optional: true,
					Elem:     &Schema{Type: TypeInt},
					Set:      func(a interface{}) int { return a.(int) },
				},
			},

			State: nil,

			Diff: nil,

			Key:   "ports",
			Value: []interface{}{},
			Ok:    false,
		},

		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeSet,
					Optional: true,
					Elem:     &Schema{Type: TypeInt},
					Set:      func(a interface{}) int { return a.(int) },
				},
			},

			State: nil,

			Diff: nil,

			Key:   "ports.0",
			Value: 0,
			Ok:    false,
		},
	}

	for i, tc := range cases {
		d, err := schemaMap(tc.Schema).Data(tc.State, tc.Diff)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		v, ok := d.GetOk(tc.Key)
		if s, ok := v.(*Set); ok {
			v = s.List()
		}

		if !reflect.DeepEqual(v, tc.Value) {
			t.Fatalf("Bad: %d\n\n%#v", i, v)
		}
		if ok != tc.Ok {
			t.Fatalf("Bad: %d\n\n%#v", i, ok)
		}
	}
}

func TestResourceDataHasChange(t *testing.T) {
	cases := []struct {
		Schema map[string]*Schema
		State  *terraform.InstanceState
		Diff   *terraform.InstanceDiff
		Key    string
		Change bool
	}{
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: nil,

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						Old:         "",
						New:         "foo",
						RequiresNew: true,
					},
				},
			},

			Key: "availability_zone",

			Change: true,
		},

		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"availability_zone": "foo",
				},
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						Old:         "",
						New:         "foo",
						RequiresNew: true,
					},
				},
			},

			Key: "availability_zone",

			Change: false,
		},

		{
			Schema: map[string]*Schema{
				"tags": &Schema{
					Type:     TypeMap,
					Optional: true,
					Computed: true,
				},
			},

			State: nil,

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"tags.Name": &terraform.ResourceAttrDiff{
						Old: "foo",
						New: "foo",
					},
				},
			},

			Key: "tags",

			Change: true,
		},

		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeSet,
					Optional: true,
					Elem:     &Schema{Type: TypeInt},
					Set:      func(a interface{}) int { return a.(int) },
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"ports.#":  "1",
					"ports.80": "80",
				},
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"ports.#": &terraform.ResourceAttrDiff{
						Old: "1",
						New: "0",
					},
				},
			},

			Key: "ports",

			Change: true,
		},

		// https://github.com/hashicorp/terraform/issues/927
		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeSet,
					Optional: true,
					Elem:     &Schema{Type: TypeInt},
					Set:      func(a interface{}) int { return a.(int) },
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"ports.#":  "1",
					"ports.80": "80",
				},
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"tags.foo": &terraform.ResourceAttrDiff{
						Old: "",
						New: "bar",
					},
				},
			},

			Key: "ports",

			Change: false,
		},
	}

	for i, tc := range cases {
		d, err := schemaMap(tc.Schema).Data(tc.State, tc.Diff)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		actual := d.HasChange(tc.Key)
		if actual != tc.Change {
			t.Fatalf("Bad: %d %#v", i, actual)
		}
	}
}

func TestResourceDataSet(t *testing.T) {
	cases := []struct {
		Schema   map[string]*Schema
		State    *terraform.InstanceState
		Diff     *terraform.InstanceDiff
		Key      string
		Value    interface{}
		Err      bool
		GetKey   string
		GetValue interface{}

		// GetPreProcess can be set to munge the return value before being
		// compared to GetValue
		GetPreProcess func(interface{}) interface{}
	}{
		// #0: Basic good
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: nil,

			Diff: nil,

			Key:   "availability_zone",
			Value: "foo",

			GetKey:   "availability_zone",
			GetValue: "foo",
		},

		// #1: Basic int
		{
			Schema: map[string]*Schema{
				"port": &Schema{
					Type:     TypeInt,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: nil,

			Diff: nil,

			Key:   "port",
			Value: 80,

			GetKey:   "port",
			GetValue: 80,
		},

		// #2: Basic bool
		{
			Schema: map[string]*Schema{
				"vpc": &Schema{
					Type:     TypeBool,
					Optional: true,
				},
			},

			State: nil,

			Diff: nil,

			Key:   "vpc",
			Value: true,

			GetKey:   "vpc",
			GetValue: true,
		},

		// #3
		{
			Schema: map[string]*Schema{
				"vpc": &Schema{
					Type:     TypeBool,
					Optional: true,
				},
			},

			State: nil,

			Diff: nil,

			Key:   "vpc",
			Value: false,

			GetKey:   "vpc",
			GetValue: false,
		},

		// #4: Invalid type
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: nil,

			Diff: nil,

			Key:   "availability_zone",
			Value: 80,
			Err:   true,

			GetKey:   "availability_zone",
			GetValue: "",
		},

		// #5: List of primitives, set list
		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeList,
					Computed: true,
					Elem:     &Schema{Type: TypeInt},
				},
			},

			State: nil,

			Diff: nil,

			Key:   "ports",
			Value: []int{1, 2, 5},

			GetKey:   "ports",
			GetValue: []interface{}{1, 2, 5},
		},

		// #6: List of primitives, set list with error
		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeList,
					Computed: true,
					Elem:     &Schema{Type: TypeInt},
				},
			},

			State: nil,

			Diff: nil,

			Key:   "ports",
			Value: []interface{}{1, "NOPE", 5},
			Err:   true,

			GetKey:   "ports",
			GetValue: []interface{}{},
		},

		// #7: Set a list of maps
		{
			Schema: map[string]*Schema{
				"config_vars": &Schema{
					Type:     TypeList,
					Optional: true,
					Computed: true,
					Elem: &Schema{
						Type: TypeMap,
					},
				},
			},

			State: nil,

			Diff: nil,

			Key: "config_vars",
			Value: []interface{}{
				map[string]interface{}{
					"foo": "bar",
				},
				map[string]interface{}{
					"bar": "baz",
				},
			},
			Err: false,

			GetKey: "config_vars",
			GetValue: []interface{}{
				map[string]interface{}{
					"foo": "bar",
				},
				map[string]interface{}{
					"bar": "baz",
				},
			},
		},

		// #8: Set, with list
		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeSet,
					Optional: true,
					Computed: true,
					Elem:     &Schema{Type: TypeInt},
					Set: func(a interface{}) int {
						return a.(int)
					},
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"ports.#": "3",
					"ports.0": "100",
					"ports.1": "80",
					"ports.2": "80",
				},
			},

			Key:   "ports",
			Value: []interface{}{100, 125, 125},

			GetKey:   "ports",
			GetValue: []interface{}{100, 125},
		},

		// #9: Set, with Set
		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeSet,
					Optional: true,
					Computed: true,
					Elem:     &Schema{Type: TypeInt},
					Set: func(a interface{}) int {
						return a.(int)
					},
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"ports.#":   "3",
					"ports.100": "100",
					"ports.80":  "80",
					"ports.81":  "81",
				},
			},

			Key: "ports",
			Value: &Set{
				m: map[int]interface{}{
					1: 1,
					2: 2,
				},
			},

			GetKey:   "ports",
			GetValue: []interface{}{1, 2},
		},

		// #10: Set single item
		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeSet,
					Optional: true,
					Computed: true,
					Elem:     &Schema{Type: TypeInt},
					Set: func(a interface{}) int {
						return a.(int)
					},
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"ports.#":   "2",
					"ports.100": "100",
					"ports.80":  "80",
				},
			},

			Key:   "ports.100",
			Value: 256,
			Err:   true,

			GetKey:   "ports",
			GetValue: []interface{}{80, 100},
		},

		// #11: Set with nested set
		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type: TypeSet,
					Elem: &Resource{
						Schema: map[string]*Schema{
							"port": &Schema{
								Type: TypeInt,
							},

							"set": &Schema{
								Type: TypeSet,
								Elem: &Schema{Type: TypeInt},
								Set: func(a interface{}) int {
									return a.(int)
								},
							},
						},
					},
					Set: func(a interface{}) int {
						return a.(map[string]interface{})["port"].(int)
					},
				},
			},

			State: nil,

			Key: "ports",
			Value: []interface{}{
				map[string]interface{}{
					"port": 80,
				},
			},

			GetKey: "ports",
			GetValue: []interface{}{
				map[string]interface{}{
					"port": 80,
					"set":  []interface{}{},
				},
			},

			GetPreProcess: func(v interface{}) interface{} {
				if v == nil {
					return v
				}
				s, ok := v.([]interface{})
				if !ok {
					return v
				}
				for _, v := range s {
					m, ok := v.(map[string]interface{})
					if !ok {
						continue
					}
					if m["set"] == nil {
						continue
					}
					if s, ok := m["set"].(*Set); ok {
						m["set"] = s.List()
					}
				}

				return v
			},
		},

		// #12: List of floats, set list
		{
			Schema: map[string]*Schema{
				"ratios": &Schema{
					Type:     TypeList,
					Computed: true,
					Elem:     &Schema{Type: TypeFloat},
				},
			},

			State: nil,

			Diff: nil,

			Key:   "ratios",
			Value: []float64{1.0, 2.2, 5.5},

			GetKey:   "ratios",
			GetValue: []interface{}{1.0, 2.2, 5.5},
		},

		// #12: Set of floats, set list
		{
			Schema: map[string]*Schema{
				"ratios": &Schema{
					Type:     TypeSet,
					Computed: true,
					Elem:     &Schema{Type: TypeFloat},
					Set: func(a interface{}) int {
						return int(math.Float64bits(a.(float64)))
					},
				},
			},

			State: nil,

			Diff: nil,

			Key:   "ratios",
			Value: []float64{1.0, 2.2, 5.5},

			GetKey:   "ratios",
			GetValue: []interface{}{1.0, 2.2, 5.5},
		},
	}

	for i, tc := range cases {
		d, err := schemaMap(tc.Schema).Data(tc.State, tc.Diff)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		err = d.Set(tc.Key, tc.Value)
		if (err != nil) != tc.Err {
			t.Fatalf("%d err: %s", i, err)
		}

		v := d.Get(tc.GetKey)
		if s, ok := v.(*Set); ok {
			v = s.List()
		}

		if tc.GetPreProcess != nil {
			v = tc.GetPreProcess(v)
		}

		if !reflect.DeepEqual(v, tc.GetValue) {
			t.Fatalf("Get Bad: %d\n\n%#v", i, v)
		}
	}
}

func TestResourceDataState(t *testing.T) {
	cases := []struct {
		Schema  map[string]*Schema
		State   *terraform.InstanceState
		Diff    *terraform.InstanceDiff
		Set     map[string]interface{}
		Result  *terraform.InstanceState
		Partial []string
	}{
		// #0 Basic primitive in diff
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: nil,

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						Old:         "",
						New:         "foo",
						RequiresNew: true,
					},
				},
			},

			Result: &terraform.InstanceState{
				Attributes: map[string]string{
					"availability_zone": "foo",
				},
			},
		},

		// #1 Basic primitive set override
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: nil,

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						Old:         "",
						New:         "foo",
						RequiresNew: true,
					},
				},
			},

			Set: map[string]interface{}{
				"availability_zone": "bar",
			},

			Result: &terraform.InstanceState{
				Attributes: map[string]string{
					"availability_zone": "bar",
				},
			},
		},

		// #2
		{
			Schema: map[string]*Schema{
				"vpc": &Schema{
					Type:     TypeBool,
					Optional: true,
				},
			},

			State: nil,

			Diff: nil,

			Set: map[string]interface{}{
				"vpc": true,
			},

			Result: &terraform.InstanceState{
				Attributes: map[string]string{
					"vpc": "true",
				},
			},
		},

		// #3 Basic primitive with StateFunc set
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:      TypeString,
					Optional:  true,
					Computed:  true,
					StateFunc: func(interface{}) string { return "" },
				},
			},

			State: nil,

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						Old:      "",
						New:      "foo",
						NewExtra: "foo!",
					},
				},
			},

			Result: &terraform.InstanceState{
				Attributes: map[string]string{
					"availability_zone": "foo",
				},
			},
		},

		// #4 List
		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeList,
					Required: true,
					Elem:     &Schema{Type: TypeInt},
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"ports.#": "1",
					"ports.0": "80",
				},
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"ports.#": &terraform.ResourceAttrDiff{
						Old: "1",
						New: "2",
					},
					"ports.1": &terraform.ResourceAttrDiff{
						Old: "",
						New: "100",
					},
				},
			},

			Result: &terraform.InstanceState{
				Attributes: map[string]string{
					"ports.#": "2",
					"ports.0": "80",
					"ports.1": "100",
				},
			},
		},

		// #5 List of resources
		{
			Schema: map[string]*Schema{
				"ingress": &Schema{
					Type:     TypeList,
					Required: true,
					Elem: &Resource{
						Schema: map[string]*Schema{
							"from": &Schema{
								Type:     TypeInt,
								Required: true,
							},
						},
					},
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"ingress.#":      "1",
					"ingress.0.from": "80",
				},
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"ingress.#": &terraform.ResourceAttrDiff{
						Old: "1",
						New: "2",
					},
					"ingress.0.from": &terraform.ResourceAttrDiff{
						Old: "80",
						New: "150",
					},
					"ingress.1.from": &terraform.ResourceAttrDiff{
						Old: "",
						New: "100",
					},
				},
			},

			Result: &terraform.InstanceState{
				Attributes: map[string]string{
					"ingress.#":      "2",
					"ingress.0.from": "150",
					"ingress.1.from": "100",
				},
			},
		},

		// #6 List of maps
		{
			Schema: map[string]*Schema{
				"config_vars": &Schema{
					Type:     TypeList,
					Optional: true,
					Computed: true,
					Elem: &Schema{
						Type: TypeMap,
					},
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"config_vars.#":     "2",
					"config_vars.0.#":   "2",
					"config_vars.0.foo": "bar",
					"config_vars.0.bar": "bar",
					"config_vars.1.#":   "1",
					"config_vars.1.bar": "baz",
				},
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"config_vars.0.bar": &terraform.ResourceAttrDiff{
						NewRemoved: true,
					},
				},
			},

			Set: map[string]interface{}{
				"config_vars": []map[string]interface{}{
					map[string]interface{}{
						"foo": "bar",
					},
					map[string]interface{}{
						"baz": "bang",
					},
				},
			},

			Result: &terraform.InstanceState{
				Attributes: map[string]string{
					"config_vars.#":     "2",
					"config_vars.0.#":   "1",
					"config_vars.0.foo": "bar",
					"config_vars.1.#":   "1",
					"config_vars.1.baz": "bang",
				},
			},
		},

		// #7 List of maps with removal in diff
		{
			Schema: map[string]*Schema{
				"config_vars": &Schema{
					Type:     TypeList,
					Optional: true,
					Computed: true,
					Elem: &Schema{
						Type: TypeMap,
					},
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"config_vars.#":     "1",
					"config_vars.0.FOO": "bar",
				},
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"config_vars.#": &terraform.ResourceAttrDiff{
						Old: "1",
						New: "0",
					},
					"config_vars.0.FOO": &terraform.ResourceAttrDiff{
						Old:        "bar",
						NewRemoved: true,
					},
				},
			},

			Result: &terraform.InstanceState{
				Attributes: map[string]string{
					"config_vars.#": "0",
				},
			},
		},

		// #8 Basic state with other keys
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: &terraform.InstanceState{
				ID: "bar",
				Attributes: map[string]string{
					"id": "bar",
				},
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						Old:         "",
						New:         "foo",
						RequiresNew: true,
					},
				},
			},

			Result: &terraform.InstanceState{
				ID: "bar",
				Attributes: map[string]string{
					"id":                "bar",
					"availability_zone": "foo",
				},
			},
		},

		// #9 Sets
		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeSet,
					Optional: true,
					Computed: true,
					Elem:     &Schema{Type: TypeInt},
					Set: func(a interface{}) int {
						return a.(int)
					},
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"ports.#":   "3",
					"ports.100": "100",
					"ports.80":  "80",
					"ports.81":  "81",
				},
			},

			Diff: nil,

			Result: &terraform.InstanceState{
				Attributes: map[string]string{
					"ports.#":   "3",
					"ports.80":  "80",
					"ports.81":  "81",
					"ports.100": "100",
				},
			},
		},

		// #10
		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeSet,
					Optional: true,
					Computed: true,
					Elem:     &Schema{Type: TypeInt},
					Set: func(a interface{}) int {
						return a.(int)
					},
				},
			},

			State: nil,

			Diff: nil,

			Set: map[string]interface{}{
				"ports": []interface{}{100, 80},
			},

			Result: &terraform.InstanceState{
				Attributes: map[string]string{
					"ports.#":   "2",
					"ports.80":  "80",
					"ports.100": "100",
				},
			},
		},

		// #11
		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeSet,
					Optional: true,
					Computed: true,
					Elem: &Resource{
						Schema: map[string]*Schema{
							"order": &Schema{
								Type: TypeInt,
							},

							"a": &Schema{
								Type: TypeList,
								Elem: &Schema{Type: TypeInt},
							},

							"b": &Schema{
								Type: TypeList,
								Elem: &Schema{Type: TypeInt},
							},
						},
					},
					Set: func(a interface{}) int {
						m := a.(map[string]interface{})
						return m["order"].(int)
					},
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"ports.#":        "2",
					"ports.10.order": "10",
					"ports.10.a.#":   "1",
					"ports.10.a.0":   "80",
					"ports.20.order": "20",
					"ports.20.b.#":   "1",
					"ports.20.b.0":   "100",
				},
			},

			Set: map[string]interface{}{
				"ports": []interface{}{
					map[string]interface{}{
						"order": 20,
						"b":     []interface{}{100},
					},
					map[string]interface{}{
						"order": 10,
						"a":     []interface{}{80},
					},
				},
			},

			Result: &terraform.InstanceState{
				Attributes: map[string]string{
					"ports.#":        "2",
					"ports.10.order": "10",
					"ports.10.a.#":   "1",
					"ports.10.a.0":   "80",
					"ports.10.b.#":   "0",
					"ports.20.order": "20",
					"ports.20.a.#":   "0",
					"ports.20.b.#":   "1",
					"ports.20.b.0":   "100",
				},
			},
		},

		/*
		 * PARTIAL STATES
		 */

		// #12 Basic primitive
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: nil,

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						Old:         "",
						New:         "foo",
						RequiresNew: true,
					},
				},
			},

			Partial: []string{},

			Result: &terraform.InstanceState{
				Attributes: map[string]string{},
			},
		},

		// #13 List
		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeList,
					Required: true,
					Elem:     &Schema{Type: TypeInt},
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"ports.#": "1",
					"ports.0": "80",
				},
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"ports.#": &terraform.ResourceAttrDiff{
						Old: "1",
						New: "2",
					},
					"ports.1": &terraform.ResourceAttrDiff{
						Old: "",
						New: "100",
					},
				},
			},

			Partial: []string{},

			Result: &terraform.InstanceState{
				Attributes: map[string]string{
					"ports.#": "1",
					"ports.0": "80",
				},
			},
		},

		// #14
		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeList,
					Optional: true,
					Computed: true,
					Elem:     &Schema{Type: TypeInt},
				},
			},

			State: nil,

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"ports.#": &terraform.ResourceAttrDiff{
						Old:         "",
						NewComputed: true,
					},
				},
			},

			Partial: []string{},

			Set: map[string]interface{}{
				"ports": []interface{}{},
			},

			Result: &terraform.InstanceState{
				Attributes: map[string]string{},
			},
		},

		// #15 List of resources
		{
			Schema: map[string]*Schema{
				"ingress": &Schema{
					Type:     TypeList,
					Required: true,
					Elem: &Resource{
						Schema: map[string]*Schema{
							"from": &Schema{
								Type:     TypeInt,
								Required: true,
							},
						},
					},
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"ingress.#":      "1",
					"ingress.0.from": "80",
				},
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"ingress.#": &terraform.ResourceAttrDiff{
						Old: "1",
						New: "2",
					},
					"ingress.0.from": &terraform.ResourceAttrDiff{
						Old: "80",
						New: "150",
					},
					"ingress.1.from": &terraform.ResourceAttrDiff{
						Old: "",
						New: "100",
					},
				},
			},

			Partial: []string{},

			Result: &terraform.InstanceState{
				Attributes: map[string]string{
					"ingress.#":      "1",
					"ingress.0.from": "80",
				},
			},
		},

		// #16 List of maps
		{
			Schema: map[string]*Schema{
				"config_vars": &Schema{
					Type:     TypeList,
					Optional: true,
					Computed: true,
					Elem: &Schema{
						Type: TypeMap,
					},
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"config_vars.#":     "2",
					"config_vars.0.foo": "bar",
					"config_vars.0.bar": "bar",
					"config_vars.1.bar": "baz",
				},
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"config_vars.0.bar": &terraform.ResourceAttrDiff{
						NewRemoved: true,
					},
				},
			},

			Set: map[string]interface{}{
				"config_vars": []map[string]interface{}{
					map[string]interface{}{
						"foo": "bar",
					},
					map[string]interface{}{
						"baz": "bang",
					},
				},
			},

			Partial: []string{},

			Result: &terraform.InstanceState{
				Attributes: map[string]string{
					// TODO: broken, shouldn't bar be removed?
					"config_vars.#":     "2",
					"config_vars.0.#":   "2",
					"config_vars.0.foo": "bar",
					"config_vars.0.bar": "bar",
					"config_vars.1.#":   "1",
					"config_vars.1.bar": "baz",
				},
			},
		},

		// #17 Sets
		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeSet,
					Optional: true,
					Computed: true,
					Elem:     &Schema{Type: TypeInt},
					Set: func(a interface{}) int {
						return a.(int)
					},
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"ports.#":   "3",
					"ports.100": "100",
					"ports.80":  "80",
					"ports.81":  "81",
				},
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"ports.120": &terraform.ResourceAttrDiff{
						New: "120",
					},
				},
			},

			Partial: []string{},

			Result: &terraform.InstanceState{
				Attributes: map[string]string{
					"ports.#":   "3",
					"ports.80":  "80",
					"ports.81":  "81",
					"ports.100": "100",
				},
			},
		},

		// #18
		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeSet,
					Optional: true,
					Computed: true,
					Elem:     &Schema{Type: TypeInt},
					Set: func(a interface{}) int {
						return a.(int)
					},
				},
			},

			State: nil,

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"ports.#": &terraform.ResourceAttrDiff{
						Old:         "",
						NewComputed: true,
					},
				},
			},

			Partial: []string{},

			Result: &terraform.InstanceState{
				Attributes: map[string]string{},
			},
		},

		// #19 Maps
		{
			Schema: map[string]*Schema{
				"tags": &Schema{
					Type:     TypeMap,
					Optional: true,
					Computed: true,
				},
			},

			State: nil,

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"tags.Name": &terraform.ResourceAttrDiff{
						Old: "",
						New: "foo",
					},
				},
			},

			Result: &terraform.InstanceState{
				Attributes: map[string]string{
					"tags.#":    "1",
					"tags.Name": "foo",
				},
			},
		},

		// #20
		{
			Schema: map[string]*Schema{
				"tags": &Schema{
					Type:     TypeMap,
					Optional: true,
					Computed: true,
				},
			},

			State: nil,

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"tags.Name": &terraform.ResourceAttrDiff{
						Old: "",
						New: "foo",
					},
				},
			},

			Set: map[string]interface{}{
				"tags": map[string]string{},
			},

			Result: &terraform.InstanceState{
				Attributes: map[string]string{},
			},
		},

		// #21
		{
			Schema: map[string]*Schema{
				"foo": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
				},
			},

			State: nil,

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"foo": &terraform.ResourceAttrDiff{
						NewComputed: true,
					},
				},
			},

			Result: &terraform.InstanceState{
				Attributes: map[string]string{},
			},
		},

		// #22
		{
			Schema: map[string]*Schema{
				"foo": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
				},
			},

			State: nil,

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"foo": &terraform.ResourceAttrDiff{
						NewComputed: true,
					},
				},
			},

			Set: map[string]interface{}{
				"foo": "bar",
			},

			Result: &terraform.InstanceState{
				Attributes: map[string]string{
					"foo": "bar",
				},
			},
		},

		// #23 Set of maps
		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeSet,
					Optional: true,
					Computed: true,
					Elem: &Resource{
						Schema: map[string]*Schema{
							"index": &Schema{Type: TypeInt},
							"uuids": &Schema{Type: TypeMap},
						},
					},
					Set: func(a interface{}) int {
						m := a.(map[string]interface{})
						return m["index"].(int)
					},
				},
			},

			State: nil,

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"ports.10.uuids.#": &terraform.ResourceAttrDiff{
						NewComputed: true,
					},
				},
			},

			Set: map[string]interface{}{
				"ports": []interface{}{
					map[string]interface{}{
						"index": 10,
						"uuids": map[string]interface{}{
							"80": "value",
						},
					},
				},
			},

			Result: &terraform.InstanceState{
				Attributes: map[string]string{
					"ports.#":           "1",
					"ports.10.index":    "10",
					"ports.10.uuids.#":  "1",
					"ports.10.uuids.80": "value",
				},
			},
		},

		// #24
		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeSet,
					Optional: true,
					Computed: true,
					Elem:     &Schema{Type: TypeInt},
					Set: func(a interface{}) int {
						return a.(int)
					},
				},
			},

			State: &terraform.InstanceState{
				Attributes: map[string]string{
					"ports.#":   "3",
					"ports.100": "100",
					"ports.80":  "80",
					"ports.81":  "81",
				},
			},

			Diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"ports.#": &terraform.ResourceAttrDiff{
						Old: "3",
						New: "0",
					},
				},
			},

			Result: &terraform.InstanceState{
				Attributes: map[string]string{
					"ports.#": "0",
				},
			},
		},
	}

	for i, tc := range cases {
		d, err := schemaMap(tc.Schema).Data(tc.State, tc.Diff)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		for k, v := range tc.Set {
			if err := d.Set(k, v); err != nil {
				t.Fatalf("%d err: %s", i, err)
			}
		}

		// Set an ID so that the state returned is not nil
		idSet := false
		if d.Id() == "" {
			idSet = true
			d.SetId("foo")
		}

		// If we have partial, then enable partial state mode.
		if tc.Partial != nil {
			d.Partial(true)
			for _, k := range tc.Partial {
				d.SetPartial(k)
			}
		}

		actual := d.State()

		// If we set an ID, then undo what we did so the comparison works
		if actual != nil && idSet {
			actual.ID = ""
			delete(actual.Attributes, "id")
		}

		if !reflect.DeepEqual(actual, tc.Result) {
			t.Fatalf("Bad: %d\n\n%#v\n\nExpected:\n\n%#v", i, actual, tc.Result)
		}
	}
}

func TestResourceDataSetConnInfo(t *testing.T) {
	d := &ResourceData{}
	d.SetId("foo")
	d.SetConnInfo(map[string]string{
		"foo": "bar",
	})

	expected := map[string]string{
		"foo": "bar",
	}

	actual := d.State()
	if !reflect.DeepEqual(actual.Ephemeral.ConnInfo, expected) {
		t.Fatalf("bad: %#v", actual)
	}
}

func TestResourceDataSetId(t *testing.T) {
	d := &ResourceData{}
	d.SetId("foo")

	actual := d.State()
	if actual.ID != "foo" {
		t.Fatalf("bad: %#v", actual)
	}
}

func TestResourceDataSetId_clear(t *testing.T) {
	d := &ResourceData{
		state: &terraform.InstanceState{ID: "bar"},
	}
	d.SetId("")

	actual := d.State()
	if actual != nil {
		t.Fatalf("bad: %#v", actual)
	}
}

func TestResourceDataSetId_override(t *testing.T) {
	d := &ResourceData{
		state: &terraform.InstanceState{ID: "bar"},
	}
	d.SetId("foo")

	actual := d.State()
	if actual.ID != "foo" {
		t.Fatalf("bad: %#v", actual)
	}
}
