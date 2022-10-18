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

type Value[T any] struct {
	Valid bool
	Item  any
}

func (option *Value[T]) Status() Status {
	return Status(option.Valid)
}

func (option *Value[T]) Unwrap() T {
	if !option.Valid {
		panic("calling `Option.Unwrap` on a None value")
	}
	return option.Item.(T)
}

func (option *Value[T]) MarshalJSON() ([]byte, error) {
	if !option.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(option.Value)
}

func (option *Value[T]) UnmarshalJSON(bytes []byte) error {
	if bytesPkg.Equal(bytesPkg.TrimSpace(bytes), []byte("null")) {
		option.Valid = false
		option.Item = nil
	}
	option.Valid = true
	if unmarshaler, ok := option.Item.(json.Unmarshaler); ok {
		return unmarshaler.UnmarshalJSON(bytes)
	}
	if reflect.TypeOf(option.Item).Kind() == reflect.Pointer {
		return json.Unmarshal(bytes, option.Item)
	} else {
		return json.Unmarshal(bytes, &option.Item)
	}
}

func (option *Value[T]) Scan(src any) error {
	if src == nil {
		option.Valid = false
		option.Item = nil
		return nil
	}
	option.Valid = true
	if scanner, ok := option.Item.(sql.Scanner); ok {
		return scanner.Scan(src)
	}
	switch v := src.(type) {
	case int64:
		switch option.Item.(type) {
		case int:
			option.Item = int(v)
		case int64:
			option.Item = v
		case uint:
			option.Item = uint(v)
		case uint64:
			option.Item = uint64(v)
		}
	case float64:
		switch option.Item.(type) {
		case float64:
			option.Item = v
		}
	case bool:
		switch option.Item.(type) {
		case bool:
			option.Item = v
		}
	case []byte:
		switch option.Item.(type) {
		case []byte:
			option.Item = v
		case string:
			option.Item = string(v)
		}
	case string:
		switch option.Item.(type) {
		case []byte:
			option.Item = []byte(v)
		case string:
			option.Item = v
		}
	case time.Time:
		switch option.Item.(type) {
		case time.Time:
			option.Item = v
		case *time.Time:
			option.Item = &v
		}
	}
	return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type %T", src, option.Item)
}

func (option *Value[T]) Value() (driver.Value, error) {
	if IsNull[T](option) {
		return nil, nil
	}
	if valuer, ok := option.Item.(driver.Valuer); ok {
		return valuer.Value()
	}
	return option.Item, nil
}

func Some[T any](value T) Option[T] {
	return &Value[T]{
		Valid: true,
		Item:  value,
	}
}

func None[T any]() Option[T] {
	return &Value[T]{
		Valid: false,
		Item:  nil,
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
