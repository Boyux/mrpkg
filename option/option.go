package option

import (
	"encoding/json"
)

type Status bool

const (
	StatusSome = true
	StatusNone = false
)

type Option[T any] interface {
	Status() Status
	Value() T
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
	Valid   bool
	Content T
}

func (option *Value[T]) Status() Status {
	return Status(option.Valid)
}

func (option *Value[T]) Value() T {
	return option.Content
}

func (option *Value[T]) MarshalJSON() ([]byte, error) {
	if !option.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(option.Content)
}

func Some[T any](value T) Option[T] {
	return &Value[T]{
		Valid:   true,
		Content: value,
	}
}

func None[T any]() Option[T] {
	return &Value[T]{Valid: false}
}

func IsNull[T any](option Option[T]) bool {
	return option == nil || option.Status().IsNone()
}

func IsNonNull[T any](option Option[T]) bool {
	return option != nil && option.Status().IsSome()
}

func Contains[T comparable](option Option[T], value T) bool {
	if IsNonNull(option) {
		return option.Value() == value
	}
	return false
}

func Map[T any, U any](option Option[T], mapFunc func(item T) U) Option[U] {
	if IsNull(option) {
		return None[U]()
	}
	return Some(mapFunc(option.Value()))
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
	return option.Value()
}
