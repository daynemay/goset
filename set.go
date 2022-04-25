package goset

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/exp/constraints"
)

var exists = struct{}{}

// Set represents a (mathematical) set of values, supporting the set concepts of Union, Intersection, Difference
type Set[T constraints.Ordered] struct {
	members map[T]struct{}
}

// New returns a new Set, optionally initialized with some members
func New[T constraints.Ordered](members ...T) Set[T] {
	newSet := Set[T]{
		members: map[T]struct{}{},
	}
	for _, entry := range members {
		newSet.members[entry] = exists
	}
	return newSet
}

// String returns a string representation of theSet
func (theSet Set[T]) String() string {
	asListString := fmt.Sprintf("%#v", theSet.AsSortedList())
	return strings.Replace(asListString, "[]string", "Set", 1)
}

// Add adds a member to a Set, ignoring it if is already present
func (theSet Set[T]) Add(members ...T) Set[T] {
	for _, member := range members {
		theSet.members[member] = exists
	}
	return theSet
}

// Contains returns a boolean indicating whether theSet contains all the given strs
func (theSet Set[T]) Contains(values ...T) bool {
	for _, s := range values {
		if _, ok := theSet.members[s]; !ok {
			return false
		}
	}
	return true
}

// Equals returns a boolean indicating whether theSet is set-equal to other
func (theSet Set[T]) Equals(other Set[T]) bool {
	if theSet.Count() != other.Count() {
		return false
	}
	asList := theSet.AsList()
	for _, entry := range asList {
		if !other.Contains(entry) {
			return false
		}
	}
	return true
}

// AsList returns a slice of values in theSet
func (theSet Set[T]) AsList() []T {
	asList := []T{}
	for entry := range theSet.members {
		asList = append(asList, entry)
	}
	return asList
}

// AsSortedList returns a sorted slice of values in theSet
func (theSet Set[T]) AsSortedList() []T {
	asList := theSet.AsList()
	isLess := func(i, j int) bool {
		return asList[i] < asList[j]
	}
	sort.SliceStable(asList, isLess)
	return asList
}

// Intersect returns a new Set resulting from the set intersection of theSet and other
func (theSet Set[T]) Intersect(other Set[T]) Set[T] {
	commonMembers := []T{}
	for member := range theSet.members {
		if other.Contains(member) {
			commonMembers = append(commonMembers, member)
		}
	}
	return New(commonMembers...)
}

// Minus returns a new set representing the set difference theSet - other
func (theSet Set[T]) Minus(other Set[T]) Set[T] {
	difference := New[T]()
	for member := range theSet.members {
		if !other.Contains(member) {
			difference.Add(member)
		}
	}
	return difference
}

// Clone returns a copy of this Set
func (theSet Set[T]) Clone() Set[T] {
	return New(theSet.AsList()...)
}

// Union returns a new Set resulting from the set union of theSet and other
func (theSet Set[T]) Union(other Set[T]) Set[T] {
	union := theSet.Clone()
	union.Add(other.AsList()...)
	return union
}

func (theSet Set[T]) IsSubsetOf(other Set[T]) bool {
	return theSet.Intersect(other).Equals(theSet)
}

func (theSet Set[T]) IsProperSubsetOf(other Set[T]) bool {
	return theSet.IsSubsetOf(other) && !theSet.Equals(other)
}

func (theSet Set[T]) IsSupersetOf(other Set[T]) bool {
	return other.IsSubsetOf(theSet)
}

func (theSet Set[T]) IsProperSupersetOf(other Set[T]) bool {
	return theSet.IsSupersetOf(other) && !theSet.Equals(other)
}

// Count returns the set cardinality of theSet
func (theSet Set[T]) Count() int {
	return len(theSet.members)
}
