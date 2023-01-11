package computed

import "github.com/hashicorp/terraform/internal/plans"

// Diff captures the computed diff for a single block, element or attribute.
//
// It essentially merges common functionality across all types of changes,
// namely the replace logic and the action / change type. Any remaining
// behaviour can be offloaded to the renderer which will be unique for the
// various change types (eg. maps, objects, lists, blocks, primitives, etc.).
type Diff struct {
	// Renderer captures the uncommon functionality across the different kinds
	// of changes. Each type of change (lists, blocks, sets, etc.) will have a
	// unique renderer.
	Renderer DiffRenderer

	// Action is the action described by this change (such as create, delete,
	// update, etc.).
	Action plans.Action

	// Replace tells the Change that it should add the `# forces replacement`
	// suffix.
	//
	// Every single change could potentially add this suffix, so we embed it in
	// the change as common functionality instead of in the specific renderers.
	Replace bool
}

// NewDiff creates a new Diff object with the provided renderer, action and
// replace context.
func NewDiff(renderer DiffRenderer, action plans.Action, replace bool) Diff {
	return Diff{
		Renderer: renderer,
		Action:   action,
		Replace:  replace,
	}
}

// RenderHuman prints the Change into a human-readable string referencing the
// specified RenderOpts.
//
// If the returned string is a single line, then indent should be ignored.
//
// If the return string is multiple lines, then indent should be used to offset
// the beginning of all lines but the first by the specified amount.
func (diff Diff) RenderHuman(indent int, opts RenderHumanOpts) string {
	return diff.Renderer.RenderHuman(diff, indent, opts)
}

// WarningsHuman returns a list of strings that should be rendered as warnings
// before a given change is rendered.
//
// As with the RenderHuman function, the indent should only be applied on
// multiline warnings and on the second and following lines.
func (diff Diff) WarningsHuman(indent int) []string {
	return diff.Renderer.WarningsHuman(diff, indent)
}

type DiffRenderer interface {
	RenderHuman(diff Diff, indent int, opts RenderHumanOpts) string
	WarningsHuman(diff Diff, indent int) []string
}

// RenderHumanOpts contains options that can control how the human render
// function of the DiffRenderer will function.
type RenderHumanOpts struct {
	// OverrideNullSuffix tells the Renderer not to display the `-> null` suffix
	// that is normally displayed when an element, attribute, or block is
	// deleted.
	OverrideNullSuffix bool

	// OverrideForcesReplacement tells the Renderer to display the
	// `# forces replacement` suffix, even if a diff doesn't have the Replace
	// field set.
	//
	// Some renderers (like the Set renderer) don't display the suffix
	// themselves but force their child diffs to display it instead.
	OverrideForcesReplacement bool

	// ShowUnchangedChildren instructs the Renderer to render all children of a
	// given complex change, instead of hiding unchanged items and compressing
	// them into a single line.
	ShowUnchangedChildren bool
}

// Clone returns a new RenderOpts object, that matches the original but can be
// edited without changing the original.
func (opts RenderHumanOpts) Clone() RenderHumanOpts {
	return RenderHumanOpts{
		OverrideNullSuffix:    opts.OverrideNullSuffix,
		ShowUnchangedChildren: opts.ShowUnchangedChildren,

		// OverrideForcesReplacement is a special case in that it doesn't
		// cascade. So each diff should decide independently whether it's direct
		// children should override their internal Replace logic, instead of
		// an ancestor making the switch and affecting the entire tree.
		OverrideForcesReplacement: false,
	}
}
