package mrpkg

import (
	"reflect"
	"testing"
)

func TestVector_SortBy(t *testing.T) {
	var vector Vector[int]
	vector.Append(Iter([]int{1, 2, 3, 4, 5}))
	vector.SortBy(func(l, r int) bool {
		return l >= r
	})
	expect, got := []int{5, 4, 3, 2, 1}, ToGoSlice(vector.ListIterator())
	if !reflect.DeepEqual(got, expect) {
		t.Errorf("Vector.SortBy: \n\texpect=%v; \n\tgot=%v;\n",
			expect, got)
	}
}

func TestVector_Insert(t *testing.T) {
	var vector Vector[int]
	vector.Append(Iter([]int{1, 2, 3, 4, 5}))
	vector.Insert(1, 2)
	expect, got := []int{1, 2, 2, 3, 4, 5}, ToGoSlice(vector.ListIterator())
	if !reflect.DeepEqual(got, expect) {
		t.Errorf("Vector.Insert: \n\texpect=%v; \n\tgot=%v;\n",
			expect, got)
	}
}
