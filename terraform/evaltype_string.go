// generated by stringer -type=EvalType eval_type.go; DO NOT EDIT

package terraform

import "fmt"

const (
	_EvalType_name_0 = "EvalTypeInvalid"
	_EvalType_name_1 = "EvalTypeNull"
	_EvalType_name_2 = "EvalTypeConfig"
	_EvalType_name_3 = "EvalTypeResourceProvider"
)

var (
	_EvalType_index_0 = [...]uint8{0, 15}
	_EvalType_index_1 = [...]uint8{0, 12}
	_EvalType_index_2 = [...]uint8{0, 14}
	_EvalType_index_3 = [...]uint8{0, 24}
)

func (i EvalType) String() string {
	switch {
	case i == 0:
		return _EvalType_name_0
	case i == 2:
		return _EvalType_name_1
	case i == 4:
		return _EvalType_name_2
	case i == 8:
		return _EvalType_name_3
	default:
		return fmt.Sprintf("EvalType(%d)", i)
	}
}
