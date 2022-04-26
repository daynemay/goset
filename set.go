package goset

import (
	"fmt"
	"strings"

	"golang.org/x/exp/maps"
)

var exists = struct{}{}

type Comparator func(a, b interface{}) int

// Set represents a (mathematical) set of values, supporting the set concepts of Union, Intersection, Difference
type Set[T comparable] struct {
	members    map[T]struct{}
	comparator Comparator
}

// New returns a new Set, optionally initialized with some members
func New[T comparable](members ...T) Set[T] {
	newSet := Set[T]{
		members:    map[T]struct{}{},
		comparator: nil,
	}
	for _, entry := range members {
		newSet.members[entry] = exists
	}

	return newSet
}

func NewWithComparator[T comparable](c Comparator, members ...T) Set[T] {
	newSet := Set[T]{
		members:    map[T]struct{}{},
		comparator: c,
	}
	for _, entry := range members {
		newSet.members[entry] = exists
	}

	return newSet
}

// String returns a string representation of theSet
func (theSet Set[T]) String() string {
	var sb strings.Builder
	members := theSet.AsSortedList()

	sb.WriteString(fmt.Sprintf("%T{", theSet))
	for idx, value := range members {
		sb.WriteString(fmt.Sprintf("%v", value))
		if idx < len(members)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("}")
	return sb.String()
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
	return maps.Keys(theSet.members)
}

// AsSortedList returns a slice of values in theSet in a stable sorted order.
func (theSet Set[T]) AsSortedList() []T {
	return sortComparable(theSet.AsList())
}

/*
func (theSet Set[T]) AsSortedList() []T {
	asList := theSet.AsList()
	var isLess func(i, j int) bool

	if theSet.comparator != nil {
		isLess = func(i, j int) bool {
			return theSet.comparator(asList[i], asList[j]) < 0
		}
	} else {
		isLess = func(i, j int) bool {
			const bitSize = 64

			si := fmt.Sprintf("%#v", asList[i])
			sj := fmt.Sprintf("%#v", asList[j])
			fi, erri := strconv.ParseFloat(si, bitSize)
			fj, errj := strconv.ParseFloat(sj, bitSize)
			if erri == nil && errj == nil {
				return fi < fj
			} else {
				return si < sj
			}
		}
		/*
			isLess = func(i, j int) bool {
				ifi := ((interface{})(asList[i]))
				ifj := ((interface{})(asList[j]))

				switch ti := ifi.(type) {
				case string:
					si, oki := (ifi).(string)
					sj, okj := (ifj).(string)
					if oki && okj {
						return si < sj
					}
					break

				default:
					break
				}
				return true
	}
	sort.SliceStable(asList, isLess)
	return asList
}

// AsSortedList returns a sorted slice of values in theSet
/*
func (theSet Set[T]) AsSortedList(sif sort.Interface) []T {
	asList := theSet.AsList()

	if sif != nil {
		sort.SliceStable(asList, sif.Less)
		return asList
	}

	isLess := func(i, j int) bool {
		const bitSize = 64

		si := fmt.Sprintf("%#v", asList[i])
		sj := fmt.Sprintf("%#v", asList[j])
		fi, erri := strconv.ParseFloat(si, bitSize)
		fj, errj := strconv.ParseFloat(sj, bitSize)
		if erri == nil && errj == nil {
			return fi < fj
		} else {
			return si < sj
		}
	}
	sort.SliceStable(asList, isLess)
	return asList
}
*/

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
