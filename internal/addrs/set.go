package addrs

// Set represents a set of addresses of types that implement UniqueKeyer.
//
// Modify the set only by the methods on this type. This type exposes its
// internals for convenience during reading, such as iterating over set elements
// by ranging over the map values, but making direct modifications could
// potentially make the set data invalid or inconsistent, leading to undefined
// behavior elsewhere.
type Set[T UniqueKeyer] map[UniqueKey]UniqueKeyer

// Has returns true if and only if the set includes the given address.
func (s Set[T]) Has(addr UniqueKeyer) bool {
	_, exists := s[addr.UniqueKey()]
	return exists
}

// Add inserts the given address into the set, if not already present. If
// an equivalent address is already in the set, this replaces that address
// with the new value.
func (s Set[T]) Add(addr UniqueKeyer) {
	s[addr.UniqueKey()] = addr
}

// Remove deletes the given address from the set, if present. If not present,
// this is a no-op.
func (s Set[T]) Remove(addr UniqueKeyer) {
	delete(s, addr.UniqueKey())
}

// Union returns a new set which contains the union of all of the elements
// of both the reciever and the given other set.
func (s Set[T]) Union(other Set[T]) Set[T] {
	ret := make(Set[T])
	for k, addr := range s {
		ret[k] = addr
	}
	for k, addr := range other {
		ret[k] = addr
	}
	return ret
}

// Intersection returns a new set which contains the intersection of all of the
// elements of both the reciever and the given other set.
func (s Set[T]) Intersection(other Set[T]) Set[T] {
	ret := make(Set[T])
	for k, addr := range s {
		if _, exists := other[k]; exists {
			ret[k] = addr
		}
	}
	for k, addr := range other {
		if _, exists := s[k]; exists {
			ret[k] = addr
		}
	}
	return ret
}
