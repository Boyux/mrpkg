package mrpkg

import (
	"reflect"
	"sort"
	"testing"
)

var (
	list1 = []int{1, 2, 3, 4}
	list2 = []int{2, 3, 4, 5}
)

func TestConcurrentSet_Union(t *testing.T) {
	var set1, set2 ConcurrentSet[int]
	set1.BatchAdd(Iter(list1))
	set2.BatchAdd(Iter(list2))
	target := ToGoSlice(set1.Union(&set2).ListIterator())
	sort.Ints(target)
	if !reflect.DeepEqual(target, []int{1, 2, 3, 4, 5}) {
		t.Errorf("ConcurrentSet.Union: \n\tset1=%v; \n\tset2=%v; \n\ttarget=%v",
			ToGoSlice(set1.ListIterator()),
			ToGoSlice(set2.ListIterator()),
			target)
	}
}

func TestConcurrentSet_Intersection(t *testing.T) {
	var set1, set2 ConcurrentSet[int]
	set1.BatchAdd(Iter(list1))
	set2.BatchAdd(Iter(list2))
	target := ToGoSlice(set1.Intersection(&set2).ListIterator())
	sort.Ints(target)
	if !reflect.DeepEqual(target, []int{2, 3, 4}) {
		t.Errorf("ConcurrentSet.Intersection: \n\tset1=%v; \n\tset2=%v; \n\ttarget=%v",
			ToGoSlice(set1.ListIterator()),
			ToGoSlice(set2.ListIterator()),
			target)
	}
}

func TestConcurrentSet_Difference(t *testing.T) {
	var set1, set2 ConcurrentSet[int]
	set1.BatchAdd(Iter(list1))
	set2.BatchAdd(Iter(list2))

	target1 := ToGoSlice(set1.Difference(&set2).ListIterator())
	sort.Ints(target1)
	if !reflect.DeepEqual(target1, []int{1}) {
		t.Errorf("ConcurrentSet.Difference: \n\tset1=%v; \n\tset2=%v; \n\ttarget1=%v",
			ToGoSlice(set1.ListIterator()),
			ToGoSlice(set2.ListIterator()),
			target1)
	}

	target2 := ToGoSlice(set2.Difference(&set1).ListIterator())
	sort.Ints(target2)
	if !reflect.DeepEqual(target2, []int{5}) {
		t.Errorf("ConcurrentSet.Difference: \n\tset2=%v; \n\tset1=%v; \n\ttarget2=%v",
			ToGoSlice(set2.ListIterator()),
			ToGoSlice(set1.ListIterator()),
			target2)
	}
}

func TestConcurrentSet_SymmetricDifference(t *testing.T) {
	var set1, set2 ConcurrentSet[int]
	set1.BatchAdd(Iter(list1))
	set2.BatchAdd(Iter(list2))
	target := ToGoSlice(set1.SymmetricDifference(&set2).ListIterator())
	sort.Ints(target)
	if !reflect.DeepEqual(target, []int{1, 5}) {
		t.Errorf("ConcurrentSet.SymmetricDifference: \n\tset1=%v; \n\tset2=%v; \n\ttarget=%v",
			ToGoSlice(set1.ListIterator()),
			ToGoSlice(set2.ListIterator()),
			target)
	}
}
