package depsfile

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/internal/getproviders"
	"github.com/hashicorp/terraform/tfdiags"
	"github.com/hashicorp/terraform/version"
)

// LoadLocksFromFile reads locks from the given file, expecting it to be a
// valid dependency lock file, or returns error diagnostics explaining why
// that was not possible.
//
// The returned locks are a snapshot of what was present on disk at the time
// the method was called. It does not take into account any subsequent writes
// to the file, whether through this package's functions or by external
// writers.
//
// If the returned diagnostics contains errors then the returned Locks may
// be incomplete or invalid.
func LoadLocksFromFile(filename string) (*Locks, tfdiags.Diagnostics) {
	ret := NewLocks()

	var diags tfdiags.Diagnostics

	parser := hclparse.NewParser()
	f, hclDiags := parser.ParseHCLFile(filename)
	ret.sources = parser.Sources()
	diags = diags.Append(hclDiags)

	moreDiags := decodeLocksFromHCL(ret, f.Body)
	diags = diags.Append(moreDiags)
	return ret, diags
}

// SaveLocksToFile writes the given locks object to the given file,
// entirely replacing any content already in that file, or returns error
// diagnostics explaining why that was not possible.
//
// SaveLocksToFile attempts an atomic replacement of the file, as an aid
// to external tools such as text editor integrations that might be monitoring
// the file as a signal to invalidate cached metadata. Consequently, other
// temporary files may be temporarily created in the same directory as the
// given filename during the operation.
func SaveLocksToFile(locks *Locks, filename string) tfdiags.Diagnostics {
	var diags tfdiags.Diagnostics

	// In other uses of the "hclwrite" package we typically try to make
	// surgical updates to the author's existing files, preserving their
	// block ordering, comments, etc. We intentionally don't do that here
	// to reinforce the fact that this file primarily belongs to Terraform,
	// and to help ensure that VCS diffs of the file primarily reflect
	// changes that actually affect functionality rather than just cosmetic
	// changes, by maintaining it in a highly-normalized form.

	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()

	// End-users _may_ edit the lock file in exceptional situations, like
	// working around potential dependency selection bugs, but we intend it
	// to be primarily maintained automatically by the "terraform init"
	// command.
	rootBody.AppendUnstructuredTokens(hclwrite.Tokens{
		{
			Type:  hclsyntax.TokenComment,
			Bytes: []byte("# This file is maintained automatically by \"terraform init\".\n"),
		},
		{
			Type:  hclsyntax.TokenComment,
			Bytes: []byte("# Manual edits may be lost in future updates.\n"),
		},
	})

	providers := make([]addrs.Provider, 0, len(locks.providers))
	for provider := range locks.providers {
		providers = append(providers, provider)
	}
	sort.Slice(providers, func(i, j int) bool {
		return providers[i].LessThan(providers[j])
	})

	for _, provider := range providers {
		lock := locks.providers[provider]
		rootBody.AppendNewline()
		block := rootBody.AppendNewBlock("provider", []string{lock.addr.String()})
		body := block.Body()
		body.SetAttributeValue("version", cty.StringVal(lock.version.String()))
		if constraintsStr := getproviders.VersionConstraintsString(lock.versionConstraints); constraintsStr != "" {
			body.SetAttributeValue("constraints", cty.StringVal(constraintsStr))
		}
		if len(lock.hashes) != 0 {
			hashVals := make([]cty.Value, 0, len(lock.hashes))
			for _, str := range lock.hashes {
				hashVals = append(hashVals, cty.StringVal(str))
			}
			// We're using a set rather than a list here because the order
			// isn't significant and SetAttributeValue will automatically
			// write the set elements in a consistent lexical order.
			hashSet := cty.SetVal(hashVals)
			body.SetAttributeValue("hashes", hashSet)
		}
	}

	newContent := f.Bytes()

	// TODO: Create the content in a new file and atomically pivot it into
	// the target, so that there isn't a brief period where an incomplete
	// file can be seen at the given location.
	// But for now, this gets us started.
	err := ioutil.WriteFile(filename, newContent, os.ModePerm)
	if err != nil {
		diags = diags.Append(tfdiags.Sourceless(
			tfdiags.Error,
			"Failed to update dependency lock file",
			fmt.Sprintf("Error while writing new dependency lock information to %s: %s.", filename, err),
		))
		return diags
	}

	return diags
}

func decodeLocksFromHCL(locks *Locks, body hcl.Body) tfdiags.Diagnostics {
	var diags tfdiags.Diagnostics

	content, hclDiags := body.Content(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "provider",
				LabelNames: []string{"source_addr"},
			},

			// "module" is just a placeholder for future enhancement, so we
			// can mostly-ignore the this block type we intend to add in
			// future, but warn in case someone tries to use one e.g. if they
			// downgraded to an earlier version of Terraform.
			{
				Type:       "module",
				LabelNames: []string{"path"},
			},
		},
	})
	diags = diags.Append(hclDiags)

	seenProviders := make(map[addrs.Provider]hcl.Range)
	seenModule := false
	for _, block := range content.Blocks {

		switch block.Type {
		case "provider":
			lock, moreDiags := decodeProviderLockFromHCL(block)
			diags = diags.Append(moreDiags)
			if lock == nil {
				continue
			}
			if previousRng, exists := seenProviders[lock.addr]; exists {
				diags = diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Duplicate provider lock",
					Detail:   fmt.Sprintf("This lockfile already declared a lock for provider %s at %s.", lock.addr.String(), previousRng.String()),
					Subject:  block.TypeRange.Ptr(),
				})
				continue
			}
			locks.providers[lock.addr] = lock
			seenProviders[lock.addr] = block.DefRange

		case "module":
			// We'll just take the first module block to use for a single warning,
			// because that's sufficient to get the point across without swamping
			// the output with warning noise.
			if !seenModule {
				currentVersion := version.SemVer.String()
				diags = diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagWarning,
					Summary:  "Dependency locks for modules are not yet supported",
					Detail:   fmt.Sprintf("Terraform v%s only supports dependency locks for providers, not for modules. This configuration may be intended for a later version of Terraform that also supports dependency locks for modules.", currentVersion),
					Subject:  block.TypeRange.Ptr(),
				})
				seenModule = true
			}

		default:
			// Shouldn't get here because this should be exhaustive for
			// all of the block types in the schema above.
		}

	}

	return diags
}

func decodeProviderLockFromHCL(block *hcl.Block) (*ProviderLock, tfdiags.Diagnostics) {
	ret := &ProviderLock{}
	var diags tfdiags.Diagnostics

	rawAddr := block.Labels[0]
	addr, moreDiags := addrs.ParseProviderSourceString(rawAddr)
	if moreDiags.HasErrors() {
		// The diagnostics from ParseProviderSourceString are, as the name
		// suggests, written with an intended audience of someone who is
		// writing a "source" attribute in a provider requirement, not
		// our lock file. Therefore we're using a less helpful, fixed error
		// here, which is non-ideal but hopefully okay for now because we
		// don't intend end-users to typically be hand-editing these anyway.
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid provider source address",
			Detail:   "The provider source address for a provider lock must be a valid, fully-qualified address of the form \"hostname/namespace/type\".",
			Subject:  block.LabelRanges[0].Ptr(),
		})
		return nil, diags
	}
	if !ProviderIsLockable(addr) {
		if addr.IsBuiltIn() {
			// A specialized error for built-in providers, because we have an
			// explicit explanation for why those are not allowed.
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid provider source address",
				Detail:   fmt.Sprintf("Cannot lock a version for built-in provider %s. Built-in providers are bundled inside Terraform itself, so you can't select a version for them independently of the Terraform release you are currently running.", addr),
				Subject:  block.LabelRanges[0].Ptr(),
			})
			return nil, diags
		}
		// Otherwise, we'll use a generic error message.
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid provider source address",
			Detail:   fmt.Sprintf("Provider source address %s is a special provider that is not eligible for dependency locking.", addr),
			Subject:  block.LabelRanges[0].Ptr(),
		})
		return nil, diags
	}
	if canonAddr := addr.String(); canonAddr != rawAddr {
		// We also require the provider addresses in the lock file to be
		// written in fully-qualified canonical form, so that it's totally
		// clear to a reader which provider each block relates to. Again,
		// we expect hand-editing of these to be atypical so it's reasonable
		// to be stricter in parsing these than we would be in the main
		// configuration.
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Non-normalized provider source address",
			Detail:   fmt.Sprintf("The provider source address for this provider lock must be written as %q, the fully-qualified and normalized form.", canonAddr),
			Subject:  block.LabelRanges[0].Ptr(),
		})
		return nil, diags
	}

	ret.addr = addr

	content, hclDiags := block.Body.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "version", Required: true},
			{Name: "constraints"},
			{Name: "hashes"},
		},
	})
	diags = diags.Append(hclDiags)

	version, moreDiags := decodeProviderVersionArgument(addr, content.Attributes["version"])
	ret.version = version
	diags = diags.Append(moreDiags)

	constraints, moreDiags := decodeProviderVersionConstraintsArgument(addr, content.Attributes["constraints"])
	ret.versionConstraints = constraints
	diags = diags.Append(moreDiags)

	hashes, moreDiags := decodeProviderHashesArgument(addr, content.Attributes["hashes"])
	ret.hashes = hashes
	diags = diags.Append(moreDiags)

	return ret, diags
}

func decodeProviderVersionArgument(provider addrs.Provider, attr *hcl.Attribute) (getproviders.Version, tfdiags.Diagnostics) {
	var diags tfdiags.Diagnostics
	if attr == nil {
		// It's not okay to omit this argument, but the caller should already
		// have generated diagnostics about that.
		return getproviders.UnspecifiedVersion, diags
	}
	expr := attr.Expr

	var raw *string
	hclDiags := gohcl.DecodeExpression(expr, nil, &raw)
	diags = diags.Append(hclDiags)
	if hclDiags.HasErrors() {
		return getproviders.UnspecifiedVersion, diags
	}
	if raw == nil {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Missing required argument",
			Detail:   "A provider lock block must contain a \"version\" argument.",
			Subject:  expr.Range().Ptr(), // the range for a missing argument's expression is the body's missing item range
		})
		return getproviders.UnspecifiedVersion, diags
	}
	version, err := getproviders.ParseVersion(*raw)
	if err != nil {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid provider version number",
			Detail:   fmt.Sprintf("The selected version number for provider %s is invalid: %s.", provider, err),
			Subject:  expr.Range().Ptr(),
		})
	}
	if canon := version.String(); canon != *raw {
		// Canonical forms are required in the lock file, to reduce the risk
		// that a file diff will show changes that are entirely cosmetic.
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid provider version number",
			Detail:   fmt.Sprintf("The selected version number for provider %s must be written in normalized form: %q.", provider, canon),
			Subject:  expr.Range().Ptr(),
		})
	}
	return version, diags
}

func decodeProviderVersionConstraintsArgument(provider addrs.Provider, attr *hcl.Attribute) (getproviders.VersionConstraints, tfdiags.Diagnostics) {
	var diags tfdiags.Diagnostics
	if attr == nil {
		// It's okay to omit this argument.
		return nil, diags
	}
	expr := attr.Expr

	var raw string
	hclDiags := gohcl.DecodeExpression(expr, nil, &raw)
	diags = diags.Append(hclDiags)
	if hclDiags.HasErrors() {
		return nil, diags
	}
	constraints, err := getproviders.ParseVersionConstraints(raw)
	if err != nil {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid provider version constraints",
			Detail:   fmt.Sprintf("The recorded version constraints for provider %s are invalid: %s.", provider, err),
			Subject:  expr.Range().Ptr(),
		})
	}
	if canon := getproviders.VersionConstraintsString(constraints); canon != raw {
		// Canonical forms are required in the lock file, to reduce the risk
		// that a file diff will show changes that are entirely cosmetic.
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid provider version constraints",
			Detail:   fmt.Sprintf("The recorded version constraints for provider %s must be written in normalized form: %q.", provider, canon),
			Subject:  expr.Range().Ptr(),
		})
	}

	return constraints, diags
}

func decodeProviderHashesArgument(provider addrs.Provider, attr *hcl.Attribute) ([]string, tfdiags.Diagnostics) {
	var diags tfdiags.Diagnostics
	if attr == nil {
		// It's okay to omit this argument.
		return nil, diags
	}
	expr := attr.Expr

	// We'll decode this argument using the HCL static analysis mode, because
	// there's no reason for the hashes list to be dynamic and this way we can
	// give more precise feedback on individual elements that are invalid,
	// with direct source locations.
	hashExprs, hclDiags := hcl.ExprList(expr)
	diags = diags.Append(hclDiags)
	if hclDiags.HasErrors() {
		return nil, diags
	}
	if len(hashExprs) == 0 {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid provider hash set",
			Detail:   "The \"hashes\" argument must either be omitted or contain at least one hash value.",
			Subject:  expr.Range().Ptr(),
		})
		return nil, diags
	}

	ret := make([]string, 0, len(hashExprs))
	for _, hashExpr := range hashExprs {
		var raw string
		hclDiags := gohcl.DecodeExpression(hashExpr, nil, &raw)
		diags = diags.Append(hclDiags)
		if hclDiags.HasErrors() {
			continue
		}
		// TODO: Validate the hash syntax, but not the actual hash schemes
		// because we expect to support different hash formats over time and
		// will silently ignore ones that we no longer prefer.
		ret = append(ret, raw)
	}

	return ret, diags
}
