package goset

import (
	"fmt"
	"sort"
	"strings"
)

var exists = struct{}{}

// Set represents a (mathematical) set of strings, supporting the set concepts of Union, Intersection, Difference
type Set struct {
	members map[string]struct{}
}

// New returns a new Set, optionally initialized with some members
func New(members ...string) Set {
	newSet := Set{
		members: map[string]struct{}{},
	}
	for _, entry := range members {
		newSet.members[entry] = exists
	}
	return newSet
}

// String returns a string representation of theSet
func (theSet Set) String() string {
	asListString := fmt.Sprintf("%#v", theSet.AsSortedList())
	return strings.Replace(asListString, "[]string", "Set", 1)
}

// Add adds a member to a Set, ignoring it if is already present
func (theSet Set) Add(members ...string) Set {
	for _, member := range members {
		theSet.members[member] = exists
	}
	return theSet
}

// Contains returns a boolean indicating whether theSet contains all the given strs
func (theSet Set) Contains(strs ...string) bool {
	for _, s := range strs {
		if _, ok := theSet.members[s]; !ok {
			return false
		}
	}
	return true
}

// Equals returns a boolean indicating whether theSet is set-equal to other
func (theSet Set) Equals(other Set) bool {
	asList := theSet.AsList()
	sort.Strings(asList)
	otherAsList := other.AsList()
	sort.Strings(otherAsList)
	if len(asList) != len(otherAsList) {
		return false
	}
	for idx := range asList {
		if asList[idx] != otherAsList[idx] {
			return false
		}
	}
	return true
}

// AsList returns a slice of strings in theSet
func (theSet Set) AsList() []string {
	asList := []string{}
	for entry := range theSet.members {
		asList = append(asList, entry)
	}
	return asList
}

// AsSortedList returns a sorted slice of strings in theSet
func (theSet Set) AsSortedList() []string {
	asList := theSet.AsList()
	sort.Strings(asList)
	return asList
}

// Intersect returns a new Set resulting from the set intersection of theSet and other
func (theSet Set) Intersect(other Set) Set {
	commonMembers := []string{}
	for member := range theSet.members {
		if other.Contains(member) {
			commonMembers = append(commonMembers, member)
		}
	}
	return New(commonMembers...)
}

// Minus returns a new set representing the set difference theSet - other
func (theSet Set) Minus(other Set) Set {
	difference := New()
	for member := range theSet.members {
		if !other.Contains(member) {
			difference.Add(member)
		}
	}
	return difference
}

// Clone returns a copy of this Set
func (theSet Set) Clone() Set {
	return New(theSet.AsList()...)
}

// Union returns a new Set resulting from the set union of theSet and other
func (theSet Set) Union(other Set) Set {
	union := theSet.Clone()
	union.Add(other.AsList()...)
	return union
}

func (theSet Set) IsSubsetOf(other Set) bool {
	return theSet.Intersect(other).Equals(theSet)
}

func (theSet Set) IsProperSubsetOf(other Set) bool {
	return theSet.IsSubsetOf(other) && !theSet.Equals(other)
}

func (theSet Set) IsSupersetOf(other Set) bool {
	return other.IsSubsetOf(theSet)
}

func (theSet Set) IsProperSupersetOf(other Set) bool {
	return theSet.IsSupersetOf(other) && !theSet.Equals(other)
}

// Count returns the set cardinality of theSet
func (theSet Set) Count() int {
	return len(theSet.members)
}
