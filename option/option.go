package option

import (
	bytesPkg "bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

type Status bool

const (
	StatusSome = true
	StatusNone = false
)

type Option[T any] interface {
	Status() Status
	Unwrap() T
}

func (status Status) String() string {
	switch status {
	case StatusSome:
		return "Some"
	case StatusNone:
		return "None"
	default:
		return "None"
	}
}

func (status Status) IsSome() bool {
	return status == StatusSome
}

func (status Status) IsNone() bool {
	return status == StatusNone
}

func New[T any](value T) Value[T] {
	return Value[T]{
		valid: true,
		value: value,
	}
}

func NewRef[T any](value T) *Value[T] {
	return &Value[T]{
		valid: true,
		value: value,
	}
}

func NewNone[T any]() Value[T] {
	return Value[T]{
		valid: false,
		value: nil,
	}
}

func NewNoneRef[T any]() *Value[T] {
	return &Value[T]{
		valid: false,
		value: nil,
	}
}

type Value[T any] struct {
	valid bool
	value any
}

func (option *Value[T]) Status() Status {
	if option == nil {
		return StatusNone
	}
	return Status(option.valid)
}

func (option *Value[T]) Unwrap() T {
	if option == nil || !option.valid {
		panic("calling `Option.Unwrap` on a None value")
	}
	return option.value.(T)
}

func (option *Value[T]) Set(value T) {
	option.valid = true
	option.value = value
}

func (option *Value[T]) MarshalJSON() ([]byte, error) {
	if !option.valid {
		return json.Marshal(nil)
	}
	return json.Marshal(option.value)
}

func (option *Value[T]) UnmarshalJSON(bytes []byte) error {
	if bytesPkg.Equal(bytesPkg.TrimSpace(bytes), []byte("null")) {
		option.valid = false
		option.value = nil
		return nil
	}
	var x T
	option.valid = true
	if xt := reflect.TypeOf(x); xt.Kind() == reflect.Pointer {
		option.value = reflect.New(xt.Elem()).Interface().(T)
		return json.Unmarshal(bytes, option.value)
	} else {
		if err := json.Unmarshal(bytes, &x); err != nil {
			return err
		}
		option.value = x
		return nil
	}
}

func (option *Value[T]) Scan(src any) error {
	if src == nil {
		option.valid = false
		option.value = nil
		return nil
	}
	var x T
	option.valid = true
	var maybeScanner any
	if xt := reflect.TypeOf(x); xt.Kind() == reflect.Pointer {
		option.value = reflect.New(xt.Elem()).Interface().(T)
		maybeScanner = option.value
		if scanner, ok := maybeScanner.(sql.Scanner); ok {
			return scanner.Scan(src)
		}
	} else {
		maybeScanner = &x
		if scanner, ok := maybeScanner.(sql.Scanner); ok {
			if err := scanner.Scan(src); err != nil {
				return err
			}
			option.value = x
			return nil
		}
	}
	option.value = x
	switch v := src.(type) {
	case int64:
		switch option.value.(type) {
		case int:
			option.value = int(v)
			return nil
		case int64:
			option.value = v
			return nil
		case uint:
			option.value = uint(v)
			return nil
		case uint64:
			option.value = uint64(v)
			return nil
		}
	case float64:
		switch option.value.(type) {
		case float64:
			option.value = v
			return nil
		}
	case bool:
		switch option.value.(type) {
		case bool:
			option.value = v
			return nil
		}
	case []byte:
		switch option.value.(type) {
		case []byte:
			option.value = v
			return nil
		case string:
			option.value = string(v)
			return nil
		}
	case string:
		switch option.value.(type) {
		case []byte:
			option.value = []byte(v)
			return nil
		case string:
			option.value = v
			return nil
		}
	case time.Time:
		switch option.value.(type) {
		case time.Time:
			option.value = v
			return nil
		case *time.Time:
			option.value = &v
			return nil
		}
	}
	return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type %T", src, option.value)
}

func (option *Value[T]) Value() (driver.Value, error) {
	if IsNull[T](option) {
		return nil, nil
	}
	var maybeValuer any
	if reflect.TypeOf(option.value).Kind() == reflect.Pointer {
		x := option.value.(T)
		maybeValuer = &x
		if valuer, ok := maybeValuer.(driver.Valuer); ok {
			return valuer.Value()
		}
	}
	maybeValuer = option.value.(T)
	if valuer, ok := maybeValuer.(driver.Valuer); ok {
		return valuer.Value()
	}
	return option.value, nil
}

func Some[T any](value T) Option[T] {
	return &Value[T]{
		valid: true,
		value: value,
	}
}

func None[T any]() Option[T] {
	return &Value[T]{
		valid: false,
		value: nil,
	}
}

func IsNull[T any](option Option[T]) bool {
	return option == nil || option.Status().IsNone()
}

func IsNonNull[T any](option Option[T]) bool {
	return option != nil && option.Status().IsSome()
}

func Contains[T comparable](option Option[T], value T) bool {
	if IsNonNull(option) {
		return option.Unwrap() == value
	}
	return false
}

func Map[T any, U any](option Option[T], mapFunc func(item T) U) Option[U] {
	if IsNull(option) {
		return None[U]()
	}
	return Some(mapFunc(option.Unwrap()))
}

func And[T any](a Option[T], b Option[T]) Option[T] {
	if IsNull(a) {
		return a
	} else {
		return b
	}
}

func Or[T any](a Option[T], b Option[T]) Option[T] {
	if IsNonNull(a) {
		return a
	} else {
		return b
	}
}

func GetOrDefault[T any](option Option[T], defaultValue T) T {
	if IsNull(option) {
		return defaultValue
	}
	return option.Unwrap()
}
