// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tofu

import (
	"testing"
)

func TestCallbackUIOutput_impl(t *testing.T) {
	var _ UIOutput = new(CallbackUIOutput)
}
