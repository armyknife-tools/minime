// Code generated by "stringer -type=GraphType context_graph_type.go"; DO NOT EDIT.

package terraform

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[GraphTypeInvalid-0]
	_ = x[GraphTypePlan-1]
	_ = x[GraphTypePlanDestroy-2]
	_ = x[GraphTypeApply-3]
	_ = x[GraphTypeValidate-4]
	_ = x[GraphTypeEval-5]
}

const _GraphType_name = "GraphTypeInvalidGraphTypePlanGraphTypePlanDestroyGraphTypeApplyGraphTypeValidateGraphTypeEval"

var _GraphType_index = [...]uint8{0, 16, 29, 49, 63, 80, 93}

func (i GraphType) String() string {
	if i >= GraphType(len(_GraphType_index)-1) {
		return "GraphType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _GraphType_name[_GraphType_index[i]:_GraphType_index[i+1]]
}
