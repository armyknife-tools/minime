// Code generated by "stringer -type ProvisionerOnFailure"; DO NOT EDIT.

package configs

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ProvisionerOnFailureInvalid-0]
	_ = x[ProvisionerOnFailureContinue-1]
	_ = x[ProvisionerOnFailureFail-2]
}

const _ProvisionerOnFailure_name = "ProvisionerOnFailureInvalidProvisionerOnFailureContinueProvisionerOnFailureFail"

var _ProvisionerOnFailure_index = [...]uint8{0, 27, 55, 79}

func (i ProvisionerOnFailure) String() string {
	if i < 0 || i >= ProvisionerOnFailure(len(_ProvisionerOnFailure_index)-1) {
		return "ProvisionerOnFailure(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ProvisionerOnFailure_name[_ProvisionerOnFailure_index[i]:_ProvisionerOnFailure_index[i+1]]
}
