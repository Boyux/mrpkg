package mrpkg

import (
	"reflect"
	"sync"
)

var poolMap = ConcurrentMap[reflect.Type, *sync.Pool]{
	DefaultFunc: func(typ reflect.Type) *sync.Pool {
		return &sync.Pool{
			New: func() any {
				return newType(typ).Interface()
			},
		}
	},
}

func GetObj[T any]() (v T) {
	return poolMap.
		GetOrDefault(
			reflect.TypeOf(v),
		).
		Get().(T)
}

func PutObj[T any](x T) {
	poolMap.
		GetOrDefault(
			reflect.TypeOf(x),
		).
		Put(x)
}
