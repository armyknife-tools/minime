// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"reflect"
	"testing"

	"github.com/opentofu/opentofu/internal/legacy/opentf"
)

func TestDiffFieldReader_impl(t *testing.T) {
	var _ FieldReader = new(DiffFieldReader)
}

func TestDiffFieldReader_NestedSetUpdate(t *testing.T) {
	hashFn := func(a interface{}) int {
		m := a.(map[string]interface{})
		return m["val"].(int)
	}

	schema := map[string]*Schema{
		"list_of_sets_1": &Schema{
			Type: TypeList,
			Elem: &Resource{
				Schema: map[string]*Schema{
					"nested_set": &Schema{
						Type: TypeSet,
						Elem: &Resource{
							Schema: map[string]*Schema{
								"val": &Schema{
									Type: TypeInt,
								},
							},
						},
						Set: hashFn,
					},
				},
			},
		},
		"list_of_sets_2": &Schema{
			Type: TypeList,
			Elem: &Resource{
				Schema: map[string]*Schema{
					"nested_set": &Schema{
						Type: TypeSet,
						Elem: &Resource{
							Schema: map[string]*Schema{
								"val": &Schema{
									Type: TypeInt,
								},
							},
						},
						Set: hashFn,
					},
				},
			},
		},
	}

	r := &DiffFieldReader{
		Schema: schema,
		Diff: &opentf.InstanceDiff{
			Attributes: map[string]*opentf.ResourceAttrDiff{
				"list_of_sets_1.0.nested_set.1.val": &opentf.ResourceAttrDiff{
					Old:        "1",
					New:        "0",
					NewRemoved: true,
				},
				"list_of_sets_1.0.nested_set.2.val": &opentf.ResourceAttrDiff{
					New: "2",
				},
			},
		},
	}

	r.Source = &MultiLevelFieldReader{
		Readers: map[string]FieldReader{
			"diff": r,
			"set":  &MapFieldReader{Schema: schema},
			"state": &MapFieldReader{
				Map: &BasicMapReader{
					"list_of_sets_1.#":                  "1",
					"list_of_sets_1.0.nested_set.#":     "1",
					"list_of_sets_1.0.nested_set.1.val": "1",
					"list_of_sets_2.#":                  "1",
					"list_of_sets_2.0.nested_set.#":     "1",
					"list_of_sets_2.0.nested_set.1.val": "1",
				},
				Schema: schema,
			},
		},
		Levels: []string{"state", "config"},
	}

	out, err := r.ReadField([]string{"list_of_sets_2"})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	s := &Set{F: hashFn}
	s.Add(map[string]interface{}{"val": 1})
	expected := s.List()

	l := out.Value.([]interface{})
	i := l[0].(map[string]interface{})
	actual := i["nested_set"].(*Set).List()

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("bad: NestedSetUpdate\n\nexpected: %#v\n\ngot: %#v\n\n", expected, actual)
	}
}

// https://github.com/hashicorp/terraform/issues/914
func TestDiffFieldReader_MapHandling(t *testing.T) {
	schema := map[string]*Schema{
		"tags": &Schema{
			Type: TypeMap,
		},
	}
	r := &DiffFieldReader{
		Schema: schema,
		Diff: &opentf.InstanceDiff{
			Attributes: map[string]*opentf.ResourceAttrDiff{
				"tags.%": &opentf.ResourceAttrDiff{
					Old: "1",
					New: "2",
				},
				"tags.baz": &opentf.ResourceAttrDiff{
					Old: "",
					New: "qux",
				},
			},
		},
		Source: &MapFieldReader{
			Schema: schema,
			Map: BasicMapReader(map[string]string{
				"tags.%":   "1",
				"tags.foo": "bar",
			}),
		},
	}

	result, err := r.ReadField([]string{"tags"})
	if err != nil {
		t.Fatalf("ReadField failed: %#v", err)
	}

	expected := map[string]interface{}{
		"foo": "bar",
		"baz": "qux",
	}

	if !reflect.DeepEqual(expected, result.Value) {
		t.Fatalf("bad: DiffHandling\n\nexpected: %#v\n\ngot: %#v\n\n", expected, result.Value)
	}
}

func TestDiffFieldReader_extra(t *testing.T) {
	schema := map[string]*Schema{
		"stringComputed": &Schema{Type: TypeString},

		"listMap": &Schema{
			Type: TypeList,
			Elem: &Schema{
				Type: TypeMap,
			},
		},

		"mapRemove": &Schema{Type: TypeMap},

		"setChange": &Schema{
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

		"setEmpty": &Schema{
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
	}

	r := &DiffFieldReader{
		Schema: schema,
		Diff: &opentf.InstanceDiff{
			Attributes: map[string]*opentf.ResourceAttrDiff{
				"stringComputed": &opentf.ResourceAttrDiff{
					Old:         "foo",
					New:         "bar",
					NewComputed: true,
				},

				"listMap.0.bar": &opentf.ResourceAttrDiff{
					NewRemoved: true,
				},

				"mapRemove.bar": &opentf.ResourceAttrDiff{
					NewRemoved: true,
				},

				"setChange.10.value": &opentf.ResourceAttrDiff{
					Old: "50",
					New: "80",
				},

				"setEmpty.#": &opentf.ResourceAttrDiff{
					Old: "2",
					New: "0",
				},
			},
		},

		Source: &MapFieldReader{
			Schema: schema,
			Map: BasicMapReader(map[string]string{
				"listMap.#":     "2",
				"listMap.0.foo": "bar",
				"listMap.0.bar": "baz",
				"listMap.1.baz": "baz",

				"mapRemove.foo": "bar",
				"mapRemove.bar": "bar",

				"setChange.#":        "1",
				"setChange.10.index": "10",
				"setChange.10.value": "50",

				"setEmpty.#":        "2",
				"setEmpty.10.index": "10",
				"setEmpty.10.value": "50",
				"setEmpty.20.index": "20",
				"setEmpty.20.value": "50",
			}),
		},
	}

	cases := map[string]struct {
		Addr   []string
		Result FieldReadResult
		Err    bool
	}{
		"stringComputed": {
			[]string{"stringComputed"},
			FieldReadResult{
				Value:    "",
				Exists:   true,
				Computed: true,
			},
			false,
		},

		"listMapRemoval": {
			[]string{"listMap"},
			FieldReadResult{
				Value: []interface{}{
					map[string]interface{}{
						"foo": "bar",
					},
					map[string]interface{}{
						"baz": "baz",
					},
				},
				Exists: true,
			},
			false,
		},

		"mapRemove": {
			[]string{"mapRemove"},
			FieldReadResult{
				Value: map[string]interface{}{
					"foo": "bar",
				},
				Exists:   true,
				Computed: false,
			},
			false,
		},

		"setChange": {
			[]string{"setChange"},
			FieldReadResult{
				Value: []interface{}{
					map[string]interface{}{
						"index": 10,
						"value": "80",
					},
				},
				Exists: true,
			},
			false,
		},

		"setEmpty": {
			[]string{"setEmpty"},
			FieldReadResult{
				Value:  []interface{}{},
				Exists: true,
			},
			false,
		},
	}

	for name, tc := range cases {
		out, err := r.ReadField(tc.Addr)
		if err != nil != tc.Err {
			t.Fatalf("%s: err: %s", name, err)
		}
		if s, ok := out.Value.(*Set); ok {
			// If it is a set, convert to a list so its more easily checked.
			out.Value = s.List()
		}
		if !reflect.DeepEqual(tc.Result, out) {
			t.Fatalf("%s: bad: %#v", name, out)
		}
	}
}

func TestDiffFieldReader(t *testing.T) {
	testFieldReader(t, func(s map[string]*Schema) FieldReader {
		return &DiffFieldReader{
			Schema: s,
			Diff: &opentf.InstanceDiff{
				Attributes: map[string]*opentf.ResourceAttrDiff{
					"bool": &opentf.ResourceAttrDiff{
						Old: "",
						New: "true",
					},

					"int": &opentf.ResourceAttrDiff{
						Old: "",
						New: "42",
					},

					"float": &opentf.ResourceAttrDiff{
						Old: "",
						New: "3.1415",
					},

					"string": &opentf.ResourceAttrDiff{
						Old: "",
						New: "string",
					},

					"stringComputed": &opentf.ResourceAttrDiff{
						Old:         "foo",
						New:         "bar",
						NewComputed: true,
					},

					"list.#": &opentf.ResourceAttrDiff{
						Old: "0",
						New: "2",
					},

					"list.0": &opentf.ResourceAttrDiff{
						Old: "",
						New: "foo",
					},

					"list.1": &opentf.ResourceAttrDiff{
						Old: "",
						New: "bar",
					},

					"listInt.#": &opentf.ResourceAttrDiff{
						Old: "0",
						New: "2",
					},

					"listInt.0": &opentf.ResourceAttrDiff{
						Old: "",
						New: "21",
					},

					"listInt.1": &opentf.ResourceAttrDiff{
						Old: "",
						New: "42",
					},

					"map.foo": &opentf.ResourceAttrDiff{
						Old: "",
						New: "bar",
					},

					"map.bar": &opentf.ResourceAttrDiff{
						Old: "",
						New: "baz",
					},

					"mapInt.%": &opentf.ResourceAttrDiff{
						Old: "",
						New: "2",
					},
					"mapInt.one": &opentf.ResourceAttrDiff{
						Old: "",
						New: "1",
					},
					"mapInt.two": &opentf.ResourceAttrDiff{
						Old: "",
						New: "2",
					},

					"mapIntNestedSchema.%": &opentf.ResourceAttrDiff{
						Old: "",
						New: "2",
					},
					"mapIntNestedSchema.one": &opentf.ResourceAttrDiff{
						Old: "",
						New: "1",
					},
					"mapIntNestedSchema.two": &opentf.ResourceAttrDiff{
						Old: "",
						New: "2",
					},

					"mapFloat.%": &opentf.ResourceAttrDiff{
						Old: "",
						New: "1",
					},
					"mapFloat.oneDotTwo": &opentf.ResourceAttrDiff{
						Old: "",
						New: "1.2",
					},

					"mapBool.%": &opentf.ResourceAttrDiff{
						Old: "",
						New: "2",
					},
					"mapBool.True": &opentf.ResourceAttrDiff{
						Old: "",
						New: "true",
					},
					"mapBool.False": &opentf.ResourceAttrDiff{
						Old: "",
						New: "false",
					},

					"set.#": &opentf.ResourceAttrDiff{
						Old: "0",
						New: "2",
					},

					"set.10": &opentf.ResourceAttrDiff{
						Old: "",
						New: "10",
					},

					"set.50": &opentf.ResourceAttrDiff{
						Old: "",
						New: "50",
					},

					"setDeep.#": &opentf.ResourceAttrDiff{
						Old: "0",
						New: "2",
					},

					"setDeep.10.index": &opentf.ResourceAttrDiff{
						Old: "",
						New: "10",
					},

					"setDeep.10.value": &opentf.ResourceAttrDiff{
						Old: "",
						New: "foo",
					},

					"setDeep.50.index": &opentf.ResourceAttrDiff{
						Old: "",
						New: "50",
					},

					"setDeep.50.value": &opentf.ResourceAttrDiff{
						Old: "",
						New: "bar",
					},
				},
			},

			Source: &MapFieldReader{
				Schema: s,
				Map: BasicMapReader(map[string]string{
					"listMap.#":     "2",
					"listMap.0.foo": "bar",
					"listMap.0.bar": "baz",
					"listMap.1.baz": "baz",
				}),
			},
		}
	})
}
