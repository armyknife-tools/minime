// Code generated by "stringer -type=countHookAction hook_count_action.go"; DO NOT EDIT.

package local

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[countHookActionAdd-0]
	_ = x[countHookActionChange-1]
	_ = x[countHookActionRemove-2]
}

const _countHookAction_name = "countHookActionAddcountHookActionChangecountHookActionRemove"

var _countHookAction_index = [...]uint8{0, 18, 39, 60}

func (i countHookAction) String() string {
	if i >= countHookAction(len(_countHookAction_index)-1) {
		return "countHookAction(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _countHookAction_name[_countHookAction_index[i]:_countHookAction_index[i+1]]
}
