// Code generated by "stringer -type=ResourceInstanceChangeActionReason changes.go"; DO NOT EDIT.

package plans

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ResourceInstanceChangeNoReason-0]
	_ = x[ResourceInstanceReplaceBecauseTainted-84]
	_ = x[ResourceInstanceReplaceByRequest-82]
	_ = x[ResourceInstanceReplaceByTriggers-68]
	_ = x[ResourceInstanceReplaceBecauseCannotUpdate-70]
	_ = x[ResourceInstanceDeleteBecauseNoResourceConfig-78]
	_ = x[ResourceInstanceDeleteBecauseWrongRepetition-87]
	_ = x[ResourceInstanceDeleteBecauseCountIndex-67]
	_ = x[ResourceInstanceDeleteBecauseEachKey-69]
	_ = x[ResourceInstanceDeleteBecauseNoModule-77]
	_ = x[ResourceInstanceDeleteBecauseNoMoveTarget-65]
	_ = x[ResourceInstanceReadBecauseConfigUnknown-63]
	_ = x[ResourceInstanceReadBecauseDependencyPending-33]
	_ = x[ResourceInstanceReadBecauseCheckNested-35]
}

const (
	_ResourceInstanceChangeActionReason_name_0 = "ResourceInstanceChangeNoReason"
	_ResourceInstanceChangeActionReason_name_1 = "ResourceInstanceReadBecauseDependencyPending"
	_ResourceInstanceChangeActionReason_name_2 = "ResourceInstanceReadBecauseCheckNested"
	_ResourceInstanceChangeActionReason_name_3 = "ResourceInstanceReadBecauseConfigUnknown"
	_ResourceInstanceChangeActionReason_name_4 = "ResourceInstanceDeleteBecauseNoMoveTarget"
	_ResourceInstanceChangeActionReason_name_5 = "ResourceInstanceDeleteBecauseCountIndexResourceInstanceReplaceByTriggersResourceInstanceDeleteBecauseEachKeyResourceInstanceReplaceBecauseCannotUpdate"
	_ResourceInstanceChangeActionReason_name_6 = "ResourceInstanceDeleteBecauseNoModuleResourceInstanceDeleteBecauseNoResourceConfig"
	_ResourceInstanceChangeActionReason_name_7 = "ResourceInstanceReplaceByRequest"
	_ResourceInstanceChangeActionReason_name_8 = "ResourceInstanceReplaceBecauseTainted"
	_ResourceInstanceChangeActionReason_name_9 = "ResourceInstanceDeleteBecauseWrongRepetition"
)

var (
	_ResourceInstanceChangeActionReason_index_5 = [...]uint8{0, 39, 72, 108, 150}
	_ResourceInstanceChangeActionReason_index_6 = [...]uint8{0, 37, 82}
)

func (i ResourceInstanceChangeActionReason) String() string {
	switch {
	case i == 0:
		return _ResourceInstanceChangeActionReason_name_0
	case i == 33:
		return _ResourceInstanceChangeActionReason_name_1
	case i == 35:
		return _ResourceInstanceChangeActionReason_name_2
	case i == 63:
		return _ResourceInstanceChangeActionReason_name_3
	case i == 65:
		return _ResourceInstanceChangeActionReason_name_4
	case 67 <= i && i <= 70:
		i -= 67
		return _ResourceInstanceChangeActionReason_name_5[_ResourceInstanceChangeActionReason_index_5[i]:_ResourceInstanceChangeActionReason_index_5[i+1]]
	case 77 <= i && i <= 78:
		i -= 77
		return _ResourceInstanceChangeActionReason_name_6[_ResourceInstanceChangeActionReason_index_6[i]:_ResourceInstanceChangeActionReason_index_6[i+1]]
	case i == 82:
		return _ResourceInstanceChangeActionReason_name_7
	case i == 84:
		return _ResourceInstanceChangeActionReason_name_8
	case i == 87:
		return _ResourceInstanceChangeActionReason_name_9
	default:
		return "ResourceInstanceChangeActionReason(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
