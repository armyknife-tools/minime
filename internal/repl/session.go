package repl

import (
	"strings"

	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform/internal/lang"
	"github.com/hashicorp/terraform/internal/lang/marks"
	"github.com/hashicorp/terraform/internal/tfdiags"
)

// Session represents the state for a single REPL session.
type Session struct {
	// Scope is the evaluation scope where expressions will be evaluated.
	Scope *lang.Scope
}

// Handle handles a single line of input from the REPL.
//
// This is a stateful operation if a command is given (such as setting
// a variable). This function should not be called in parallel.
//
// The return value is the output and the error to show.
func (s *Session) Handle(line string) (string, bool, tfdiags.Diagnostics) {
	switch {
	case strings.TrimSpace(line) == "":
		return "", false, nil
	case strings.TrimSpace(line) == "exit":
		return "", true, nil
	case strings.TrimSpace(line) == "help":
		ret, diags := s.handleHelp()
		return ret, false, diags
	default:
		ret, diags := s.handleEval(line)
		return ret, false, diags
	}
}

func (s *Session) handleEval(line string) (string, tfdiags.Diagnostics) {
	var diags tfdiags.Diagnostics

	// Parse the given line as an expression
	expr, parseDiags := hclsyntax.ParseExpression([]byte(line), "<console-input>", hcl.Pos{Line: 1, Column: 1})
	diags = diags.Append(parseDiags)
	if parseDiags.HasErrors() {
		return "", diags
	}

	val, valDiags := s.Scope.EvalExpr(expr, cty.DynamicPseudoType)
	diags = diags.Append(valDiags)
	if valDiags.HasErrors() {
		return "", diags
	}

	// The raw mark is used only by the console-only `type` function, in order
	// to allow display of a string value representation of the type without the
	// usual HCL formatting. If we receive a string value with this mark, we do
	// not want to format it any further.
	//
	// Due to mark propagation in cty, calling `type` as part of a larger
	// expression can lead to other values being marked, which can in turn lead
	// to unpredictable results. If any non-string value has the raw mark, we
	// return a diagnostic explaining that this use of `type` is not permitted.
	if marks.Contains(val, marks.Raw) {
		if val.Type().Equals(cty.String) {
			raw, _ := val.Unmark()
			return raw.AsString(), diags
		} else {
			diags = diags.Append(tfdiags.Sourceless(
				tfdiags.Error,
				"Invalid use of type function",
				"The console-only \"type\" function cannot be used as part of an expression.",
			))
			return "", diags
		}
	}

	return FormatValue(val, 0), diags
}

func (s *Session) handleHelp() (string, tfdiags.Diagnostics) {
	text := `
The Terraform console allows you to experiment with Terraform interpolations.
You may access resources in the state (if you have one) just as you would
from a configuration. For example: "aws_instance.foo.id" would evaluate
to the ID of "aws_instance.foo" if it exists in your state.

Type in the interpolation to test and hit <enter> to see the result.

To exit the console, type "exit" and hit <enter>, or use Control-C or
Control-D.
`

	return strings.TrimSpace(text), nil
}
