package moduletest

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"

	"github.com/hashicorp/terraform/internal/addrs"
	"github.com/hashicorp/terraform/internal/configs"
	"github.com/hashicorp/terraform/internal/configs/configschema"
	"github.com/hashicorp/terraform/internal/plans"
	"github.com/hashicorp/terraform/internal/providers"
	"github.com/hashicorp/terraform/internal/states"
	"github.com/hashicorp/terraform/internal/tfdiags"
)

type Run struct {
	Config *configs.TestRun

	Verbose *Verbose

	Name   string
	Index  int
	Status Status

	Diagnostics tfdiags.Diagnostics
}

// Verbose is a utility struct that holds all the information required for a run
// to render the results verbosely.
//
// At the moment, this basically means printing out the plan. To do that we need
// all the information within this struct.
type Verbose struct {
	Plan         *plans.Plan
	State        *states.State
	Config       *configs.Config
	Providers    map[addrs.Provider]providers.ProviderSchema
	Provisioners map[string]*configschema.Block
}

func (run *Run) GetTargets() ([]addrs.Targetable, tfdiags.Diagnostics) {
	var diagnostics tfdiags.Diagnostics
	var targets []addrs.Targetable

	for _, target := range run.Config.Options.Target {
		addr, diags := addrs.ParseTarget(target)
		diagnostics = diagnostics.Append(diags)
		if addr != nil {
			targets = append(targets, addr.Subject)
		}
	}

	return targets, diagnostics
}

func (run *Run) GetReplaces() ([]addrs.AbsResourceInstance, tfdiags.Diagnostics) {
	var diagnostics tfdiags.Diagnostics
	var replaces []addrs.AbsResourceInstance

	for _, replace := range run.Config.Options.Replace {
		addr, diags := addrs.ParseAbsResourceInstance(replace)
		diagnostics = diagnostics.Append(diags)
		if diags.HasErrors() {
			continue
		}

		if addr.Resource.Resource.Mode != addrs.ManagedResourceMode {
			diagnostics = diagnostics.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Can only target managed resources for forced replacements.",
				Detail:   addr.String(),
				Subject:  replace.SourceRange().Ptr(),
			})
			continue
		}

		replaces = append(replaces, addr)
	}

	return replaces, diagnostics
}

func (run *Run) GetReferences() ([]*addrs.Reference, tfdiags.Diagnostics) {
	var diagnostics tfdiags.Diagnostics
	var references []*addrs.Reference

	for _, rule := range run.Config.CheckRules {
		for _, variable := range rule.Condition.Variables() {
			reference, diags := addrs.ParseRef(variable)
			diagnostics = diagnostics.Append(diags)
			if reference != nil {
				references = append(references, reference)
			}
		}
		for _, variable := range rule.ErrorMessage.Variables() {
			reference, diags := addrs.ParseRef(variable)
			diagnostics = diagnostics.Append(diags)
			if reference != nil {
				references = append(references, reference)
			}
		}
	}

	return references, diagnostics
}

// ValidateExpectedFailures steps through the provided diagnostics (which should
// be the result of a plan or an apply operation), and does 3 things:
//  1. Removes diagnostics that match the expected failures from the config.
//  2. Upgrades warnings from check blocks into errors where appropriate so the
//     test will fail later.
//  3. Adds diagnostics for any expected failures that were not satisfied.
//
// Point 2 is a bit complicated so worth expanding on. In normal Terraform
// execution, any error that originates within a check block (either from an
// assertion or a scoped data source) is wrapped up as a Warning to be
// identified to the user but not to fail the actual Terraform operation. During
// test execution, we want to upgrade (or rollback) these warnings into errors
// again so the test will fail. We do that as part of this function as we are
// already processing the diagnostics from check blocks in here anyway.
//
// The way the function works out which diagnostics are relevant to expected
// failures is by using the tfdiags Extra functionality to detect which
// diagnostics were generated by custom conditions. Terraform adds the
// addrs.CheckRule that generated each diagnostic to the diagnostic itself so we
// can tell which diagnostics can be expected.
func (run *Run) ValidateExpectedFailures(originals tfdiags.Diagnostics) tfdiags.Diagnostics {

	// We're going to capture all the checkable objects that are referenced
	// from the expected failures.
	expectedFailures := addrs.MakeMap[addrs.Referenceable, bool]()
	sourceRanges := addrs.MakeMap[addrs.Referenceable, tfdiags.SourceRange]()

	for _, traversal := range run.Config.ExpectFailures {
		// Ignore the diagnostics returned from the reference parsing, these
		// references will have been checked earlier in the process by the
		// validate stage so we don't need to do that again here.
		reference, _ := addrs.ParseRefFromTestingScope(traversal)
		expectedFailures.Put(reference.Subject, false)
		sourceRanges.Put(reference.Subject, reference.SourceRange)
	}

	var diags tfdiags.Diagnostics
	for _, diag := range originals {

		if rule, ok := addrs.DiagnosticOriginatesFromCheckRule(diag); ok {
			switch rule.Container.CheckableKind() {
			case addrs.CheckableOutputValue:
				addr := rule.Container.(addrs.AbsOutputValue)
				if !addr.Module.IsRoot() {
					// failures can only be expected against checkable objects
					// in the root module. This diagnostic will be added into
					// returned set below.
					break
				}

				if diag.Severity() == tfdiags.Warning {
					// Warnings don't count as errors. This diagnostic will be
					// added into the returned set below.
					break
				}

				if expectedFailures.Has(addr.OutputValue) {
					// Then this failure is expected! Mark the original map as
					// having found a failure and swallow this error by
					// continuing and not adding it into the returned set of
					// diagnostics.
					expectedFailures.Put(addr.OutputValue, true)
					continue
				}

				// Otherwise, this isn't an expected failure so just fall out
				// and add it into the returned set of diagnostics below.

			case addrs.CheckableInputVariable:
				addr := rule.Container.(addrs.AbsInputVariableInstance)
				if !addr.Module.IsRoot() {
					// failures can only be expected against checkable objects
					// in the root module. This diagnostic will be added into
					// returned set below.
					break
				}

				if diag.Severity() == tfdiags.Warning {
					// Warnings don't count as errors. This diagnostic will be
					// added into the returned set below.
					break
				}
				if expectedFailures.Has(addr.Variable) {
					// Then this failure is expected! Mark the original map as
					// having found a failure and swallow this error by
					// continuing and not adding it into the returned set of
					// diagnostics.
					expectedFailures.Put(addr.Variable, true)
					continue
				}

				// Otherwise, this isn't an expected failure so just fall out
				// and add it into the returned set of diagnostics below.

			case addrs.CheckableResource:
				addr := rule.Container.(addrs.AbsResourceInstance)
				if !addr.Module.IsRoot() {
					// failures can only be expected against checkable objects
					// in the root module. This diagnostic will be added into
					// returned set below.
					break
				}

				if diag.Severity() == tfdiags.Warning {
					// Warnings don't count as errors. This diagnostic will be
					// added into the returned set below.
					break
				}

				if expectedFailures.Has(addr.Resource) {
					// Then this failure is expected! Mark the original map as
					// having found a failure and swallow this error by
					// continuing and not adding it into the returned set of
					// diagnostics.
					expectedFailures.Put(addr.Resource, true)
					continue
				}

				if expectedFailures.Has(addr.Resource.Resource) {
					// We can also blanket expect failures in all instances for
					// a resource so we check for that here as well.
					expectedFailures.Put(addr.Resource.Resource, true)
					continue
				}

				// Otherwise, this isn't an expected failure so just fall out
				// and add it into the returned set of diagnostics below.

			case addrs.CheckableCheck:
				addr := rule.Container.(addrs.AbsCheck)

				// Check blocks are a bit more difficult than the others. Check
				// block diagnostics could be from a nested data block, or
				// from a failed assertion, and have all been marked as just
				// warning severity.
				//
				// For diagnostics from failed assertions, we want to check if
				// it was expected and skip it if it was. But if it wasn't
				// expected we want to upgrade the diagnostic from a warning
				// into an error so the test case will fail overall.
				//
				// For diagnostics from nested data blocks, we have two
				// categories of diagnostics. First, diagnostics that were
				// originally errors and we mapped into warnings. Second,
				// diagnostics that were originally warnings and stayed that
				// way. For the first case, we want to turn these back to errors
				// and use them as part of the expected failures functionality.
				// The second case should remain as warnings and be ignored by
				// the expected failures functionality.
				//
				// Note, as well that we still want to upgrade failed checks
				// from child modules into errors, so in the other branches we
				// just do a simple blanket skip off all diagnostics not
				// from the root module. We're more selective here, only
				// diagnostics from the root module are considered for the
				// expect failures functionality but we do also upgrade
				// diagnostics from child modules back into errors.

				if rule.Type == addrs.CheckAssertion {
					// Then this diagnostic is from a check block assertion, it
					// is something we want to treat as an error even though it
					// is actually claiming to be a warning.

					if addr.Module.IsRoot() && expectedFailures.Has(addr.Check) {
						// Then this failure is expected! Mark the original map as
						// having found a failure and continue.
						expectedFailures.Put(addr.Check, true)
						continue
					}

					// Otherwise, let's package this up as an error and move on.
					diags = diags.Append(tfdiags.Override(diag, tfdiags.Error, nil))
					continue
				} else if rule.Type == addrs.CheckDataResource {
					// Then the diagnostic we have was actually overridden so
					// let's get back to the original.
					original := tfdiags.UndoOverride(diag)

					// This diagnostic originated from a scoped data source.
					if addr.Module.IsRoot() && original.Severity() == tfdiags.Error {
						// Okay, we have a genuine error from the root module,
						// so we can now check if we want to ignore it or not.
						if expectedFailures.Has(addr.Check) {
							// Then this failure is expected! Mark the original map as
							// having found a failure and continue.
							expectedFailures.Put(addr.Check, true)
							continue
						}
					}

					// In all other cases, we want to add the original error
					// into the set we return to the testing framework and move
					// onto the next one.
					diags = diags.Append(original)
					continue
				} else {
					panic("invalid CheckType: " + rule.Type.String())
				}
			default:
				panic("unrecognized CheckableKind: " + rule.Container.CheckableKind().String())
			}
		}

		// If we get here, then we're not modifying the original diagnostic at
		// all. We just want the testing framework to treat it as normal.
		diags = diags.Append(diag)
	}

	// Okay, we've checked all our diagnostics to see if any were expected.
	// Now, let's make sure that all the checkable objects we expected to fail
	// actually did!

	for _, elem := range expectedFailures.Elems {
		addr := elem.Key
		failed := elem.Value

		if !failed {
			// Then we expected a failure, and it did not occur. Add it to the
			// diagnostics.
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Missing expected failure",
				Detail:   fmt.Sprintf("The checkable object, %s, was expected to report an error but did not.", addr.String()),
				Subject:  sourceRanges.Get(addr).ToHCL().Ptr(),
			})
		}
	}

	return diags
}
