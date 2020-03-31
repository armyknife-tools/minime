package getproviders

import (
	"fmt"
	"strings"

	svchost "github.com/hashicorp/terraform-svchost"

	"github.com/hashicorp/terraform/addrs"
)

// MultiSource is a Source that wraps a series of other sources and combines
// their sets of available providers and provider versions.
//
// A MultiSource consists of a sequence of selectors that each specify an
// underlying source to query and a set of matching patterns to decide which
// providers can be retrieved from which sources. If multiple selectors find
// a given provider version then the earliest one in the sequence takes
// priority for deciding the package metadata for the provider.
//
// For underlying sources that make network requests, consider wrapping each
// one in a MemoizeSource so that availability information retrieved in
// AvailableVersions can be reused in PackageMeta.
type MultiSource []MultiSourceSelector

var _ Source = MultiSource(nil)

// AvailableVersions retrieves all of the versions of the given provider
// that are available across all of the underlying selectors, while respecting
// each selector's matching patterns.
func (s MultiSource) AvailableVersions(provider addrs.Provider) (VersionList, error) {
	if len(s) == 0 { // Easy case: there can be no available versions
		return nil, nil
	}

	// We will return the union of all versions reported by the nested
	// sources that have matching patterns that accept the given provider.
	vs := make(map[Version]struct{})
	for _, selector := range s {
		if !selector.CanHandleProvider(provider) {
			continue // doesn't match the given patterns
		}
		thisSourceVersions, err := selector.Source.AvailableVersions(provider)
		switch err.(type) {
		case nil:
			// okay
		case ErrProviderNotKnown:
			continue // ignore, then
		default:
			return nil, err
		}
		for _, v := range thisSourceVersions {
			vs[v] = struct{}{}
		}
	}

	if len(vs) == 0 {
		return nil, ErrProviderNotKnown{provider}
	}
	ret := make(VersionList, 0, len(vs))
	for v := range vs {
		ret = append(ret, v)
	}
	ret.Sort()

	return ret, nil
}

// PackageMeta retrieves the package metadata for the requested provider package
// from the first selector that indicates availability of it.
func (s MultiSource) PackageMeta(provider addrs.Provider, version Version, target Platform) (PackageMeta, error) {
	if len(s) == 0 { // Easy case: no providers exist at all
		return PackageMeta{}, ErrProviderNotKnown{provider}
	}

	for _, selector := range s {
		if !selector.CanHandleProvider(provider) {
			continue // doesn't match the given patterns
		}
		meta, err := selector.Source.PackageMeta(provider, version, target)
		switch err.(type) {
		case nil:
			return meta, nil
		case ErrProviderNotKnown, ErrPlatformNotSupported:
			continue // ignore, then
		default:
			return PackageMeta{}, err
		}
	}

	// If we fall out here then none of the sources have the requested
	// package.
	return PackageMeta{}, ErrPlatformNotSupported{
		Provider: provider,
		Version:  version,
		Platform: target,
	}
}

// MultiSourceSelector is an element of the source selection configuration on
// MultiSource. A MultiSource has zero or more of these to configure which
// underlying sources it should consult for a given provider.
type MultiSourceSelector struct {
	// Source is the underlying source that this selector applies to.
	Source Source

	// Include and Exclude are sets of provider matching patterns that
	// together define which providers are eligible to be potentially
	// installed from the corresponding Source.
	Include, Exclude MultiSourceMatchingPatterns
}

// MultiSourceMatchingPatterns is a set of patterns that together define a
// set of providers by matching on the segments of the provider FQNs.
//
// The Provider address values in a MultiSourceMatchingPatterns are special in
// that any of Hostname, Namespace, or Type can be getproviders.Wildcard
// to indicate that any concrete value is permitted for that segment.
type MultiSourceMatchingPatterns []addrs.Provider

// ParseMultiSourceMatchingPatterns parses a slice of strings containing the
// string form of provider matching patterns and, if all the given strings
// are valid, returns the corresponding MultiSourceMatchingPatterns value.
func ParseMultiSourceMatchingPatterns(strs []string) (MultiSourceMatchingPatterns, error) {
	if len(strs) == 0 {
		return nil, nil
	}

	ret := make(MultiSourceMatchingPatterns, len(strs))
	for i, str := range strs {
		parts := strings.Split(str, "/")
		if len(parts) < 2 || len(parts) > 3 {
			return nil, fmt.Errorf("invalid provider matching pattern %q: must have either two or three slash-separated segments", str)
		}
		host := defaultRegistryHost
		explicitHost := len(parts) == 3
		if explicitHost {
			givenHost := parts[0]
			if givenHost == "*" {
				host = svchost.Hostname(Wildcard)
			} else {
				normalHost, err := svchost.ForComparison(givenHost)
				if err != nil {
					return nil, fmt.Errorf("invalid hostname in provider matching pattern %q: %s", str, err)
				}

				// The remaining code below deals only with the namespace/type portions.
				host = normalHost
			}

			parts = parts[1:]
		}

		if !validProviderNameOrWildcard(parts[1]) {
			return nil, fmt.Errorf("invalid provider type %q in provider matching pattern %q: must either be the wildcard * or a provider type name", parts[1], str)
		}
		if !validProviderNameOrWildcard(parts[0]) {
			return nil, fmt.Errorf("invalid registry namespace %q in provider matching pattern %q: must either be the wildcard * or a literal namespace", parts[1], str)
		}

		ret[i] = addrs.Provider{
			Hostname:  host,
			Namespace: parts[0],
			Type:      parts[1],
		}

		if ret[i].Hostname == svchost.Hostname(Wildcard) && !(ret[i].Namespace == Wildcard && ret[i].Type == Wildcard) {
			return nil, fmt.Errorf("invalid provider matching pattern %q: hostname can be a wildcard only if both namespace and provider type are also wildcards", str)
		}
		if ret[i].Namespace == Wildcard && ret[i].Type != Wildcard {
			return nil, fmt.Errorf("invalid provider matching pattern %q: namespace can be a wildcard only if the provider type is also a wildcard", str)
		}
	}
	return ret, nil
}

// CanHandleProvider returns true if and only if the given provider address
// is both included by the selector's include patterns and _not_ excluded
// by its exclude patterns.
//
// The absense of any include patterns is treated the same as a pattern
// that matches all addresses. Exclusions take priority over inclusions.
func (s MultiSourceSelector) CanHandleProvider(addr addrs.Provider) bool {
	switch {
	case s.Exclude.MatchesProvider(addr):
		return false
	case len(s.Include) > 0:
		return s.Include.MatchesProvider(addr)
	default:
		return true
	}
}

// MatchesProvider tests whether the receiving matching patterns match with
// the given concrete provider address.
func (ps MultiSourceMatchingPatterns) MatchesProvider(addr addrs.Provider) bool {
	for _, pattern := range ps {
		hostMatch := (pattern.Hostname == svchost.Hostname(Wildcard) || pattern.Hostname == addr.Hostname)
		namespaceMatch := (pattern.Namespace == Wildcard || pattern.Namespace == addr.Namespace)
		typeMatch := (pattern.Type == Wildcard || pattern.Type == addr.Type)
		if hostMatch && namespaceMatch && typeMatch {
			return true
		}
	}
	return false
}

// Wildcard is a string value representing a wildcard element in the Include
// and Exclude patterns used with MultiSource. It is not valid to use Wildcard
// anywhere else.
const Wildcard string = "*"

// We'll read the default registry host from over in the addrs package, to
// avoid duplicating it. A "default" provider uses the default registry host
// by definition.
var defaultRegistryHost = addrs.DefaultRegistryHost

func validProviderNameOrWildcard(s string) bool {
	if s == Wildcard {
		return true
	}
	_, err := addrs.ParseProviderPart(s)
	return err == nil
}
