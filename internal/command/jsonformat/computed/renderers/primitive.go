package renderers

import (
	"fmt"
	"strings"

	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/terraform/internal/command/format"
	"github.com/hashicorp/terraform/internal/command/jsonformat/collections"
	"github.com/hashicorp/terraform/internal/command/jsonformat/computed"
	"github.com/hashicorp/terraform/internal/plans"
)

var _ computed.DiffRenderer = (*primitiveRenderer)(nil)

func Primitive(before, after interface{}, ctype cty.Type) computed.DiffRenderer {
	return &primitiveRenderer{
		before: before,
		after:  after,
		ctype:  ctype,
	}
}

type primitiveRenderer struct {
	NoWarningsRenderer

	before interface{}
	after  interface{}
	ctype  cty.Type
}

func (renderer primitiveRenderer) RenderHuman(diff computed.Diff, indent int, opts computed.RenderHumanOpts) string {
	if renderer.ctype == cty.String {
		return renderer.renderStringDiff(diff, indent, opts)
	}

	beforeValue := renderPrimitiveValue(renderer.before, renderer.ctype)
	afterValue := renderPrimitiveValue(renderer.after, renderer.ctype)

	switch diff.Action {
	case plans.Create:
		return fmt.Sprintf("%s%s", afterValue, forcesReplacement(diff.Replace, opts.OverrideForcesReplacement))
	case plans.Delete:
		return fmt.Sprintf("%s%s%s", beforeValue, nullSuffix(opts.OverrideNullSuffix, diff.Action), forcesReplacement(diff.Replace, opts.OverrideForcesReplacement))
	case plans.NoOp:
		return fmt.Sprintf("%s%s", beforeValue, forcesReplacement(diff.Replace, opts.OverrideForcesReplacement))
	default:
		return fmt.Sprintf("%s [yellow]->[reset] %s%s", beforeValue, afterValue, forcesReplacement(diff.Replace, opts.OverrideForcesReplacement))
	}
}

func renderPrimitiveValue(value interface{}, t cty.Type) string {
	if value == nil {
		return "[dark_gray]null[reset]"
	}

	switch {
	case t == cty.Bool:
		if value.(bool) {
			return "true"
		}
		return "false"
	case t == cty.Number:
		return fmt.Sprintf("%g", value)
	default:
		panic("unrecognized primitive type: " + t.FriendlyName())
	}
}

func (renderer primitiveRenderer) renderStringDiff(diff computed.Diff, indent int, opts computed.RenderHumanOpts) string {

	// We process multiline strings at the end of the switch statement.
	var lines []string

	switch diff.Action {
	case plans.Create, plans.NoOp:
		str := evaluatePrimitiveString(renderer.after)

		if str.Json != nil {
			if diff.Action == plans.NoOp {
				return renderer.renderStringDiffAsJson(diff, indent, opts, str, str)
			} else {
				return renderer.renderStringDiffAsJson(diff, indent, opts, evaluatedString{}, str)
			}
		}

		if !str.IsMultiline {
			return fmt.Sprintf("%q%s", str.String, forcesReplacement(diff.Replace, opts.OverrideForcesReplacement))
		}

		// We are creating a single multiline string, so let's split by the new
		// line character. While we are doing this, we are going to insert our
		// indents and make sure each line is formatted correctly.
		lines = strings.Split(strings.ReplaceAll(str.String, "\n", fmt.Sprintf("\n%s%s ", formatIndent(indent), format.DiffActionSymbol(plans.NoOp))), "\n")

		// We now just need to do the same for the first entry in lines, because
		// we split on the new line characters which won't have been at the
		// beginning of the first line.
		lines[0] = fmt.Sprintf("%s%s %s", formatIndent(indent), format.DiffActionSymbol(plans.NoOp), lines[0])
	case plans.Delete:
		str := evaluatePrimitiveString(renderer.before)

		if str.Json != nil {
			return renderer.renderStringDiffAsJson(diff, indent, opts, str, evaluatedString{})
		}

		if !str.IsMultiline {
			return fmt.Sprintf("%q%s%s", str.String, nullSuffix(opts.OverrideNullSuffix, diff.Action), forcesReplacement(diff.Replace, opts.OverrideForcesReplacement))
		}

		// We are creating a single multiline string, so let's split by the new
		// line character. While we are doing this, we are going to insert our
		// indents and make sure each line is formatted correctly.
		lines = strings.Split(strings.ReplaceAll(str.String, "\n", fmt.Sprintf("\n%s%s ", formatIndent(indent), format.DiffActionSymbol(plans.NoOp))), "\n")

		// We now just need to do the same for the first entry in lines, because
		// we split on the new line characters which won't have been at the
		// beginning of the first line.
		lines[0] = fmt.Sprintf("%s%s %s", formatIndent(indent), format.DiffActionSymbol(plans.NoOp), lines[0])
	default:
		beforeString := evaluatePrimitiveString(renderer.before)
		afterString := evaluatePrimitiveString(renderer.after)

		if beforeString.Json != nil && afterString.Json != nil {
			return renderer.renderStringDiffAsJson(diff, indent, opts, beforeString, afterString)
		}

		if beforeString.Json != nil || afterString.Json != nil {
			// This means one of the strings is JSON and one isn't. We're going
			// to be a little inefficient here, but we can just reuse another
			// renderer for this so let's keep it simple.
			return computed.NewDiff(
				TypeChange(
					computed.NewDiff(Primitive(renderer.before, nil, cty.String), plans.Delete, false),
					computed.NewDiff(Primitive(nil, renderer.after, cty.String), plans.Create, false)),
				diff.Action,
				diff.Replace).RenderHuman(indent, opts)
		}

		if !beforeString.IsMultiline && !afterString.IsMultiline {
			return fmt.Sprintf("%q [yellow]->[reset] %q%s", beforeString.String, afterString.String, forcesReplacement(diff.Replace, opts.OverrideForcesReplacement))
		}

		beforeLines := strings.Split(beforeString.String, "\n")
		afterLines := strings.Split(afterString.String, "\n")

		processIndices := func(beforeIx, afterIx int) {
			if beforeIx < 0 || beforeIx >= len(beforeLines) {
				lines = append(lines, fmt.Sprintf("%s%s %s", formatIndent(indent), format.DiffActionSymbol(plans.Create), afterLines[afterIx]))
				return
			}

			if afterIx < 0 || afterIx >= len(afterLines) {
				lines = append(lines, fmt.Sprintf("%s%s %s", formatIndent(indent), format.DiffActionSymbol(plans.Delete), beforeLines[beforeIx]))
				return
			}

			lines = append(lines, fmt.Sprintf("%s%s %s", formatIndent(indent), format.DiffActionSymbol(plans.NoOp), beforeLines[beforeIx]))
		}
		isObjType := func(_ string) bool {
			return false
		}

		collections.ProcessSlice(beforeLines, afterLines, processIndices, isObjType)
	}

	// We return early if we find non-multiline strings or JSON strings, so we
	// know here that we just render the lines slice properly.
	return fmt.Sprintf("<<-EOT%s\n%s\n%sEOT%s",
		forcesReplacement(diff.Replace, opts.OverrideForcesReplacement),
		strings.Join(lines, "\n"),
		formatIndent(indent),
		nullSuffix(opts.OverrideNullSuffix, diff.Action))
}

func (renderer primitiveRenderer) renderStringDiffAsJson(diff computed.Diff, indent int, opts computed.RenderHumanOpts, before evaluatedString, after evaluatedString) string {
	jsonDiff := RendererJsonOpts().Transform(before.Json, after.Json)

	var whitespace, replace string
	if jsonDiff.Action == plans.NoOp && diff.Action == plans.Update {
		// Then this means we are rendering a whitespace only change. The JSON
		// differ will have ignored the whitespace changes so that makes the
		// diff we are about to print out very confusing without extra
		// explanation.
		if diff.Replace {
			whitespace = " # whitespace changes force replacement"
		} else {
			whitespace = " # whitespace changes"
		}
	} else {
		// We only show the replace suffix if we didn't print something out
		// about whitespace changes.
		replace = forcesReplacement(diff.Replace, opts.OverrideForcesReplacement)
	}

	renderedJsonDiff := jsonDiff.RenderHuman(indent, opts)

	if strings.Contains(renderedJsonDiff, "\n") {
		return fmt.Sprintf("jsonencode(%s\n%s%s %s%s\n%s)", whitespace, formatIndent(indent), format.DiffActionSymbol(diff.Action), renderedJsonDiff, replace, formatIndent(indent))
	}
	return fmt.Sprintf("jsonencode(%s)%s%s", renderedJsonDiff, whitespace, replace)
}
