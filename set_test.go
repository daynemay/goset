package goset

import (
	"fmt"
	"reflect"
	"testing"
)

func expect(t *testing.T, condition bool, description string, subs ...interface{}) {
	if !condition {
		t.Errorf(description, subs...)
	}
}

type Person struct {
	name string
	age  int
}

// A Comparator (see set.go) to order a collection of Person by age.
func byPersonAge(a, b interface{}) bool {
	p1, ok1 := a.(Person)
	p2, ok2 := b.(Person)
	if ok1 && ok2 {
		return p1.age < p2.age
	} else {
		return false
	}
}

func TestNew(t *testing.T) {

	t.Run("New should return an empty set by default", func(t *testing.T) {
		count := New[string]().Count()
		expect(t, count == 0, "NewSet().Count() = %v, expected 0", count)
	})

	t.Run("New should include supplied members", func(t *testing.T) {
		set := New("balrog", "blanka", "cammy", "guile")
		count := set.Count()
		expected := 4
		expect(t, count == expected, "NewSet(...).Count() = %v, expected %v", count, expected)

	})

	t.Run("New should ignore repeated supplied members", func(t *testing.T) {
		set := New("balrog", "balrog", "cammy", "balrog")
		count := set.Count()
		expected := 2
		expect(t, count == expected, "NewSet(...).Count() = %v, expected %v", count, expected)
	})

	t.Run("New can create a Set[int]", func(t *testing.T) {
		set := New(1, 2, 3, 4)
		count := set.Count()
		expected := 4
		expect(t, count == expected, "NewSet(...).Count() = %v, expected %v", count, expected)
	})

	t.Run("Try with (simple) composite type Person", func(t *testing.T) {
		people := []Person{
			{"Jeff", 58}, {"Rick", 55}, {"Kim", 3},
			{"Lara", 52}, {"Chris", 47}, {"Greg", 45},
		}
		set := NewWithComparator(byPersonAge, people...)
		count := set.Count()
		expected := len(people)
		expect(t, count == expected, "NewSet(...).Count() = %v, expected %v", count, expected)
		sl := set.AsSortedList()
		expect(t, count == expected, "SortedList Count() = %v, expected %v", len(sl), expected)
		listName := sl[0].name
		expectedName := "Kim"
		expect(t, listName == expectedName, "First name in sorted list = %v, expected %v",
			listName, expectedName)
		listName = sl[len(sl)-1].name
		expectedName = "Jeff"
		expect(t, listName == expectedName, "Last name in sorted list = %v, expected %v",
			listName, expectedName)
	})
}

func TestSet_String(t *testing.T) {
	t.Run("String() of an empty set", func(t *testing.T) {
		actual := New[string]().String()
		expected := "goset.Set[string]{}"
		expect(t, actual == expected, "Expected empty set to String to %s, was %v", expected, actual)
	})

	t.Run("String() shows ordered members", func(t *testing.T) {
		actual := New("ryu", "ken", "balrog", "cammy").String()
		expected := "goset.Set[string]{balrog, cammy, ken, ryu}"
		expect(t, actual == expected, "Expected String results to be ordered (%s), got %s", expected, actual)
	})

	t.Run("String() reflects the type of members", func(t *testing.T) {
		actual := New(1, 2, 3, 4).String()
		expected := "goset.Set[int]{1, 2, 3, 4}"
		expect(t, actual == expected, "Expected String results to be show the data type (%s), got %s", expected, actual)
	})

	t.Run("String() works like a String()", func(t *testing.T) {
		set := New("balrog", "cammy", "ken", "ryu")
		actual := fmt.Sprintf("%v", set)
		expected := "goset.Set[string]{balrog, cammy, ken, ryu}"
		expect(t, actual == expected, "Expected String() to be used for fmt.Sprintf: expected %s, was %v", expected, actual)
	})
}

func TestSet_AsSortedList(t *testing.T) {
	t.Run("Set[string].AsSortedList will return a sorted []string]", func(t *testing.T) {
		set := New("cammy", "ken", "ryu", "balrog")
		actual := set.AsSortedList()
		expected := []string{"balrog", "cammy", "ken", "ryu"}
		expect(t, reflect.DeepEqual(actual, expected), "Expected %v, got %v", expected, actual)
	})

	t.Run("Set[int].AsSortedList will return a sorted []int", func(t *testing.T) {
		set := New(44, -12, 3, -5, 0)
		actual := set.AsSortedList()
		expected := []int{-12, -5, 0, 3, 44}
		expect(t, reflect.DeepEqual(actual, expected), "Expected %v, got %v", expected, actual)
	})

	t.Run("Set[float].AsSortedList will return a sorted []int", func(t *testing.T) {
		set := New(44.44, -12.12, 3.3, -5.5, 0.0)
		actual := set.AsSortedList()
		expected := []float64{-12.12, -5.5, 0.0, 3.3, 44.44}
		expect(t, reflect.DeepEqual(actual, expected), "Expected %v, got %v", expected, actual)
	})

	t.Run("Set[bool].AsSortedList will return a sorted []bool", func(t *testing.T) {
		set := New(true, false)
		actual := set.AsSortedList()
		expected := []bool{false, true}
		expect(t, reflect.DeepEqual(actual, expected), "Expected %v, got %v", expected, actual)
	})

	t.Run("Set[struct{...}].AsSortedList will return a sorted []struct{...}", func(t *testing.T) {
		type fighter struct {
			name string
		}
		balrog := fighter{"balrog"}
		cammy := fighter{"cammy"}
		ken := fighter{"ken"}
		ryu := fighter{"ryu"}
		set := New(ken, ryu, cammy, balrog)
		actual := set.AsSortedList()
		expected := []fighter{balrog, cammy, ken, ryu}
		expect(t, reflect.DeepEqual(actual, expected), "Expected %v, got %v", expected, actual)
	})
}

func TestSet_Add(t *testing.T) {
	t.Run("Adding a new member should result increase the size of the set", func(t *testing.T) {
		set := New[string]()
		expect(t, set.Count() == 0, "Sanity check, expected to be empty")
		set.Add("guile")
		expect(t, set.Count() == 1, "Expect set to increase in size after Add()ing new member")
	})

	t.Run("Adding a new member should result in the presence of the member in the set", func(t *testing.T) {
		set := New[string]()
		expect(t, !set.Contains("guile"), "Expect set not to contain new member initially")
		set.Add("guile")
		expect(t, set.Contains("guile"), "Expect set to contain new member after Add()ing new member")
	})

	t.Run("Adding an existing member should not increase the size of the set", func(t *testing.T) {
		set := New("guile")
		expect(t, set.Count() == 1, "Expected to initially be size == 1")
		set.Add("guile")
		expect(t, set.Count() == 1, "Expected to remain at size == 1")
	})

	t.Run("Add returns the modified original set", func(t *testing.T) {
		airForce := New("guile")
		usa := airForce.Add("ken") // Ken joins the USAF!

		expect(t, airForce.Equals(usa), "Expected Add to return the modified, original set")
	})
}

func TestSet_Contains(t *testing.T) {
	t.Run("Set.Contains() members used to create it", func(t *testing.T) {
		set := New("balrog", "guile")
		expect(t, set.Contains("balrog"), "Expected set to contain balrog")
	})

	t.Run("Set.Contains() members Add()ed to it", func(t *testing.T) {
		set := New[string]()
		set.Add("balrog", "guile")
		expect(t, set.Contains("guile"), "Expected set to contain guile")
	})

	t.Run("Set.Contains() must contain all arguments", func(t *testing.T) {
		set := New[string]()
		set.Add("balrog", "guile")
		expect(t, !set.Contains("guile", "honda"), "Expected set not to contain guile-and-honda")
	})
}

func TestSet_Equals(t *testing.T) {
	t.Run("Two empty sets should be equal", func(t *testing.T) {
		expect(t, New[string]().Equals(New[string]()), "Expect two empty sets to be equal")
	})

	t.Run("Two sets with the same members should be Equal()", func(t *testing.T) {
		first := New("bison", "guile", "fei long")
		second := New("guile", "fei long", "bison")
		expect(t, first.Equals(second), "Expect two sets with the same members to be Equal()")
	})

	t.Run("Two sets with the same members should be Equal()", func(t *testing.T) {
		first := New("bison", "bison", "bison", "bison", "bison", "bison", "bison", "bison", "bison", "bison")
		second := New("bison")
		expect(t, first.Equals(second), "Expect repetition of members not to affect Equals()")
	})

	t.Run("Two sets with different members should not be Equal()", func(t *testing.T) {
		first := New("ryu", "ken", "honda")
		second := New("guile", "fei long", "bison")
		expect(t, !first.Equals(second), "Expect two different sets to not be Equal()")
	})
}

func TestSet_Intersect(t *testing.T) {
	t.Run("Set.Intersect() disparate sets should be empty", func(t *testing.T) {
		worldWarriors := New("ryu", "ken", "guile", "chun-li")
		bosses := New("balrog", "vega", "sagat", "bison")
		expect(t, worldWarriors.Intersect(bosses).Count() == 0, "Expected no intersection to disparate sets")
	})

	t.Run("Intersection with self should be equal to self", func(t *testing.T) {
		characters := New("ryu", "ken", "guile")
		intersection := characters.Intersect(characters)
		expect(t, intersection.Equals(characters), "Intersection with self should Equals(self)")
	})

	t.Run("Intersection with other should contain only the common members", func(t *testing.T) {
		first := New("ryu", "ken", "guile")
		second := New("ken", "guile", "balrog")
		intersection := first.Intersect(second)
		expected := New("guile", "ken")
		expect(t, intersection.Equals(expected), "Expected intersection to contain only common members")
	})
}

func TestSet_Union(t *testing.T) {
	t.Run("Union with disparate set contains all members of either set", func(t *testing.T) {
		worldWarriors := New("ryu", "ken", "guile", "chun-li")
		bosses := New("balrog", "vega", "sagat", "bison")
		expected := New("ryu", "ken", "guile", "chun-li", "balrog", "vega", "sagat", "bison")
		expect(t, worldWarriors.Union(bosses).Equals(expected), "Expected union to contain all members of either set")
	})

	t.Run("Union with self should be equal to self", func(t *testing.T) {
		characters := New("ryu", "ken", "guile")
		union := characters.Union(characters)
		expect(t, union.Equals(characters), "Union with self should Equals(self)")
	})
}

func TestSet_Minus(t *testing.T) {
	t.Run("Empty set minus empty set should be empty set", func(t *testing.T) {
		empty := New[string]()
		difference := empty.Minus(empty)
		expect(t, difference.Equals(empty), "Empty set minus empty set should be empty set")
	})

	t.Run("Empty set minus non-empty set should be empty set", func(t *testing.T) {
		empty := New[string]()
		nonEmpty := New("dhalsim", "honda", "vega")
		difference := empty.Minus(nonEmpty)
		expect(t, difference.Equals(empty), "Empty set minus non-empty set should be empty set")
	})

	t.Run("Non-empty minus empty set should be equal to the original set", func(t *testing.T) {
		nonEmpty := New("dhalsim", "honda", "vega")
		empty := New[string]()
		difference := nonEmpty.Minus(empty)
		expect(t, difference.Equals(nonEmpty), "Original set minus empty should be equal to original set")
	})

	t.Run("Set difference should return the members in the first that are not present in the second", func(t *testing.T) {
		first := New("ken", "honda", "ryu")
		second := New("honda", "chun-li", "cammy")
		difference := first.Minus(second)
		expected := New("ken", "ryu")
		expect(t, difference.Equals(expected), "First minus second should return elements in first that are not in second")
	})
}

func TestSet_Clone(t *testing.T) {
	t.Run("Clone of empty set should be empty set", func(t *testing.T) {
		empty := New[string]()
		clone := empty.Clone()
		expect(t, clone.Count() == 0, "Clone of empty set should be of size zero")
	})

	t.Run("Clone of set should be equal to original", func(t *testing.T) {
		original := New("balrog", "vega", "sagat", "bison")
		clone := original.Clone()
		expect(t, clone.Count() == original.Count(), "Clone of set should be of same size")
		expect(t, clone.Equals(original), "Clone of set should be equal to original")
	})

	t.Run("Mutation of clone should not affect original", func(t *testing.T) {
		original := New[string]()
		clone := original.Clone()
		clone.Add("deejay")
		expect(t, clone.Contains("deejay"), "Clone should contain new member")
		expect(t, !original.Contains("deejay"), "Modification of clone should not change the original set")
	})

	t.Run("Mutation of original should not affect clone", func(t *testing.T) {
		original := New[string]()
		clone := original.Clone()
		original.Add("cammy")
		expect(t, original.Contains("cammy"), "Original should contain new member")
		expect(t, !clone.Contains("cammy"), "Modification of original should not change the clone")
	})
}

func TestSet_IsSubsetOf(t *testing.T) {
	t.Run("Empty set is a subset of empty set", func(t *testing.T) {
		empty := New[string]()
		otherEmpty := New[string]()
		expect(t, empty.IsSubsetOf(otherEmpty), "Empty set should be a subset of empty set")
	})

	t.Run("Empty set is a subset of a non-empty set", func(t *testing.T) {
		empty := New[string]()
		nonEmpty := New("dhalsim", "honda", "vega")
		expect(t, empty.IsSubsetOf(nonEmpty), "Empty set should be a subset of non-empty set")
	})

	t.Run("Proper subset is a subset of a proper super set", func(t *testing.T) {
		sub := New("dhalsim", "honda")
		super := New("dhalsim", "honda", "vega")
		expect(t, sub.IsSubsetOf(super), "Populated proper subset should be identified as a subset of a superset")
	})

	t.Run("Equal sets are subsets of each other", func(t *testing.T) {
		a := New("dhalsim", "honda", "vega")
		b := New("dhalsim", "honda", "vega")
		expect(t, a.Equals(b), "sanity check: a and b should be equal")
		expect(t, a.IsSubsetOf(b), "a should be a subset of equal set b")
		expect(t, b.IsSubsetOf(a), "b should be a subset of equal set a")
	})

	t.Run("Element in A missing from B prevents A from being a subset of B", func(t *testing.T) {
		sub := New("ken", "honda")
		super := New("dhalsim", "honda", "vega")
		expect(t, !sub.IsSubsetOf(super), "Element in %s should prevent it from being a subset of %s", sub, super)
	})
}

func TestSet_IsProperSubsetOf(t *testing.T) {
	t.Run("Empty set is not a proper subset of empty set", func(t *testing.T) {
		empty := New[string]()
		otherEmpty := New[string]()
		expect(t, !empty.IsProperSubsetOf(otherEmpty), "Empty set is not a proper subset of empty set")
	})

	t.Run("Empty set is a proper subset of a non-empty set", func(t *testing.T) {
		empty := New[string]()
		nonEmpty := New("dhalsim", "honda", "vega")
		expect(t, empty.IsProperSubsetOf(nonEmpty), "Empty set is a proper subset of a non-empty set")
	})

	t.Run("Proper subset is correctly identified", func(t *testing.T) {
		sub := New("dhalsim", "honda")
		super := New("dhalsim", "honda", "vega")
		expect(t, sub.IsProperSubsetOf(super), "Proper subset is correctly identified")
	})

	t.Run("Equal sets are not proper subsets of each other", func(t *testing.T) {
		a := New("dhalsim", "honda", "vega")
		b := New("dhalsim", "honda", "vega")
		expect(t, a.Equals(b), "sanity check: a and b should be equal")
		expect(t, !a.IsProperSubsetOf(b), "a should be a subset of equal set b")
		expect(t, !b.IsProperSubsetOf(a), "b should be a subset of equal set a")
	})

	t.Run("Element in A missing from B prevents A from being a proper subset of B", func(t *testing.T) {
		sub := New("ken", "honda")
		super := New("dhalsim", "honda", "vega")
		expect(t, !sub.IsProperSubsetOf(super), "Element in %s should prevent it from being a subset of %s", sub, super)
	})
}

func TestSet_IsSupersetOf(t *testing.T) {
	t.Run("Empty set is a superset of empty set", func(t *testing.T) {
		empty := New[string]()
		otherEmpty := New[string]()
		expect(t, empty.IsSupersetOf(otherEmpty), "Empty set should be a superset of empty set")
	})

	t.Run("Non-empty is a superset of empty set", func(t *testing.T) {
		empty := New[string]()
		nonEmpty := New("dhalsim", "honda", "vega")
		expect(t, nonEmpty.IsSupersetOf(empty), "Non-empty should be a superset of empty set")
	})

	t.Run("Proper superset is identified as a superset", func(t *testing.T) {
		sub := New("dhalsim", "honda")
		super := New("dhalsim", "honda", "vega")
		expect(t, super.IsSupersetOf(sub), "Populated proper subset should be identified as a superset of a proper subset")
	})

	t.Run("Equal sets are supersets of each other", func(t *testing.T) {
		a := New("dhalsim", "honda", "vega")
		b := New("dhalsim", "honda", "vega")
		expect(t, a.Equals(b), "sanity check: a and b should be equal")
		expect(t, a.IsSupersetOf(b), "a should be a superset of equal set b")
		expect(t, b.IsSupersetOf(a), "b should be a superset of equal set a")
	})

	t.Run("Element in A missing from B prevents B from being a superset of A", func(t *testing.T) {
		sub := New("ken", "honda")
		super := New("dhalsim", "honda", "vega")
		expect(t, !super.IsSupersetOf(sub), "Element in %s should prevent %s from being a superset", sub, super)
	})
}

func TestSet_IsProperSupersetOf(t *testing.T) {
	t.Run("Empty set is not a proper superset of empty set", func(t *testing.T) {
		empty := New[string]()
		otherEmpty := New[string]()
		expect(t, !empty.IsProperSupersetOf(otherEmpty), "Empty set is not a proper superset of empty set")
	})

	t.Run("Non-empty set is a proper superset of empty set", func(t *testing.T) {
		empty := New[string]()
		nonEmpty := New("dhalsim", "honda", "vega")
		expect(t, nonEmpty.IsProperSupersetOf(empty), "Non-empty set is a proper superset of empty set")
	})

	t.Run("Proper superset is correctly identified", func(t *testing.T) {
		sub := New("dhalsim", "honda")
		super := New("dhalsim", "honda", "vega")
		expect(t, super.IsProperSupersetOf(sub), "Proper superset is correctly identified")
	})

	t.Run("Equal sets are not proper supersets of each other", func(t *testing.T) {
		a := New("dhalsim", "honda", "vega")
		b := New("dhalsim", "honda", "vega")
		expect(t, a.Equals(b), "sanity check: a and b should be equal")
		expect(t, !a.IsProperSupersetOf(b), "a should be a subset of equal set b")
		expect(t, !b.IsProperSupersetOf(a), "b should be a subset of equal set a")
	})

	t.Run("Element in A missing from B prevents B from being a proper superset of A", func(t *testing.T) {
		sub := New("ken", "honda")
		super := New("dhalsim", "honda", "vega")
		expect(t, !super.IsSupersetOf(sub), "Element in %s should prevent %s from being a superset", sub, super)
	})
}
