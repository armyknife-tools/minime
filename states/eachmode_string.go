// Code generated by "stringer -type EachMode"; DO NOT EDIT.

package states

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[NoEach-0]
	_ = x[EachList-76]
	_ = x[EachMap-77]
}

const (
	_EachMode_name_0 = "NoEach"
	_EachMode_name_1 = "EachListEachMap"
)

var (
	_EachMode_index_1 = [...]uint8{0, 8, 15}
)

func (i EachMode) String() string {
	switch {
	case i == 0:
		return _EachMode_name_0
	case 76 <= i && i <= 77:
		i -= 76
		return _EachMode_name_1[_EachMode_index_1[i]:_EachMode_index_1[i+1]]
	default:
		return "EachMode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
