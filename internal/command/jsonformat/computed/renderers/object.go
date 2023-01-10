package renderers

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform/internal/command/jsonformat/computed"

	"github.com/hashicorp/terraform/internal/command/format"
	"github.com/hashicorp/terraform/internal/plans"
)

var _ computed.DiffRenderer = (*objectRenderer)(nil)

func Object(attributes map[string]computed.Diff) computed.DiffRenderer {
	return &objectRenderer{
		attributes:         attributes,
		overrideNullSuffix: true,
	}
}

func NestedObject(attributes map[string]computed.Diff) computed.DiffRenderer {
	return &objectRenderer{
		attributes:         attributes,
		overrideNullSuffix: false,
	}
}

type objectRenderer struct {
	NoWarningsRenderer

	attributes         map[string]computed.Diff
	overrideNullSuffix bool
}

func (renderer objectRenderer) RenderHuman(diff computed.Diff, indent int, opts computed.RenderHumanOpts) string {
	if len(renderer.attributes) == 0 {
		return fmt.Sprintf("{}%s%s", nullSuffix(opts.OverrideNullSuffix, diff.Action), forcesReplacement(diff.Replace))
	}

	attributeOpts := opts.Clone()
	attributeOpts.OverrideNullSuffix = renderer.overrideNullSuffix

	// We need to keep track of our keys in two ways. The first is the order in
	// which we will display them. The second is a mapping to their safely
	// escaped equivalent.

	maximumKeyLen := 0
	var keys []string
	escapedKeys := make(map[string]string)
	for key := range renderer.attributes {
		keys = append(keys, key)
		escapedKey := ensureValidAttributeName(key)
		escapedKeys[key] = escapedKey
		if maximumKeyLen < len(escapedKey) {
			maximumKeyLen = len(escapedKey)
		}
	}
	sort.Strings(keys)

	unchangedAttributes := 0
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("{%s\n", forcesReplacement(diff.Replace)))
	for _, key := range keys {
		attribute := renderer.attributes[key]

		if attribute.Action == plans.NoOp && !opts.ShowUnchangedChildren {
			// Don't render NoOp operations when we are compact display.
			unchangedAttributes++
			continue
		}

		for _, warning := range attribute.WarningsHuman(indent + 1) {
			buf.WriteString(fmt.Sprintf("%s%s\n", formatIndent(indent+1), warning))
		}
		buf.WriteString(fmt.Sprintf("%s%s %-*s = %s\n", formatIndent(indent+1), format.DiffActionSymbol(attribute.Action), maximumKeyLen, escapedKeys[key], attribute.RenderHuman(indent+1, attributeOpts)))
	}

	if unchangedAttributes > 0 {
		buf.WriteString(fmt.Sprintf("%s%s %s\n", formatIndent(indent+1), format.DiffActionSymbol(plans.NoOp), unchanged("attribute", unchangedAttributes)))
	}

	buf.WriteString(fmt.Sprintf("%s%s }%s", formatIndent(indent), format.DiffActionSymbol(plans.NoOp), nullSuffix(opts.OverrideNullSuffix, diff.Action)))
	return buf.String()
}
