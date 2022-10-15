package mrpkg

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"reflect"
	"strings"
)

var constructorMap ConcurrentMap[reflect.Type, func() any]

func init() {
	RegisterTypeConstructor[string](func() any { return "" })
	RegisterTypeConstructor[bool](func() any { return false })
	RegisterTypeConstructor[int](func() any { return 0 })
	RegisterTypeConstructor[int8](func() any { return 0 })
	RegisterTypeConstructor[int16](func() any { return 0 })
	RegisterTypeConstructor[int32](func() any { return 0 })
	RegisterTypeConstructor[int64](func() any { return 0 })
	RegisterTypeConstructor[uint](func() any { return 0 })
	RegisterTypeConstructor[uint8](func() any { return 0 })
	RegisterTypeConstructor[uint16](func() any { return 0 })
	RegisterTypeConstructor[uint32](func() any { return 0 })
	RegisterTypeConstructor[uint64](func() any { return 0 })
	RegisterTypeConstructor[float32](func() any { return 0.00 })
	RegisterTypeConstructor[float64](func() any { return 0.00 })
	RegisterTypeConstructor[complex64](func() any { return complex(0, 0) })
	RegisterTypeConstructor[complex128](func() any { return complex(0, 0) })
}

// New
//
//	valid -> mrpkg.New[type]()
//	valid -> mrpkg.New[*type]()
//	invalid -> mrpkg.New[**type]()
func New[T any]() (v T) {
	rv := reflect.ValueOf(&v).Elem()

	if construct, ok := constructorMap.Get(rv.Type()); ok {
		return construct().(T)
	}

	if !isNilKind(rv.Kind()) {
		return v
	}

	rv.Set(newType(rv.Type()))

	return v
}

func RegisterTypeConstructor[T any](construct func() any) {
	var v T
	constructorMap.Set(reflect.TypeOf(v), construct)
}

func isNilKind(k reflect.Kind) bool {
	if k == reflect.Invalid {
		panic("mrpkg.isNilKind: invalid reflect.Kind")
	}

	switch k {
	case reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.Pointer,
		reflect.Slice:
		return true
	}

	return false
}

func newType(typ reflect.Type) reflect.Value {
	switch typ.Kind() {
	case reflect.Slice:
		return reflect.MakeSlice(typ, 0, 0)
	case reflect.Map:
		return reflect.MakeMap(typ)
	case reflect.Chan:
		return reflect.MakeChan(typ, 0)
	case reflect.Func:
		return reflect.MakeFunc(typ, func(_ []reflect.Value) (results []reflect.Value) {
			results = make([]reflect.Value, typ.NumOut())
			for i := 0; i < typ.NumOut(); i++ {
				results[i] = newType(typ.Out(i))
			}
			return results
		})
	case reflect.Pointer:
		return reflect.New(typ.Elem())
	default:
		return reflect.Zero(typ)
	}
}

// IsInstance
// v: type/*type
// typ: type.Type
// return true while type/*type is a receiver of typ(also pointer receiver)
//
//	valid -> mrpkg.IsInstance(obj, type.Type)
//	valid -> mrpkg.IsInstance(obj, (*type).Type)
//	valid -> mrpkg.IsInstance(&obj, type.Type)
//	valid -> mrpkg.IsInstance(&obj, (*type).Type)
func IsInstance(v any, typ any) bool {
	if v == nil {
		return false
	}
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Pointer {
		rv = reflect.Indirect(rv)
	}
	return rv.Type() == typeOf(typ)
}

type ifaceType interface {
	Type()
}

var iface = reflect.TypeOf((*ifaceType)(nil)).Elem()

func typeOf(typObj any) reflect.Type {
	rt := reflect.TypeOf(typObj)
	if rt.Kind() != reflect.Func || rt.NumIn() != 1 {
		panic("typeOf: typeObj should be kind of 'XXX.Type' method")
	}

	typ := rt.In(0)
	if !typ.Implements(iface) {
		panic(fmt.Errorf("typeOf: type '%s' not implements 'ifaceType'", typ.String()))
	}

	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	return typ
}

func Max[T constraints.Ordered](x, y T) T {
	if x <= y {
		return y
	} else {
		return x
	}
}

func Min[T constraints.Ordered](x, y T) T {
	if x <= y {
		return x
	} else {
		return y
	}
}

func MinByKey[T any, U constraints.Ordered](x, y T, mapFunc func(item T) U) T {
	if mapFunc(x) <= mapFunc(y) {
		return x
	} else {
		return y
	}
}

func Sum[T constraints.Ordered](values ...T) (sum T) {
	for _, value := range values {
		sum += value
	}
	return sum
}

func Map[T, U any](src []T, f func(T) U) []U {
	dst := make([]U, len(src))
	for i := 0; i < len(src); i++ {
		dst[i] = f(src[i])
	}
	return dst
}

func Filter[T any](src []T, f func(T) bool) []T {
	dst := make([]T, 0, len(src))
	for i := 0; i < len(src); i++ {
		if f(src[i]) {
			dst = append(dst, src[i])
		}
	}
	return dst
}

func CollectMap[K comparable, T any](items []T, getKey func(T) K) (m map[K]T) {
	m = make(map[K]T, len(items))
	for _, item := range items {
		m[getKey(item)] = item
	}
	return m
}

func In[T comparable](item T, choices ...T) bool {
	for _, choice := range choices {
		if item == choice {
			return true
		}
	}
	return false
}

func Sequence[T any](items ...T) []T {
	return items
}

func Chunk[T any](slice []T, size int) (chunks [][]T) {
	if size <= 0 || size > len(slice) {
		return append(chunks, slice)
	}

	chunks = make([][]T, 0, len(slice)+size-1/size)
	for size < len(slice) {
		slice, chunks = slice[size:], append(chunks, slice[0:size:size])
	}

	chunks = append(chunks, slice)
	return chunks
}

func Merge[T any](slices ...[]T) (merged []T) {
	var capacity int
	for i := 0; i < len(slices); i++ {
		capacity += len(slices[i])
	}

	merged = make([]T, 0, capacity)
	for _, slice := range slices {
		for _, item := range slice {
			merged = append(merged, item)
		}
	}

	return merged
}

func Unique[T comparable](slice []T) (unique []T) {
	unique = make([]T, 0, len(slice))

	setMap := make(map[T]struct{}, len(slice))
	for _, item := range slice {
		setMap[item] = struct{}{}
	}

	for item := range setMap {
		unique = append(unique, item)
	}

	return unique
}

func Distinct[T any](list []T) []T {
	concurrentSet := new(ConcurrentSet[T])
	for i := 0; i < len(list); i++ {
		concurrentSet.Add(list[i])
	}
	return ToGoSlice(concurrentSet.ListIterator())
}

func UniqueByKey[T any, U comparable](slice []T, f func(T) U) (unique []T) {
	unique = make([]T, 0, len(slice))

	mapset := make(map[U]struct{}, len(slice))
	for _, item := range slice {
		key := f(item)
		if _, ok := mapset[key]; !ok {
			unique = append(unique, item)
			mapset[key] = struct{}{}
		}
	}

	return unique
}

func NewSlice[T any]() []T {
	return make([]T, 0)
}

func NewMap[K comparable, V any]() map[K]V {
	return make(map[K]V)
}

func GetOrInsertWith[K comparable, V any](m map[K]V, k K, New func() V) V {
	if _, ok := m[k]; !ok {
		m[k] = New()
	}

	return m[k]
}

func GetOrInsert[K comparable, V any](m map[K]V, k K, v V) V {
	if _, ok := m[k]; !ok {
		m[k] = v
	}

	return m[k]
}

func MapKeys[K comparable, V any](m map[K]V) (keys []K) {
	keys = make([]K, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

func Default[T any](currentValue, defaultValue T) T {
	if reflect.ValueOf(currentValue).IsZero() {
		return defaultValue
	}
	return currentValue
}

type ItemGetter interface {
	GetFrom(reflect.Value) reflect.Value
}

type Field string

func (field Field) GetFrom(value reflect.Value) reflect.Value {
	paths := strings.Split(string(field), ".")
	for i := 0; i < len(paths); i++ {
	Switch:
		switch kind := value.Kind(); kind {
		case reflect.Pointer:
			value = value.Elem()
			goto Switch
		case reflect.Struct:
			value = value.FieldByName(paths[i])
		default:
			panic(fmt.Errorf("ItemGetter: Field getter expects Struct, got %s", kind))
		}
	}
	return value
}

type Index int

func (index Index) GetFrom(value reflect.Value) reflect.Value {
	switch kind := value.Kind(); kind {
	case reflect.Array, reflect.Slice, reflect.String:
		return value.Index(int(index))
	default:
		panic(fmt.Errorf("ItemGetter: Index getter expects Array, Slice, or String, got %s", kind))
	}
}

func Getter[T, U any](getter ItemGetter) func(T) U {
	return func(item T) U {
		rv := reflect.Indirect(reflect.ValueOf(item))
		return getter.GetFrom(rv).Interface().(U)
	}
}

func Unwrap[T any](x T, err error) T {
	if err != nil {
		panic(err)
	}
	return x
}
