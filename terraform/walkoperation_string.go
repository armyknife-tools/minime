// generated by stringer -type=walkOperation graph_walk_operation.go; DO NOT EDIT

package terraform

import "fmt"

const _walkOperation_name = "walkInvalidwalkInputwalkApplywalkPlanwalkPlanDestroywalkRefreshwalkValidatewalkDestroy"

var _walkOperation_index = [...]uint8{0, 11, 20, 29, 37, 52, 63, 75, 86}

func (i walkOperation) String() string {
	if i >= walkOperation(len(_walkOperation_index)-1) {
		return fmt.Sprintf("walkOperation(%d)", i)
	}
	return _walkOperation_name[_walkOperation_index[i]:_walkOperation_index[i+1]]
}
