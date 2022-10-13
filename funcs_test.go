package mrpkg

import (
	"reflect"
	"testing"
)

type structType struct {
	Inner struct{}
}

func (structType) Type() {}

func TestNew(t *testing.T) {
	vv := New[structType]()
	vv.Inner = struct{}{}
	if !IsInstance(&vv, structType.Type) {
		t.Errorf("TestNew: value(structType) is not an instance of type 'structType'\n")
	}

	vp := New[*structType]()
	vp.Inner = struct{}{}
	if !IsInstance(&vp, structType.Type) {
		t.Errorf("TestNew: value(*structType) is not an instance of type 'structType'\n")
	}
}

type (
	type1 int
	type2 string
	type3 float64
	type4 bool
	type5 struct{}
	type6 []struct{}
	type7 map[string]struct{}
	type8 struct{}
)

func (type1) Type()  {}
func (type2) Type()  {}
func (type3) Type()  {}
func (type4) Type()  {}
func (type5) Type()  {}
func (type6) Type()  {}
func (type7) Type()  {}
func (*type8) Type() {}

var isTestCases = []struct {
	v   any
	typ any
}{
	{type1(0), type1.Type},
	{type2(""), type2.Type},
	{type3(0.00), type3.Type},
	{type4(false), type4.Type},
	{type5(struct{}{}), type5.Type},
	{type6(nil), type6.Type},
	{type7(nil), type7.Type},
	{type8(struct{}{}), (*type8).Type},
}

func TestIs(t *testing.T) {
	for _, testCase := range isTestCases {
		if !IsInstance(testCase.v, testCase.typ) {
			t.Errorf("TestIs(v=%s, typ=%s) failed\n", reflect.TypeOf(testCase.v), typeOf(testCase.typ))
		}
	}

	var (
		ty1 = type1(0)
		ty2 = type2("")
		ty3 = type3(0.00)
		ty4 = type4(false)
		ty5 = type5(struct{}{})
		ty6 = type6(nil)
		ty7 = type7(nil)
		ty8 = type8(struct{}{})
	)

	if !IsInstance(&ty1, type1.Type) {
		t.Errorf("TestIs(v=%s, typ=%s) failed\n", reflect.TypeOf(&ty1), typeOf(type1.Type))
	}

	if !IsInstance(&ty2, type2.Type) {
		t.Errorf("TestIs(v=%s, typ=%s) failed\n", reflect.TypeOf(&ty2), typeOf(type2.Type))
	}

	if !IsInstance(&ty3, type3.Type) {
		t.Errorf("TestIs(v=%s, typ=%s) failed\n", reflect.TypeOf(&ty3), typeOf(type3.Type))
	}

	if !IsInstance(&ty4, type4.Type) {
		t.Errorf("TestIs(v=%s, typ=%s) failed\n", reflect.TypeOf(&ty4), typeOf(type4.Type))
	}

	if !IsInstance(&ty5, type5.Type) {
		t.Errorf("TestIs(v=%s, typ=%s) failed\n", reflect.TypeOf(&ty5), typeOf(type5.Type))
	}

	if !IsInstance(&ty6, type6.Type) {
		t.Errorf("TestIs(v=%s, typ=%s) failed\n", reflect.TypeOf(&ty6), typeOf(type6.Type))
	}

	if !IsInstance(&ty7, type7.Type) {
		t.Errorf("TestIs(v=%s, typ=%s) failed\n", reflect.TypeOf(&ty7), typeOf(type7.Type))
	}

	if !IsInstance(&ty8, (*type8).Type) {
		t.Errorf("TestIs(v=%s, typ=%s) failed\n", reflect.TypeOf(&ty8), typeOf((*type8).Type))
	}
}

var (
	seq1 = []int{1, 2, 3}
	seq2 = [3]int{1, 2, 3}
	seq3 = "123"
)

func TestGetter(t *testing.T) {
	if v := Getter[[]int, int](Index(0))(seq1); v != 1 {
		t.Errorf("Getter[[]int, int](Index(0))(seq1) = %v", v)
	}

	if v := Getter[[3]int, int](Index(0))(seq2); v != 1 {
		t.Errorf("Getter[[3]int, int](Index(0))(seq2) = %v", v)
	}

	if v := Getter[string, byte](Index(0))(seq3); v != '1' {
		t.Errorf("Getter[string, byte](Index(0))(seq3) = %v", v)
	}
}
