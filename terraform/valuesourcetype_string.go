// Code generated by "stringer -type ValueSourceType"; DO NOT EDIT.

package terraform

import "strconv"

const (
	_ValueSourceType_name_0 = "ValueFromUnknown"
	_ValueSourceType_name_1 = "ValueFromCLIArg"
	_ValueSourceType_name_2 = "ValueFromConfig"
	_ValueSourceType_name_3 = "ValueFromEnvVarValueFromAutoFile"
	_ValueSourceType_name_4 = "ValueFromInput"
	_ValueSourceType_name_5 = "ValueFromNamedFile"
	_ValueSourceType_name_6 = "ValueFromPlan"
	_ValueSourceType_name_7 = "ValueFromCaller"
)

var (
	_ValueSourceType_index_3 = [...]uint8{0, 15, 32}
)

func (i ValueSourceType) String() string {
	switch {
	case i == 0:
		return _ValueSourceType_name_0
	case i == 65:
		return _ValueSourceType_name_1
	case i == 67:
		return _ValueSourceType_name_2
	case 69 <= i && i <= 70:
		i -= 69
		return _ValueSourceType_name_3[_ValueSourceType_index_3[i]:_ValueSourceType_index_3[i+1]]
	case i == 73:
		return _ValueSourceType_name_4
	case i == 78:
		return _ValueSourceType_name_5
	case i == 80:
		return _ValueSourceType_name_6
	case i == 83:
		return _ValueSourceType_name_7
	default:
		return "ValueSourceType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
