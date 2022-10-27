package mrpkg

import (
	"reflect"
	"testing"
)

type toNamedArgType struct{}

func (*toNamedArgType) ToNamedArgs() map[string]any {
	return map[string]any{
		"key1": 1,
		"key2": 2,
	}
}

type structArgType struct {
	Key int `db:"key3"`
}

func TestMergeNamedArgs(t *testing.T) {
	var (
		arg1 = new(toNamedArgType)
		arg2 = structArgType{
			Key: 3,
		}
		arg3 = map[string]any{
			"key4": 4,
			"key5": 5,
		}
	)

	expect := map[string]any{
		"key1": 1,
		"key2": 2,
		"key3": 3,
		"key4": 4,
		"key5": 5,
		"key6": 6,
	}

	got := MergeNamedArgs(map[string]any{
		"arg1": arg1,
		"arg2": arg2,
		"arg3": arg3,
		"key6": 6,
	})

	if !reflect.DeepEqual(expect, got) {
		t.Errorf("MergeNamedArgs: expect=%v; got=%v", expect, got)
	}
}
