// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package differ

import (
	"github.com/zclconf/go-cty/cty"

	"github.com/opentofu/opentofu/internal/command/jsonformat/computed"
	"github.com/opentofu/opentofu/internal/command/jsonformat/computed/renderers"
	"github.com/opentofu/opentofu/internal/command/jsonformat/structured"
	"github.com/opentofu/opentofu/internal/command/jsonprovider"
)

func checkForUnknownType(change structured.Change, ctype cty.Type) (computed.Diff, bool) {
	return change.CheckForUnknown(
		false,
		processUnknown,
		createProcessUnknownWithBefore(func(value structured.Change) computed.Diff {
			return ComputeDiffForType(value, ctype)
		}))
}

func checkForUnknownNestedAttribute(change structured.Change, attribute *jsonprovider.NestedType) (computed.Diff, bool) {
	// We want our child attributes to show up as computed instead of deleted.
	// Let's populate that here.
	childUnknown := make(map[string]interface{})
	for key := range attribute.Attributes {
		childUnknown[key] = true
	}

	return change.CheckForUnknown(
		childUnknown,
		processUnknown,
		createProcessUnknownWithBefore(func(value structured.Change) computed.Diff {
			return computeDiffForNestedAttribute(value, attribute)
		}))
}

func checkForUnknownBlock(change structured.Change, block *jsonprovider.Block) (computed.Diff, bool) {
	// We want our child attributes to show up as computed instead of deleted.
	// Let's populate that here.
	childUnknown := make(map[string]interface{})
	for key := range block.Attributes {
		childUnknown[key] = true
	}

	return change.CheckForUnknown(
		childUnknown,
		processUnknown,
		createProcessUnknownWithBefore(func(value structured.Change) computed.Diff {
			return ComputeDiffForBlock(value, block)
		}))
}

func processUnknown(current structured.Change) computed.Diff {
	return asDiff(current, renderers.Unknown(computed.Diff{}))
}

func createProcessUnknownWithBefore(computeDiff func(value structured.Change) computed.Diff) structured.ProcessUnknownWithBefore {
	return func(current structured.Change, before structured.Change) computed.Diff {
		return asDiff(current, renderers.Unknown(computeDiff(before)))
	}
}
