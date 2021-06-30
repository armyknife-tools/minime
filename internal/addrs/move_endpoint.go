package addrs

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/terraform/internal/tfdiags"
)

// MoveEndpoint is to AbsMoveable and ConfigMoveable what Target is to
// Targetable: a wrapping struct that captures the result of decoding an HCL
// traversal representing a relative path from the current module to
// a moveable object.
//
// Its name reflects that its primary purpose is for the "from" and "to"
// addresses in a "moved" statement in the configuration, but it's also
// valid to use MoveEndpoint for other similar mechanisms that give
// Terraform hints about historical configuration changes that might
// prompt creating a different plan than Terraform would by default.
//
// To obtain a full address from a MoveEndpoint you must use
// either the package function UnifyMoveEndpoints (to get an AbsMovable) or
// the method ConfigMoveable (to get a ConfigMoveable).
type MoveEndpoint struct {
	// SourceRange is the location of the physical endpoint address
	// in configuration, if this MoveEndpoint was decoded from a
	// configuration expresson.
	SourceRange tfdiags.SourceRange

	// Internally we (ab)use AbsMovable as the representation of our
	// relative address, even though everywhere else in Terraform
	// AbsMovable always represents a fully-absolute address.
	// In practice, due to the implementation of ParseMoveEndpoint,
	// this is always either a ModuleInstance or an AbsResourceInstance,
	// and we only consider the possibility of interpreting it as
	// a AbsModuleCall or an AbsResource in UnifyMoveEndpoints.
	// This is intentionally unexported to encapsulate this unusual
	// meaning of AbsMovable.
	relSubject AbsMoveable
}

func (e *MoveEndpoint) String() string {
	// Our internal pseudo-AbsMovable representing the relative
	// address (either ModuleInstance or AbsResourceInstance) is
	// a good enough proxy for the relative move endpoint address
	// serialization.
	return e.relSubject.String()
}

// ConfigMovable transforms the reciever into a ConfigMovable by resolving it
// relative to the given base module, which should be the module where
// the MoveEndpoint expression was found.
//
// The result is useful for finding the target object in the configuration,
// but it's not sufficient for fully interpreting a move statement because
// it lacks the specific module and resource instance keys.
func (e *MoveEndpoint) ConfigMoveable(baseModule Module) ConfigMoveable {
	addr := e.relSubject
	switch addr := addr.(type) {
	case ModuleInstance:
		ret := make(Module, 0, len(baseModule)+len(addr))
		ret = append(ret, baseModule...)
		ret = append(ret, addr.Module()...)
		return ret
	case AbsResourceInstance:
		moduleAddr := make(Module, 0, len(baseModule)+len(addr.Module))
		moduleAddr = append(moduleAddr, baseModule...)
		moduleAddr = append(moduleAddr, addr.Module.Module()...)
		return ConfigResource{
			Module:   moduleAddr,
			Resource: addr.Resource.Resource,
		}
	default:
		// The above should be exhaustive for all of the types
		// that ParseMoveEndpoint produces as our intermediate
		// address representation.
		panic(fmt.Sprintf("unsupported address type %T", addr))
	}

}

// ParseMoveEndpoint attempts to interpret the given traversal as a
// "move endpoint" address, which is a relative path from the module containing
// the traversal to a movable object in either the same module or in some
// child module.
//
// This deals only with the syntactic element of a move endpoint expression
// in configuration. Before the result will be useful you'll need to combine
// it with the address of the module where it was declared in order to get
// an absolute address relative to the root module.
func ParseMoveEndpoint(traversal hcl.Traversal) (*MoveEndpoint, tfdiags.Diagnostics) {
	path, remain, diags := parseModuleInstancePrefix(traversal)
	if diags.HasErrors() {
		return nil, diags
	}

	rng := tfdiags.SourceRangeFromHCL(traversal.SourceRange())

	if len(remain) == 0 {
		return &MoveEndpoint{
			relSubject:  path,
			SourceRange: rng,
		}, diags
	}

	riAddr, moreDiags := parseResourceInstanceUnderModule(path, remain)
	diags = diags.Append(moreDiags)
	if diags.HasErrors() {
		return nil, diags
	}

	return &MoveEndpoint{
		relSubject:  riAddr,
		SourceRange: rng,
	}, diags
}

// UnifyMoveEndpoints takes a pair of MoveEndpoint objects representing the
// "from" and "to" addresses in a moved block, and returns a pair of
// AbsMoveable addresses guaranteed to be of the same dynamic type
// that represent what the two MoveEndpoint addresses refer to.
//
// moduleAddr must be the address of the module instance where the move
// was declared.
//
// This function deals both with the conversion from relative to absolute
// addresses and with resolving the ambiguity between no-key instance
// addresses and whole-object addresses, returning the least specific
// address type possible.
//
// Not all combinations of addresses are unifyable: the two addresses must
// either both include resources or both just be modules. If the two
// given addresses are incompatible then UnifyMoveEndpoints returns (nil, nil),
// in which case the caller should typically report an error to the user
// stating the unification constraints.
func UnifyMoveEndpoints(moduleAddr ModuleInstance, relFrom, relTo *MoveEndpoint) (absFrom, absTo AbsMoveable) {

	// First we'll make a decision about which address type we're
	// ultimately trying to unify to. For our internal purposes
	// here we're going to borrow TargetableAddrType just as a
	// convenient way to talk about our address types, even though
	// targetable address types are not 100% aligned with moveable
	// address types.
	fromType := relFrom.internalAddrType()
	toType := relTo.internalAddrType()
	var wantType TargetableAddrType

	// Our goal here is to choose the whole-resource or whole-module-call
	// addresses if both agree on it, but to use specific instance addresses
	// otherwise. This is a somewhat-arbitrary way to resolve syntactic
	// ambiguity between the two situations which allows both for renaming
	// whole resources and for switching from a single-instance object to
	// a multi-instance object.
	switch {
	case fromType == AbsResourceInstanceAddrType || toType == AbsResourceInstanceAddrType:
		wantType = AbsResourceInstanceAddrType
	case fromType == AbsResourceAddrType || toType == AbsResourceAddrType:
		wantType = AbsResourceAddrType
	case fromType == ModuleInstanceAddrType || toType == ModuleInstanceAddrType:
		wantType = ModuleInstanceAddrType
	case fromType == ModuleAddrType || toType == ModuleAddrType:
		// NOTE: We're fudging a little here and using
		// ModuleAddrType to represent AbsModuleCall rather
		// than Module.
		wantType = ModuleAddrType
	default:
		panic("unhandled move address types")
	}

	absFrom = relFrom.prepareAbsMoveable(moduleAddr, wantType)
	absTo = relTo.prepareAbsMoveable(moduleAddr, wantType)
	if absFrom == nil || absTo == nil {
		// if either of them failed then they both failed, to make the
		// caller's life a little easier.
		return nil, nil
	}
	return absFrom, absTo
}

func (e *MoveEndpoint) prepareAbsMoveable(moduleAddr ModuleInstance, wantType TargetableAddrType) AbsMoveable {
	// relAddr can only be either AbsResourceInstance or ModuleInstance, the
	// internal intermediate representation produced by ParseMoveEndpoint.
	relAddr := e.relSubject

	switch relAddr := relAddr.(type) {
	case ModuleInstance:
		switch wantType {
		case ModuleInstanceAddrType:
			ret := make(ModuleInstance, 0, len(moduleAddr)+len(relAddr))
			ret = append(ret, moduleAddr...)
			ret = append(ret, relAddr...)
			return ret
		case ModuleAddrType:
			// NOTE: We're fudging a little here and using
			// ModuleAddrType to represent AbsModuleCall rather
			// than Module.
			callerAddr := make(ModuleInstance, 0, len(moduleAddr)+len(relAddr)-1)
			callerAddr = append(callerAddr, moduleAddr...)
			callerAddr = append(callerAddr, relAddr[:len(relAddr)-1]...)
			return AbsModuleCall{
				Module: callerAddr,
				Call: ModuleCall{
					Name: relAddr[len(relAddr)-1].Name,
				},
			}
		default:
			return nil // can't make any other types from a ModuleInstance
		}
	case AbsResourceInstance:
		callerAddr := make(ModuleInstance, 0, len(moduleAddr)+len(relAddr.Module))
		callerAddr = append(callerAddr, moduleAddr...)
		callerAddr = append(callerAddr, relAddr.Module...)
		switch wantType {
		case AbsResourceInstanceAddrType:
			return AbsResourceInstance{
				Module:   callerAddr,
				Resource: relAddr.Resource,
			}
		case AbsResourceAddrType:
			return AbsResource{
				Module:   callerAddr,
				Resource: relAddr.Resource.Resource,
			}
		default:
			return nil // can't make any other types from an AbsResourceInstance
		}
	default:
		panic(fmt.Sprintf("unhandled address type %T", relAddr))
	}
}

// internalAddrType helps facilitate our slight abuse of TargetableAddrType
// as a way to talk about our different possible result address types in
// UnifyMoveEndpoints.
//
// It's not really correct to use TargetableAddrType in this way, because
// it's for Targetable rather than for AbsMoveable, but as long as the two
// remain aligned enough it saves introducing yet another enumeration with
// similar members that would be for internal use only anyway.
func (e *MoveEndpoint) internalAddrType() TargetableAddrType {
	switch addr := e.relSubject.(type) {
	case ModuleInstance:
		if !addr.IsRoot() && addr[len(addr)-1].InstanceKey == NoKey {
			// NOTE: We're fudging a little here and using
			// ModuleAddrType to represent AbsModuleCall rather
			// than Module.
			return ModuleAddrType
		}
		return ModuleInstanceAddrType
	case AbsResourceInstance:
		if addr.Resource.Key == NoKey {
			return AbsResourceAddrType
		}
		return AbsResourceInstanceAddrType
	default:
		// The above should cover all of the address types produced
		// by ParseMoveEndpoint.
		panic(fmt.Sprintf("unsupported address type %T", addr))
	}
}
