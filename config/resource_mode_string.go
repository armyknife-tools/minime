// Code generated by "stringer -type=ResourceMode -output=resource_mode_string.go resource_mode.go"; DO NOT EDIT

package config

import "fmt"

const _ResourceMode_name = "ManagedResourceModeDataResourceMode"

var _ResourceMode_index = [...]uint8{0, 19, 35}

func (i ResourceMode) String() string {
	if i < 0 || i >= ResourceMode(len(_ResourceMode_index)-1) {
		return fmt.Sprintf("ResourceMode(%d)", i)
	}
	return _ResourceMode_name[_ResourceMode_index[i]:_ResourceMode_index[i+1]]
}
