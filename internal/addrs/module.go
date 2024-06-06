// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package addrs

import (
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/opentofu/opentofu/internal/tfdiags"
)

// Module is an address for a module call within configuration. This is
// the static counterpart of ModuleInstance, representing a traversal through
// the static module call tree in configuration and does not take into account
// the potentially-multiple instances of a module that might be created by
// "count" and "for_each" arguments within those calls.
//
// This type should be used only in very specialized cases when working with
// the static module call tree. Type ModuleInstance is appropriate in more cases.
//
// Although Module is a slice, it should be treated as immutable after creation.
type Module []string

// RootModule is the module address representing the root of the static module
// call tree, which is also the zero value of Module.
//
// Note that this is not the root of the dynamic module tree, which is instead
// represented by RootModuleInstance.
var RootModule Module

// IsRoot returns true if the receiver is the address of the root module,
// or false otherwise.
func (m Module) IsRoot() bool {
	return len(m) == 0
}

func (m Module) String() string {
	if len(m) == 0 {
		return ""
	}
	// Calculate necessary space.
	l := 0
	for _, step := range m {
		l += len(step)
	}
	buf := strings.Builder{}
	// 8 is len(".module.") which separates entries.
	buf.Grow(l + len(m)*8)
	sep := ""
	for _, step := range m {
		buf.WriteString(sep)
		buf.WriteString("module.")
		buf.WriteString(step)
		sep = "."
	}
	return buf.String()
}

func (m Module) Equal(other Module) bool {
	if len(m) != len(other) {
		return false
	}
	for i := range m {
		if m[i] != other[i] {
			return false
		}
	}
	return true
}

func (m Module) targetableSigil() {
	// Module is targetable
}

// TargetContains implements Targetable for Module by returning true if the given other
// address either matches the receiver, is a sub-module-instance of the
// receiver, or is a targetable absolute address within a module that
// is contained within the receiver.
func (m Module) TargetContains(other Targetable) bool {
	switch to := other.(type) {

	case Module:
		if len(to) < len(m) {
			// Can't be contained if the path is shorter
			return false
		}
		// Other is contained if its steps match for the length of our own path.
		for i, ourStep := range m {
			otherStep := to[i]
			if ourStep != otherStep {
				return false
			}
		}
		// If we fall out here then the prefixed matched, so it's contained.
		return true

	case ModuleInstance:
		return m.TargetContains(to.Module())

	case ConfigResource:
		return m.TargetContains(to.Module)

	case AbsResource:
		return m.TargetContains(to.Module)

	case AbsResourceInstance:
		return m.TargetContains(to.Module)

	default:
		return false
	}
}

func (m Module) AddrType() TargetableAddrType {
	return ModuleAddrType
}

// Child returns the address of a child call in the receiver, identified by the
// given name.
func (m Module) Child(name string) Module {
	ret := make(Module, 0, len(m)+1)
	ret = append(ret, m...)
	return append(ret, name)
}

// Parent returns the address of the parent module of the receiver, or the
// receiver itself if there is no parent (if it's the root module address).
func (m Module) Parent() Module {
	if len(m) == 0 {
		return m
	}
	return m[:len(m)-1]
}

// Call returns the module call address that corresponds to the given module
// instance, along with the address of the module that contains it.
//
// There is no call for the root module, so this method will panic if called
// on the root module address.
//
// In practice, this just turns the last element of the receiver into a
// ModuleCall and then returns a slice of the receiever that excludes that
// last part. This is just a convenience for situations where a call address
// is required, such as when dealing with *Reference and Referencable values.
func (m Module) Call() (Module, ModuleCall) {
	if len(m) == 0 {
		panic("cannot produce ModuleCall for root module")
	}

	caller, callName := m[:len(m)-1], m[len(m)-1]
	return caller, ModuleCall{
		Name: callName,
	}
}

// Ancestors returns a slice containing the receiver and all of its ancestor
// modules, all the way up to (and including) the root module.  The result is
// ordered by depth, with the root module always first.
//
// Since the result always includes the root module, a caller may choose to
// ignore it by slicing the result with [1:].
func (m Module) Ancestors() []Module {
	ret := make([]Module, 0, len(m)+1)
	for i := 0; i <= len(m); i++ {
		ret = append(ret, m[:i])
	}
	return ret
}

func (m Module) configMoveableSigil() {
	// ModuleInstance is moveable
}
func (m Module) configRemovableSigil() {
	// Empty function so Module will fulfill the requirements of the removable interface
}

// ParseModule parses a module address from the given traversal,
// which has to contain only the module address with no resource/data/variable/etc.
// This function only supports module addresses without instance keys (as the
// returned Module struct doesn't support instance keys) and will return an
// error if it encounters one.
func ParseModule(traversal hcl.Traversal) (Module, tfdiags.Diagnostics) {
	mod, remain, diags := parseModulePrefix(traversal)
	if !diags.HasErrors() && len(remain) != 0 {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Module address expected",
			Detail:   "It's not allowed to reference anything other than module here.",
			Subject:  remain[0].SourceRange().Ptr(),
		})
	}

	return mod, diags
}

// parseModulePrefix parses a module address from the given traversal,
// returning the module address and the remaining traversal.
// For example, if the input traversal is ["module","a","module","b",
// "null_resource", example_resource"], the output module will be ["a", "b"]
// and the output remaining traversal will be ["null_resource",
// "example_resource"].
// This function only supports module addresses without instance keys (as the
// returned Module struct doesn't support instance keys) and will return an
// error if it encounters one.
func parseModulePrefix(traversal hcl.Traversal) (Module, hcl.Traversal, tfdiags.Diagnostics) {
	remain := traversal
	var module Module
	var diags tfdiags.Diagnostics

	for len(remain) > 0 {
		moduleName, isModule, moduleNameDiags := getModuleName(remain)
		diags = diags.Append(moduleNameDiags)

		if !isModule || diags.HasErrors() {
			break
		}

		// Because this is a valid module address, we can safely assume that
		// the first two elements are "module" and the module name
		remain = remain[2:]

		if len(remain) > 0 {
			// We don't allow module instances as part of the module address
			if _, ok := remain[0].(hcl.TraverseIndex); ok {
				diags = diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Module instance address with keys is not allowed",
					Detail:   "Module address cannot be a module instance (e.g. \"module.a[0]\"), it must be a module instead (e.g. \"module.a\").",
					Subject:  remain[0].SourceRange().Ptr(),
				})

				return module, remain, diags
			}
		}

		module = append(module, moduleName)
	}

	var retRemain hcl.Traversal
	if len(remain) > 0 {
		retRemain = make(hcl.Traversal, len(remain))
		copy(retRemain, remain)
		// The first element here might be either a TraverseRoot or a
		// TraverseAttr, depending on whether we had a module address on the
		// front. To make life easier for callers, we'll normalize to always
		// start with a TraverseRoot.
		if tt, ok := retRemain[0].(hcl.TraverseAttr); ok {
			retRemain[0] = hcl.TraverseRoot{
				Name:     tt.Name,
				SrcRange: tt.SrcRange,
			}
		}
	}

	return module, retRemain, diags
}

func getModuleName(remain hcl.Traversal) (moduleName string, isModule bool, diags tfdiags.Diagnostics) {
	if len(remain) == 0 {
		// If the address is empty, then we can't possibly have a module address
		return moduleName, false, diags
	}

	var next string
	switch tt := remain[0].(type) {
	case hcl.TraverseRoot:
		next = tt.Name
	case hcl.TraverseAttr:
		next = tt.Name
	default:
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid address operator",
			Detail:   "Module address prefix must be followed by dot and then a name.",
			Subject:  remain[0].SourceRange().Ptr(),
		})

		return moduleName, false, diags
	}

	if next != "module" {
		return moduleName, false, diags
	}

	kwRange := remain[0].SourceRange()
	remain = remain[1:]
	// If we have the prefix "module" then we should be followed by a
	// module call name, as an attribute
	if len(remain) == 0 {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid address operator",
			Detail:   "Prefix \"module.\" must be followed by a module name.",
			Subject:  &kwRange,
		})

		return moduleName, false, diags
	}

	switch tt := remain[0].(type) {
	case hcl.TraverseAttr:
		moduleName = tt.Name
	default:
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid address operator",
			Detail:   "Prefix \"module.\" must be followed by a module name.",
			Subject:  remain[0].SourceRange().Ptr(),
		})
		return moduleName, false, diags
	}
	return moduleName, true, diags
}
